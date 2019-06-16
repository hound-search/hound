package ui

import (
	"bytes"
	"errors"
	"fmt"
	html_template "html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	text_template "text/template"

	"github.com/hound-search/hound/config"
)

// An http.Handler for the dev-mode case.
type devHandler struct {
	// A simple file server for serving non-template assets
	http.Handler

	// the collection of templated assets
	content map[string]*content

	// the root asset dir
	root string

	// the config we are running on
	cfg *config.Config
}

// An http.Handler for the prd-mode case.
type prdHandler struct {
	// The collection of templated assets w/ their templates pre-parsed
	content map[string]*content

	// The config object as a json string
	cfgJson string

	// the config we are running on
	cfg *config.Config
}

func (h *devHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	// See if we have templated content for this path
	cr := h.content[p]
	if cr == nil {
		// if not, serve up files
		h.Handler.ServeHTTP(w, r)
		return
	}

	// If so, render the HTML
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if err := renderForDev(w, h.root, cr, h.cfg, r); err != nil {
		log.Panic(err)
	}
}

// Renders a templated asset in dev-mode. This simply embeds external script tags
// for the source elements.
func renderForDev(w io.Writer, root string, c *content, cfg *config.Config, r *http.Request) error {
	var err error
	// For more context, see: https://github.com/etsy/hound/issues/239
	switch c.tplType {
	case "html":
		// Use html/template to parse the html template
		c.tpl, err = html_template.ParseFiles(filepath.Join(root, c.template))
		if err != nil {
			return err
		}
	case "xml", "text":
		// Use text/template to parse the xml or text templates
		// We are using text/template here for parsing xml to keep things
		// consistent with html/template parsing.
		c.tpl, err = text_template.ParseFiles(filepath.Join(root, c.template))
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid tplType for content")
	}

	json, err := cfg.ToJsonString()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, path := range c.sources {
		fmt.Fprintf(&buf, "<script src=\"http://localhost:8080/ui/%s\"></script>", path)
	}

	return c.tpl.Execute(w, map[string]interface{}{
		"ReactVersion":  ReactVersion,
		"jQueryVersion": JQueryVersion,
		"ReposAsJson":   json,
		"Source":        html_template.HTML(buf.String()),
		"Host":          r.Host,
	})
}

// Serve an asset over HTTP. This ensures we get proper support for range
// requests and if-modified-since checks.
func serveAsset(w http.ResponseWriter, r *http.Request, name string) {
	n, err := AssetInfo(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	a, err := Asset(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, n.Name(), n.ModTime(), bytes.NewReader(a))
}

func (h *prdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	// see if we have a templated asset for this path
	ct := h.content[p]
	if ct != nil {
		// if so, render it
		if err := renderForPrd(w, ct, h.cfgJson, r); err != nil {
			log.Panic(err)
		}
		return
	}

	// otherwise, we need to find the asset in the bundled asset
	// data. Assets are relative to the asset directory, so we need
	// to remove the leading '/' in the path.
	serveAsset(w, r, p[1:])
}

// Renders a templated asset in prd-mode. This strategy will embed
// the sources directly in a script tag on the templated page.
func renderForPrd(w io.Writer, c *content, cfgJson string, r *http.Request) error {
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
		"Source":        html_template.HTML(buf.String()),
		"Host":          r.Host,
	})
}

// Used for dev-mode only. Determime the asset directory where
// we can find all our web files for direct serving.
func assetDir() string {
	_, file, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(
		filepath.Join(filepath.Dir(file), "assets"))
	if err != nil {
		log.Panic(err)
	}
	return dir
}

// Create an http.Handler for dev-mode.
func newDevHandler(cfg *config.Config) (http.Handler, error) {
	root := assetDir()
	return &devHandler{
		Handler: http.FileServer(http.Dir(root)),
		content: contents,
		root:    root,
		cfg:     cfg,
	}, nil
}

// Create an http.Handler for prd-mode.
func newPrdHandler(cfg *config.Config) (http.Handler, error) {
	for _, cnt := range contents {
		a, err := Asset(cnt.template)
		if err != nil {
			return nil, err
		}

		// For more context, see: https://github.com/etsy/hound/issues/239
		switch cnt.tplType {
		case "html":
			// Use html/template to parse the html template
			cnt.tpl, err = html_template.New(cnt.template).Parse(string(a))
			if err != nil {
				return nil, err
			}
		case "xml", "text":
			// Use text/template to parse the xml or text templates
			// We are using text/template here for parsing xml to keep things
			// consistent with html/template parsing.
			cnt.tpl, err = text_template.New(cnt.template).Parse(string(a))
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("invalid tplType for content")
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

// Create an http.Handler for serving the web assets. If dev is true,
// the http.Handler that is returned will serve assets directly our of
// the source directories making rapid web development possible. If dev
// is false, the http.Handler will serve assets out of data embedded
// in the executable.
func Content(dev bool, cfg *config.Config) (http.Handler, error) {
	if dev {
		return newDevHandler(cfg)
	}

	return newPrdHandler(cfg)
}
