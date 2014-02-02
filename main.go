package main

import (
	"fmt"
	"flag"
	"os"
	"io/ioutil"
	"github.com/atmoz/rigid/fileutil"
)

func main() {
	var site Site

	exit := func() {
		os.Exit(1)
	}

	flag.BoolVar(&site.Verbose, "verbose", false, "Verbose output")
	flag.StringVar(&site.SourceDirPath, "source", ".", "Source dir")
	flag.StringVar(&site.TargetDirPath, "target", "./_output", "Target dir")
	flag.Parse()

	site.TemplateRegexpPattern = `\.template$`
	site.PageRegexpPattern = `\.(?:md|markdown|html)$`

	// Create temp dir
	tempDir, err := ioutil.TempDir("", "rigid-")
	if err != nil {
		fmt.Println(err)
		exit()
	}

	// Remove temp dir when done
	defer os.RemoveAll(tempDir)
	exit = func() {
		os.RemoveAll(tempDir)
		os.Exit(1)
	}

	site.Log("Created temp dir:", tempDir)

	site.Log("\nScanning dir:", site.SourceDirPath)
	if err := site.ScanDir(site.SourceDirPath, true); err != nil {
		fmt.Println(err)
		exit()
	}

	site.Log("\nBuilding pages")
	for _, page := range site.Pages {
		if err := page.Build(tempDir); err != nil {
			fmt.Println(err)
			exit()
		}
		site.Log(" + ", page.TargetRelPath)
		//site.Log(page.SourceRelPath, "-->", page.TargetRelPath)
	}

	site.Log("\nCopy files to target dir:", site.TargetDirPath)

	if site.TargetDirPath != site.SourceDirPath {
		// Delete old target dir
		_, err = os.Stat(site.TargetDirPath)
		if err == nil {
			os.RemoveAll(site.TargetDirPath)
		}

		// Create new target dir
		if err = os.MkdirAll(site.TargetDirPath, os.FileMode(0755)); err != nil {
			fmt.Println(err)
			exit()
		}

		// Copy normal files to target dir
		filter := fileutil.FileFilter{
			HiddenDirs: false,
			HiddenFiles: true,
			TemporaryFiles: false,
			Blacklist: []string{
				site.TargetDirPath,
				site.TemplateRegexpPattern,
				site.PageRegexpPattern,
			},
		}
		err = fileutil.CopyDirectory(site.SourceDirPath, site.TargetDirPath, filter)
		if err != nil {
			fmt.Println(err)
			exit()
		}
	}

	// Copy pages from temp to target dir
	filter := fileutil.FileFilter{}
	err = fileutil.CopyDirectory(tempDir, site.TargetDirPath, filter)
	if err != nil {
		fmt.Println(err)
		exit()
	}
	

}
