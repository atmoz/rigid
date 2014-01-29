package main

import (
	"fmt"
	"path"
	"path/filepath"
)

type PageData struct {
	Path string
	Meta PageMeta
}

type PageTemplateData struct {
	page Page
	Page PageData
	Content string
}

func (ptd *PageTemplateData) Init() {
	ptd.Page = PageData{
		Path: ptd.page.PublicPath,
		Meta: ptd.page.Meta,
	}
}

// Execute template
func (ptd PageTemplateData) Template(relPath string, data interface{}) (string, error) {
	result, err := ptd.page.executeTemplate(relPath, data)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// Generate sitemap
func (ptd PageTemplateData) Sitemap(glob string) (html string) {
	// TODO Sorting and filtering
	// TODO Use template, allowing template to be replaced by user
	html += "\n<ul class=\"sitemap\">\n"
	for _, page := range ptd.page.site.Pages {
		class := ""
		if page.TargetRelPath == ptd.page.TargetRelPath {
			class = `class="active"`
		}
		html += fmt.Sprintf("<li><a %s href=\"%s\">%s</a></li>\n",
			class, ptd.RelPath(page.TargetRelPath), page.Meta.Title)
	}
	html += "</ul>\n"

	return
}

// Make sure path is relative to site root, from current page position
func (ptd PageTemplateData) RelPath(filePath string) string {
	filePath, err := filepath.Rel(path.Dir(ptd.page.TargetRelPath), filePath)
	if err != nil {
		filePath = err.Error()
	}
	return filePath
}
