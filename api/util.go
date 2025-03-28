package api

import "html/template"

var tmpl *template.Template

func Init() {
	tmpl, _ = template.ParseGlob("templates/*.html")
}
