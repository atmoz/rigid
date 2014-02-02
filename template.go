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
func (ptd PageTemplateData) Sitemap(pattern string) (html string) {
	// TODO Sorting and filtering
	// TODO Use template, allowing template to be replaced by user
	html += "\n<ul class=\"sitemap\">\n"
	for _, page := range ptd.page.site.Pages {
		if pattern != "" {
			if match, err := path.Match(pattern, page.PublicPath); err != nil {
				panic(err)
			} else if !match {
				continue
			}
		}

		class := ""
		if page.PublicPath == ptd.page.PublicPath {
			class = `class="active"`
		}
		html += fmt.Sprintf("<li><a %s href=\"%s\">%s</a></li>\n",
			class, ptd.RelPath(page.PublicPath), page.Meta.Title)
	}
	html += "</ul>\n"

	return
}

// Generate sitemap
func (ptd PageTemplateData) TaggedPages(pattern string) (pages []Page) {
	// TODO Sorting and filtering
	forPages:
	for _, page := range ptd.page.site.Pages {
		if pattern != "" {
			for _, tag := range page.Meta.Tags {
				if match, err := path.Match(pattern, tag); err != nil {
					panic(err)
				} else if match {
					pages = append(pages, page)
					continue forPages
				}
			}
		}
	}

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
