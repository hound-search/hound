package ui

import (
	htemplate "html/template"
	ttemplate "text/template"
)

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
	// HTML template - used for serving HTML content
	htpl *htemplate.Template

	// This is only created in prd-mode, the pre-parsed template
	// Text template - currently used for /open_search.xml only
	ttpl *ttemplate.Template

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
				"js/common.js",
				"js/hound.js",
			},
			tplType: "html",
		},

		"/open_search.xml": &content{
			template: "open_search.tpl.xml",
			tplType: "xml",
		},

		"/excluded_files.html": &content{
			template: "excluded_files.tpl.html",
			sources: []string{
				"js/common.js",
				"js/excluded_files.js",
			},
			tplType: "html",
		},
	}
}
