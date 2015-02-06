package main

import (
	"fmt"
	"strings"
	"regexp"
	"io/ioutil"
	"path/filepath"
	"text/template"
)

type Site struct {
	Verbose bool

	SourceDirPath string
	TargetDirPath string

	TemplateRegexpPattern string
	PageRegexpPattern string

	Template template.Template
	Pages []Page
}

func (site *Site) Clean() {
	site.Template = *template.New("templates")
	site.Pages = make([]Page, 0)
}

func (site *Site) Log(a ...interface{}) {
	if site.Verbose {
		fmt.Println(a...)
	}
}

// Scan dir and init supported files
func (site *Site) ScanDir(dirPath string, recursive bool) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// Scan
	for _, fileinfo := range files {
		// @todo Use ignore list from config

		// Ignore hidden dirs (files are okay, like .htaccess)
		if fileinfo.IsDir() && strings.HasPrefix(fileinfo.Name(), ".") {
			continue
		}

		// Ignore temporary files
		if strings.HasSuffix(fileinfo.Name(), "~") {
			continue
		}

		// Get absolute path and follow symlink
		absFilePath := filepath.Join(dirPath, fileinfo.Name())
		absFilePath, err := filepath.EvalSymlinks(absFilePath)
		if err != nil {
			return err
		}

		// Get path relative to source dir
		relFilePath, err := filepath.Rel(site.SourceDirPath, absFilePath)
		if err != nil {
			return err
		}

		// Ignore target dir
		if absFilePath == site.TargetDirPath {
			continue
		}

		if recursive && fileinfo.IsDir() {
			if err := site.ScanDir(absFilePath, true); err != nil {
				return err
			}
			continue
		}

		// Parse templates
		if match, _ := regexp.MatchString(site.TemplateRegexpPattern, fileinfo.Name()); match {
			templateContent, err := ioutil.ReadFile(absFilePath)
			if err != nil {
				return err
			}

			_, err = site.Template.New("//" + relFilePath).Parse(string(templateContent))
			if err != nil {
				return err
			}

			site.Log(" T ", relFilePath)
		}

		// Init pages
		if match, _ := regexp.MatchString(site.PageRegexpPattern, fileinfo.Name()); match {

			var page Page
			if err := page.Init(site, absFilePath); err != nil {
				return err
			}

			site.Pages = append(site.Pages, page)
			site.Log(" P ", relFilePath)
		}

	}

	return nil
}
