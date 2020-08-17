package ini

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	ppt "github.com/zerodoctor/goprettyprinter"
)

// ParseFile :
func ParseFile(path string) (map[string]string, error) {

	result := make(map[string]string)

	path, err := filepath.Abs(path)
	if err != nil {
		ppt.Errorln("Couldn't not find relative path of", path)
		panic(err)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		ppt.Errorln("Couldn't not read file", path)
		return nil, err
	}

	newlineChar := "\n"
	if runtime.GOOS == "windows" {
		newlineChar = "\r\n"
	}

	lines := strings.Split(string(data), newlineChar)

	for i := 0; i < len(lines); i++ {
		j := i + 1
		for j < len(lines) && len(lines[j]) > 0 && lines[j][0] != '[' {
			if len(strings.Trim(lines[j], " ")) > 0 {
				key := strings.Trim(lines[i], "[]") // removes '[]' around the directories
				result[key] += lines[j] + ","
			}

			j++
		}
		i = j - 1
	}

	ppt.Infoln("Found", path, "...")

	return result, nil
}
