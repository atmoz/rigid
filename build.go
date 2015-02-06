package main

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/atmoz/rigid/fileutil"
)

func build(site *Site) (err error) {
	// Make paths absolute
	site.SourceDirPath, err = filepath.Abs(site.SourceDirPath)
	if err != nil {
		return err
	}
	site.TargetDirPath, err = filepath.Abs(site.TargetDirPath)
	if err != nil {
		return err
	}

	// Create temp dir
	tempDir, err := ioutil.TempDir("", "rigid-")
	if err != nil {
		return err
	}

	// Remove temp dir when done
	defer os.RemoveAll(tempDir)

	site.Log("\nScanning dir:", site.SourceDirPath)
	if err := site.ScanDir(site.SourceDirPath, true); err != nil {
		return err
	}

	site.Log("\nBuilding pages")
	for _, page := range site.Pages {
		if err := page.Build(tempDir); err != nil {
			return err
		}
		site.Log(" + ", page.TargetRelPath)
		//site.Log(page.SourceRelPath, "-->", page.TargetRelPath)
	}

	site.Log("\nCopying files to target dir:", site.TargetDirPath)

	if site.TargetDirPath != site.SourceDirPath {
		// Delete old target dir
		_, err = os.Stat(site.TargetDirPath)
		if err == nil {
			os.RemoveAll(site.TargetDirPath)
		}

		// Create new target dir
		if err = os.MkdirAll(site.TargetDirPath, os.FileMode(0755)); err != nil {
			return err
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
			return err
		}
	}

	// Copy pages from temp to target dir
	filter := fileutil.FileFilter{}
	err = fileutil.CopyDirectory(tempDir, site.TargetDirPath, filter)
	if err != nil {
		return err
	}

	site.Log("\nDone.")
	return
}
