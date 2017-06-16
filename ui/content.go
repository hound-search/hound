package ui

import "html/template"

// Current versions of some dependencies.
const (
	ReactVersion  = "0.12.2"
	JQueryVersion = "2.1.3"
)

var contents map[string]*content

type content struct {

	// The uri for accessing this asset
	uri string

	// The filename of the template relative to the asset directory
	template string

	// The JavaScript sources used in this HTML page
	sources []string

	// This is only created in prd-mode, the pre-parsed template
	tpl *template.Template
}

func init() {
	// The following are HTML assets that are rendered via
	// template.
	contents = map[string]*content{

		"/": &content{
			template: "index.tpl.html",
			sources: []string{
				"js/common.js",
				"js/hound.js",
			},
		},

		"/open_search.xml": &content{
			template: "open_search.tpl.xml",
		},

		"/excluded_files.html": &content{
			template: "excluded_files.tpl.html",
			sources: []string{
				"js/common.js",
				"js/excluded_files.js",
			},
		},

		"/preferences.html": &content{
			template: "preferences.tpl.html",
			sources: []string{
				"js/common.js",
				"js/preferences.js",
			},
		},
	}
}
