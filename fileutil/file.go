package fileutil

import (
	"os"
	"io"
)

func CopyFile(from, to string) error {
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fromInfo, err := os.Stat(from)
	if err != nil {
		return err
	}

	toFile, err := os.OpenFile(to, os.O_CREATE | os.O_RDWR, fromInfo.Mode())
	if err != nil {
		return err
	}
	defer toFile.Close()

	_, err = io.Copy(toFile, fromFile)
	return err
}

func IsSymLink(filePath string) (bool, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return false, err
	}
	return info.Mode() & os.ModeSymlink == os.ModeSymlink, nil
}
