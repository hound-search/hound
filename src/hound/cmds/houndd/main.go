package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"hound/api"
	"hound/config"
	"hound/searcher"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	info_log  *log.Logger
	error_log *log.Logger
)

type content struct {
	template string
	dest     string
	sources  []string
}

const (
	ReactVersion  = "0.12.2"
	jQueryVersion = "2.1.3"
)

func checkForJsx() error {
	return exec.Command("jsx", "--version").Run()
}

func (c *content) render(w io.Writer, root string, cfg *config.Config, prod bool) error {
	t, err := template.ParseFiles(filepath.Join(root, c.template))
	if err != nil {
		return err
	}

	json, err := cfg.ToJsonString()
	if err != nil {
		return err
	}

	var src template.HTML
	if prod {
		s, err := sourceForPrd(root, c.sources)
		if err != nil {
			return err
		}
		src = s
	} else {
		src = sourceForDev(c.sources)
	}

	return t.Execute(w, map[string]interface{}{
		"ReactVersion":  ReactVersion,
		"jQueryVersion": jQueryVersion,
		"ReposAsJson":   json,
		"Source":        src,
	})
}

func (c *content) renderToFile(filename string, root string, cfg *config.Config, prod bool) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	return c.render(w, root, cfg, prod)
}

type devHandler struct {
	http.Handler
	root    string
	cfg     *config.Config
	content map[string]*content
}

func sourceForPrd(root string, paths []string) (template.HTML, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "<script>")
	for _, path := range paths {
		cmd := exec.Command("jsx", filepath.Join(root, path))
		r, err := cmd.StdoutPipe()
		if err != nil {
			return "", err
		}

		if err := cmd.Start(); err != nil {
			return "", err
		}

		if _, err := io.Copy(&buf, r); err != nil {
			return "", err
		}

		if err := cmd.Wait(); err != nil {
			return "", err
		}
	}
	fmt.Fprintln(&buf, "</script>")
	return template.HTML(buf.String()), nil
}

func sourceForDev(paths []string) template.HTML {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "<script src=\"/assets/js/JSXTransformer-%s.js\"></script>\n", ReactVersion)
	for _, path := range paths {
		fmt.Fprintf(&buf, "<script type=\"text/jsx\" src=\"/%s\"></script>", path)
	}
	return template.HTML(buf.String())
}

func (h *devHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	c := h.content[p]
	if c == nil {
		h.Handler.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if err := c.render(w, h.root, h.cfg, false); err != nil {
		panic(err)
	}
}

func BuildContentFor(root string, prod bool, cnts []*content, cfg *config.Config) (http.Handler, error) {
	if prod {
		for _, cnt := range cnts {
			if err := cnt.renderToFile(filepath.Join(root, cnt.dest), root, cfg, prod); err != nil {
				return nil, err
			}
		}

		return http.FileServer(http.Dir(root)), nil
	}

	m := map[string]*content{}
	for _, cnt := range cnts {
		if strings.HasSuffix(cnt.dest, "index.html") {
			m[path.Clean("/"+filepath.Dir(cnt.dest)+"/")] = cnt
		} else {
			m["/"+cnt.dest] = cnt
		}
	}

	return &devHandler{
		Handler: http.FileServer(http.Dir(root)),
		root:    root,
		cfg:     cfg,
		content: m,
	}, nil
}

func makeSearchers(cfg *config.Config) (map[string]*searcher.Searcher, bool, error) {
	// Ensure we have a dbpath
	if _, err := os.Stat(cfg.DbPath); err != nil {
		if err := os.MkdirAll(cfg.DbPath, os.ModePerm); err != nil {
			return nil, false, err
		}
	}

	searchers, errs, err := searcher.MakeAll(cfg)
	if err != nil {
		return nil, false, err
	}

	if len(errs) > 0 {
		// NOTE: This mutates the original config so the repos
		// are not even seen by other code paths.
		for name, _ := range errs {
			delete(cfg.Repos, name)
		}

		return searchers, false, errors.New("One or more repos failed to index")
	}

	return searchers, true, nil
}

func makeTemplateData(cfg *config.Config) (interface{}, error) {
	var data struct {
		ReposAsJson string
	}

	res := map[string]*config.Repo{}
	for name, repo := range cfg.Repos {
		res[strings.ToLower(name)] = repo
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	data.ReposAsJson = string(b)
	return &data, nil
}

func runHttp(addr, root string, prod bool, cfg *config.Config, idx map[string]*searcher.Searcher) error {
	m := http.DefaultServeMux

	contents := []*content{
		&content{
			template: "index.tpl.html",
			dest:     "index.html",
			sources: []string{
				"assets/js/common.js",
				"assets/js/hound.js",
			},
		},
		&content{
			template: "excluded_files.tpl.html",
			dest:     "excluded_files.html",
			sources: []string{
				"assets/js/common.js",
				"assets/js/excluded_files.js",
			},
		},
	}

	handler, err := BuildContentFor(
		filepath.Join(root, "pub"),
		prod,
		contents,
		cfg)
	if err != nil {
		return err
	}

	m.Handle("/", handler)

	api.Setup(m, idx)
	return http.ListenAndServe(addr, m)
}

func findRoot(root *string) error {
	if *root == "" {
		return nil
	}

	_, file, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(
		filepath.Join(filepath.Dir(file), "../../"))
	if err != nil {
		return err
	}

	*root = dir
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	info_log = log.New(os.Stdout, "", log.LstdFlags)
	error_log = log.New(os.Stderr, "", log.LstdFlags)

	flagConf := flag.String("conf", "config.json", "")
	flagAddr := flag.String("addr", ":6080", "")
	flagRoot := flag.String("root", "", "")
	flagProd := flag.Bool("prod", false, "")

	flag.Parse()

	// In prod mode, we will need jsx.
	if *flagProd {
		if err := checkForJsx(); err != nil {
			panic("You need to install jsx. (npm install -g react-tools)")
		}
	}

	if err := findRoot(flagRoot); err != nil {
		panic(err)
	}

	var cfg config.Config
	if err := cfg.LoadFromFile(*flagConf); err != nil {
		panic(err)
	}

	idx, ok, err := makeSearchers(&cfg)
	if err != nil {
		log.Panic(err)
	}
	if !ok {
		info_log.Println("Some repos failed to index, see output above")
	} else {
		info_log.Println("All indexes built!")
	}

	host := *flagAddr
	if strings.HasPrefix(host, ":") {
		host = "localhost" + host
	}

	info_log.Printf("running server at http://%s...\n", host)

	if err := runHttp(*flagAddr, *flagRoot, *flagProd, &cfg, idx); err != nil {
		panic(err)
	}
}
