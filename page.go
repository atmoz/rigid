package main

import (
	"strings"
	"bytes"
	"path"
	"path/filepath"
	"regexp"
	"os"
	"io/ioutil"
	"text/template"
	"github.com/atmoz/rigid/fileutil"
	"launchpad.net/goyaml"
	"github.com/russross/blackfriday"
)

type Page struct {
	site *Site

	SourceRelPath string
	TargetRelPath string

	Meta PageMeta
}

func (p *Page) Init(site *Site, filePath string) (err error) {
	p.site = site

	p.SourceRelPath, err = filepath.Rel(site.SourceDirPath, filePath)
	if err != nil {
		return err
	}

	// Determine target file path
	var targetFilePath string
	regexpExt := regexp.MustCompile(path.Ext(filePath) + "$")
	strippedName := regexpExt.ReplaceAllString(filePath, "")
	if path.Ext(strippedName) == ".html" {
		targetFilePath = strippedName
	} else {
		targetFilePath = regexpExt.ReplaceAllString(
			filePath, string(os.PathSeparator) + "index.html")
	}

	p.TargetRelPath, err = filepath.Rel(site.SourceDirPath, targetFilePath)
	if err != nil {
		return err
	}

	if err = p.Meta.ReadFromFile(filePath); err != nil {
		return err
	}

	return nil
}

func (p *Page) Build(rootDirPath string) error {
	sourceAbsFilePath := path.Join(p.site.SourceDirPath, p.SourceRelPath)
	targetAbsFilePath := path.Join(rootDirPath, p.TargetRelPath)

	sourceInfo, err := os.Stat(p.site.SourceDirPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(targetAbsFilePath), sourceInfo.Mode())
	if err != nil {
		return err
	}

	err = fileutil.CopyFile(sourceAbsFilePath, targetAbsFilePath)
	if err != nil {
		return err
	}

	// At this point:
	// The original file is copied and we can start working on the copy
	// But first, remove meta data
	_, _, err = fileutil.ReadMetaData(targetAbsFilePath, true)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(targetAbsFilePath)
	if err != nil {
		return err
	}

	// Apply markup
	if content, err = p.applyMarkup(content); err != nil {
		return err
	}

	// Render templates
	if content, err = p.renderTemplate(content); err != nil {
		return err
	}

	info, err := os.Stat(targetAbsFilePath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(targetAbsFilePath, content, info.Mode())
}

func (p *Page) applyMarkup(content []byte) ([]byte, error) {
	switch path.Ext(p.SourceRelPath) {
	case ".md", ".markdown":
		content = blackfriday.MarkdownCommon(content)
		return content, nil

	default: // No markup to process
		return content, nil
	}
}

func (p *Page) renderTemplate(content []byte) ([]byte, error) {
	sourceAbsFilePath := path.Join(p.site.SourceDirPath, p.SourceRelPath)

	data := PageTemplateData{
		page: *p,
		Page: PageData{Path: p.TargetRelPath, Meta: p.Meta},
		Content: string(content),
	}

	funcMap := data.getTemplateFuncMap()

	// Override default directory templates
	if p.Meta.Template != "" {
		content, err := p.executeTemplate(p.Meta.Template, data, funcMap)
		if err != nil {
			return nil, err
		}
		return content, nil
	}

	// TODO Scan the already cached paths in Site.Template?
	templateFiles, err := fileutil.FindDirectoryTemplates(p.site.SourceDirPath, sourceAbsFilePath)
	if err != nil {
		return nil, err
	}

	for _, templatePath := range templateFiles {
		data.Content = string(content)

		content, err = p.executeTemplate(templatePath, data, funcMap)
		if err != nil {
			return nil, err
		}
	}

	return content, nil
}

func (p *Page) executeTemplate(templatePath string, data interface{}, funcMap template.FuncMap) ([]byte, error) {
	var err error
	if !path.IsAbs(templatePath) {
		templatePath = path.Join(p.site.SourceDirPath, path.Dir(p.SourceRelPath), templatePath)
	}

	// Make template relative to source dir
	if !strings.HasPrefix(templatePath, "//") {
		templatePath, err = filepath.Rel(p.site.SourceDirPath, templatePath)
		if err != nil {
			return nil, err
		}
		templatePath = "//" + templatePath
	}

	t := p.site.Template
	t.Funcs(funcMap)

	contentBuffer := bytes.NewBufferString("")
	err = t.ExecuteTemplate(contentBuffer, templatePath, data)
	if err != nil {
		return nil, err
	}
	return contentBuffer.Bytes(), nil
}


type PageMeta struct {
	Title string
	Date string
	Tags []string
	Template string
}

func (pm *PageMeta) ReadFromFile(filePath string) error {
	meta, _, err := fileutil.ReadMetaData(filePath, false)
	if err != nil {
		return err
	}

	if err := goyaml.Unmarshal(meta, &pm); err != nil {
		return err
	}

	if pm.Title == "" {
		pm.Title = path.Base(filePath)
		// TODO Try getting title from the file
		// http://golang-examples.tumblr.com/post/47426518779/parse-html
		// http://godoc.org/code.google.com/p/go.net/html
		//if content, err := pm.page.applyMarkup(rest); err != nil {
		//	return err
		//}
	}

	return nil
}

