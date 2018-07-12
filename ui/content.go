package ui

import "io"

// Current versions of some dependencies.
const (
	ReactVersion  = "0.12.2"
	JQueryVersion = "2.1.3"
)

var contents map[string]*content

// This interface abstracts the Execute method on template which is
// structurally similar in both html/template and text/template.
// We need to use an interface instead of a direct template field
// because then we will need two different fields for html template
// and text template.
type renderer interface {
	Execute(w io.Writer, data interface{}) error
}

type content struct {

	// The uri for accessing this asset
	uri string

	// The filename of the template relative to the asset directory
	template string

	// The JavaScript sources used in this HTML page
	sources []string

	// The parsed template - can be of html/template or text/template type
	tpl renderer

	// This is used to determine if a template is to be parsed as text or html
	tplType string
}

func init() {
	// The following are HTML assets that are rendered via
	// template.
	contents = map[string]*content{

		"/": &content{
			template: "index.tpl.html",
			sources: []string{
				"js/hound.js",
			},
			tplType: "html",
		},

		"/open_search.xml": &content{
			template: "open_search.tpl.xml",
			tplType:  "xml",
		},

		"/excluded_files.html": &content{
			template: "excluded_files.tpl.html",
			sources: []string{
				"js/excluded_files.js",
			},
			tplType: "html",
		},
	}
}
