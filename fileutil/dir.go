package fileutil

import (
	"fmt"
	"strings"
	"regexp"
	"os"
	"path"
	"path/filepath"
	"io/ioutil"
)

type FileFilter struct {
	HiddenDirs bool
	HiddenFiles bool
	TemporaryFiles bool
	FollowSymLinks bool
	Blacklist []string
	Whitelist []string
}

func (ff FileFilter) IsAllowed(filePath string, isDir bool) (bool, error) {
	basename := path.Base(filePath)

	// Skip symlinks
	if !ff.FollowSymLinks {
		if isSymLink, err := IsSymLink(filePath); isSymLink {
			//fmt.Println("block symlink")
			return false, nil
		} else if err != nil {
			return false, err
		}
	}

	// Ignore hidden dirs or files
	if strings.HasPrefix(basename, ".") {
		if !ff.HiddenDirs && isDir {
			//fmt.Println("block hidden dir")
			return false, nil
		}

		if !ff.HiddenFiles && !isDir {
			//fmt.Println("block hidden file")
			return false, nil
		}
	}

	// Ignore temporary files
	if strings.HasSuffix(basename, "~") {
		if !ff.TemporaryFiles && !isDir {
			//fmt.Println("block temp file")
			return false, nil
		}
	}

	// Blacklist filter
	for _, blacklistPattern := range ff.Blacklist {
		match, err := regexp.MatchString(blacklistPattern, filePath)
		if match {
			//fmt.Println("block blacklisted file")
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}

	// Whitelist filter
	if len(ff.Whitelist) > 0 {
		for _, whitelistPattern := range ff.Whitelist {
			match, err := regexp.MatchString(whitelistPattern, filePath)
			if !match {
				//fmt.Println("block not whitelisted file")
				return false, nil
			}
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

var (
	templateFinalBasename   = "_final.template"
	templatePartialBasename = "_partial.template"
	templateCurrentBasename = "_current.template"
)

// Traverses the file tree and returns directory templates along the way,
// until it either hits defined base dir or a final directory template
func FindDirectoryTemplates(base, branch string) ([]string, error) {
	var templateFiles []string

	templateExist := func(path string) bool {
		_, err := os.Stat(path)
		if err == nil {
			return true
		}
		return false
	}

	base, err := filepath.Abs(base)
	if err != nil {
		return nil, err
	}

	branch, err = filepath.Abs(branch)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(branch, base) {
		return nil, fmt.Errorf("Base path does not exist in branch path.")
	}

	info, err := os.Stat(branch)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		branch = path.Dir(branch)
	}

	// Add current template
	currentBranch := branch
	currentTemplatePath := filepath.Join(currentBranch, templateCurrentBasename)
	if templateExist(currentTemplatePath) {
		templateFiles = append(templateFiles, currentTemplatePath)
	}

	root := filepath.Base(branch)

	// For each dir, add partial templates. Stop on final template or last dir.
	for currentBranch != root {
		finalTemplatePath := filepath.Join(currentBranch, templateFinalBasename)
		partialTemplatePath := filepath.Join(currentBranch, templatePartialBasename)

		if templateExist(finalTemplatePath) {
			templateFiles = append(templateFiles, finalTemplatePath)
			break // Do not traverse further
		} else if templateExist(partialTemplatePath) {
			templateFiles = append(templateFiles, partialTemplatePath)
		}

		if currentBranch == base {
			break
		}

		// Next parent
		currentBranch = filepath.Dir(currentBranch)
	}

	return templateFiles, nil
}

func CopyDirectory(fromDir, toDir string, filter FileFilter) error {
	fromInfo, err := os.Stat(fromDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(toDir, fromInfo.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(fromDir)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		fromFile := filepath.Join(fromDir, fileInfo.Name())

		if allowed, err := filter.IsAllowed(fromFile, fileInfo.IsDir()); err != nil {
			return err
		} else if !allowed {
			continue
		}

		if filter.FollowSymLinks {
			newFromFile, err := filepath.EvalSymlinks(fromFile)
			if err != nil {
				return err
			}

			if newFromFile != fromFile {
				fromFile = newFromFile
				fileInfo, err = os.Stat(fromFile)
				if err != nil {
					return err
				}
			}
		}

		toFile := filepath.Join(toDir, fileInfo.Name())

		if fileInfo.IsDir() {
			err = CopyDirectory(fromFile, toFile, filter)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(fromFile, toFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
