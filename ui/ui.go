package ui

import (
	"bytes"
	"fmt"
	"github.com/etsy/hound/config"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"runtime"
)

type devHandler struct {
	http.Handler
	content map[string]*content
	root    string
	cfg     *config.Config
}

type prdHandler struct {
	content map[string]*content
	cfgJson string
	cfg     *config.Config
}

func (h *devHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	cr := h.content[p]
	if cr == nil {
		h.Handler.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if err := renderForDev(w, h.root, cr, h.cfg); err != nil {
		log.Panic(err)
	}
}

func renderForDev(w io.Writer, root string, c *content, cfg *config.Config) error {
	t, err := template.ParseFiles(
		filepath.Join(root, c.template))
	if err != nil {
		return err
	}

	json, err := cfg.ToJsonString()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	fmt.Fprintf(
		&buf,
		"<script src=\"/js/JSXTransformer-%s.js\"></script>\n",
		ReactVersion)
	for _, path := range c.sources {
		fmt.Fprintf(
			&buf,
			"<script type=\"text/jsx\" src=\"/%s\"></script>",
			path)
	}

	return t.Execute(w, map[string]interface{}{
		"ReactVersion":  ReactVersion,
		"jQueryVersion": JQueryVersion,
		"ReposAsJson":   json,
		"Source":        template.HTML(buf.String()),
	})
}

func (h *prdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	ct := h.content[p]
	if ct != nil {
		if err := renderForPrd(w, ct, h.cfgJson); err != nil {
			log.Panic(err)
		}
		return
	}

	a, err := Asset(p[1:])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	mt := mime.TypeByExtension(
		filepath.Ext(p))
	if mt != "" {
		w.Header().Set("Content-Type", mt)
	}

	if _, err := w.Write(a); err != nil {
		log.Panic(err)
	}
}

func renderForPrd(w io.Writer, c *content, cfgJson string) error {

	var buf bytes.Buffer
	buf.WriteString("<script>")
	for _, src := range c.sources {
		a, err := Asset(src)
		if err != nil {
			return err
		}
		buf.Write(a)
	}
	buf.WriteString("</script>")

	return c.tpl.Execute(w, map[string]interface{}{
		"ReactVersion":  ReactVersion,
		"jQueryVersion": JQueryVersion,
		"ReposAsJson":   cfgJson,
		"Source":        template.HTML(buf.String()),
	})
}

func assetDir() string {
	_, file, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(
		filepath.Join(filepath.Dir(file), "assets"))
	if err != nil {
		log.Panic(err)
	}
	return dir
}

func newDevHandler(cfg *config.Config) (http.Handler, error) {
	root := assetDir()
	return &devHandler{
		Handler: http.FileServer(http.Dir(root)),
		content: contents,
		root:    root,
		cfg:     cfg,
	}, nil
}

func newPrdHandler(cfg *config.Config) (http.Handler, error) {
	for _, cnt := range contents {
		a, err := Asset(cnt.template)
		if err != nil {
			return nil, err
		}

		cnt.tpl, err = template.New(cnt.template).Parse(string(a))
		if err != nil {
			return nil, err
		}
	}

	json, err := cfg.ToJsonString()
	if err != nil {
		return nil, err
	}

	return &prdHandler{
		content: contents,
		cfg:     cfg,
		cfgJson: json,
	}, nil
}

func Content(dev bool, cfg *config.Config) (http.Handler, error) {
	if dev {
		return newDevHandler(cfg)
	}

	return newPrdHandler(cfg)
}
