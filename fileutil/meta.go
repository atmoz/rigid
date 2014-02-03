package fileutil

import (
	"regexp"
	"os"
	"bufio"
	"io/ioutil"
)

// Read meta data from file (if any)
// If removeMeta is true, meta data is also removed from the file
// Returns meta data and the rest of the file content
func ReadMetaData(filePath string, removeMeta bool) ([]byte, []byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var lineNum int
	var meta []byte
	var rest []byte
	var isMeta bool

	addLine := func(line []byte) {
		if isMeta {
			meta = append(meta, line...)
		} else {
			rest = append(rest, line...)
		}
	}

	reader := bufio.NewReader(file)
	for {
		lineNum++
		line, readerErr := reader.ReadBytes('\n')

		// If reader returns error, it's probably EOF with missing \n
		if readerErr != nil {
			// Append what's left and stop
			addLine(line)
			break
		}

		// Begin meta data on first line
		if lineNum == 1 {
			if match, _ := regexp.MatchString(`^-{3,}\s*$`, string(line)); match {
				isMeta = true
				continue
			}
		}

		// End meta data
		if isMeta {
			if match, _ := regexp.MatchString(`^-{3,}\s*$`, string(line)); match {
				isMeta = false

				if removeMeta {
					continue
				} else {
					break
				}
			}
		}

		addLine(line)
	}

	file.Close()

	if removeMeta {
		info, err := os.Stat(filePath)
		if err != nil {
			return nil, nil, err
		}

		err = ioutil.WriteFile(filePath, rest, info.Mode())
		if err != nil {
			return nil, nil, err
		}
	}

	return meta, rest, nil
}
