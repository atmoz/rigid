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

func (s *Site) Clean() {
	s.Template = *template.New("templates")
	s.Pages = make([]Page, 0)
}

func (s *Site) Log(a ...interface{}) {
	if s.Verbose {
		fmt.Println(a...)
	}
}

// Scan dir and init supported files
func (s *Site) ScanDir(dirPath string, recursive bool) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// Scan
	for _, fileinfo := range files {
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
		relFilePath, err := filepath.Rel(s.SourceDirPath, absFilePath)
		if err != nil {
			return err
		}

		// Ignore target dir
		if absFilePath == s.TargetDirPath {
			continue
		}

		if recursive && fileinfo.IsDir() {
			if err := s.ScanDir(absFilePath, true); err != nil {
				return err
			}
			continue
		}

		// Parse templates
		if match, _ := regexp.MatchString(s.TemplateRegexpPattern, fileinfo.Name()); match {
			templateContent, err := ioutil.ReadFile(absFilePath)
			if err != nil {
				return err
			}

			_, err = s.Template.New("//" + relFilePath).Parse(string(templateContent))
			if err != nil {
				return err
			}

			s.Log(" T ", relFilePath)
		}

		// Init pages
		if match, _ := regexp.MatchString(s.PageRegexpPattern, fileinfo.Name()); match {

			var page Page
			if err := page.Init(s, absFilePath); err != nil {
				return err
			}

			s.Pages = append(s.Pages, page)
			s.Log(" P ", relFilePath)
		}

	}

	return nil
}
