package main

import "html/template"

var indexTmpls = []string{
	"templates/index.html",
	"templates/app.html",
	"templates/add-form.html",
	"templates/list.html",
	"templates/row.html",
	"templates/aside-nav.html",
	"templates/nav-row.html",
}

var appTmpls = []string{
	"templates/app.html",
	"templates/aside-nav.html",
	"templates/nav-row.html",
	"templates/add-form.html",
	"templates/list.html",
}

var listTmpls = []string{
	"templates/list.html",
	"templates/row.html",
}

var rowTmpls = []string{
	"templates/row.html",
}

var editItemTmpls = []string{
	"templates/edit-item.html",
}

var asideNavTmpls = []string{
	"templates/aside-nav.html",
	"templates/nav-row.html",
}

var navRowTmpls = []string{
	"templates/nav-row.html",
}

func GetTemplates() map[string]*template.Template {
	tmpl := make(map[string]*template.Template)

	tmpl["index.html"] = template.Must(template.ParseFiles(indexTmpls...))
	tmpl["app.html"] = template.Must(template.ParseFiles(appTmpls...))
	tmpl["list.html"] = template.Must(template.ParseFiles(listTmpls...))
	tmpl["row.html"] = template.Must(template.ParseFiles(rowTmpls...))
	tmpl["edit-item.html"] = template.Must(template.ParseFiles(editItemTmpls...))
	tmpl["aside-nav.html"] = template.Must(template.ParseFiles(asideNavTmpls...))
	tmpl["nav-row.html"] = template.Must(template.ParseFiles(navRowTmpls...))

	return tmpl
}
