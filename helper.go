package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ppt "github.com/zerodoctor/goprettyprinter"
)

// GetFilesThatMatch :
func GetFilesThatMatch(files []string, target map[string][]string) {
	// okay 3 nested for loops is not the best, but it must be done
	foundCount := 0
	for key, value := range target {
		prev := len(target[key])
		for _, pat := range value { // [.jpg, .jpeg, ...]

			if len(pat) <= 0 {
				continue
			}
			for _, file := range files { // [pic.jpg, other.jpeg, cat.jpg, ...]
				matched, err := regexp.MatchString(pat, file)
				if err != nil {
					ppt.Errorln("Failed to match", pat, "with", file)
					ppt.Errorln("\t", err.Error())
				}
				if matched {
					target[key] = append(target[key], file)
					foundCount++
				}
			}
		}

		target[key] = target[key][prev:]
	}

	ppt.Infoln("Matched", foundCount, "files...")
}

func findNewPath(file, srcFiles, dstFolder string) string {
	index := 0
	newIndex := strings.LastIndex(file, "/")
	if newIndex != -1 {
		index = newIndex + 1
	}
	newPath := dstFolder + "/" + file[index:]
	newPath, err := filepath.Abs(newPath)
	if err != nil {
		ppt.Errorln("Couldn't not find relative path of", srcFiles)
		panic(err)
	}

	return newPath
}

func checkDups(path string, dupMap map[string]bool) bool {
	_, ok := dupMap[path]
	if ok {
		ppt.Warnln("Found duplicate file", path, " would you like to overwrite? (y/n):")
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			ppt.Errorln("Failed to read input")
			panic(err)
		}
		if char == 'n' {
			return false
		}
	}
	dupMap[path] = true

	return true
}
