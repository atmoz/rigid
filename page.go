package main

import (
	"strings"
	"bytes"
	"path"
	"path/filepath"
	"regexp"
	"os"
	"io/ioutil"
	"github.com/atmoz/rigid/fileutil"
	"launchpad.net/goyaml"
	"github.com/russross/blackfriday"
)

type Page struct {
	site *Site

	SourceRelPath string
	TargetRelPath string
	PublicPath string

	Meta PageMeta
}

func (p *Page) Init(site *Site, filePath string) (err error) {
	p.site = site

	p.SourceRelPath, err = filepath.Rel(site.SourceDirPath, filePath)
	if err != nil {
		return err
	}

	var prettyPath bool

	// Determine target file path
	var targetFilePath string
	regexpExt := regexp.MustCompile(path.Ext(filePath) + "$")
	strippedName := regexpExt.ReplaceAllString(filePath, "")
	if path.Ext(strippedName) == ".html" {
		prettyPath = false
		targetFilePath = strippedName
	} else {
		prettyPath = true
		targetFilePath = regexpExt.ReplaceAllString(
			filePath, string(os.PathSeparator) + "index.html")
	}

	p.TargetRelPath, err = filepath.Rel(site.SourceDirPath, targetFilePath)
	if err != nil {
		return err
	}

	if prettyPath {
		p.PublicPath = path.Dir(p.TargetRelPath)
	} else {
		p.PublicPath = p.TargetRelPath
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

	// Render markup
	if content, err = p.renderMarkup(content); err != nil {
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

func (p *Page) renderMarkup(content []byte) ([]byte, error) {
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
		Content: string(content),
	}
	data.Init()

	// Override default directory templates
	if p.Meta.Template != "" {
		content, err := p.executeTemplate(p.Meta.Template, data)
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

		content, err = p.executeTemplate(templatePath, data)
		if err != nil {
			return nil, err
		}
	}

	return content, nil
}

func (p *Page) executeTemplate(templatePath string, data interface{}) ([]byte, error) {
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

	contentBuffer := bytes.NewBufferString("")
	err = p.site.Template.ExecuteTemplate(contentBuffer, templatePath, data)
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
		for {
			if path.Ext(pm.Title) == "" {
				break;
			}
			pm.Title = strings.TrimSuffix(pm.Title, path.Ext(pm.Title))
		}

		pm.Title = strings.Replace(pm.Title, "-", " ", -1)
		pm.Title = strings.Replace(pm.Title, "_", " ", -1)
		pm.Title = strings.Title(pm.Title)
	}

	return nil
}

