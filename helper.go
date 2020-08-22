package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ppt "github.com/zerodoctor/goprettyprinter"
)

// GetFilesThatMatch : gosh I feel this needs some work
func GetFilesThatMatch(files []string, target map[string][]string) {
	// okay 3 nested for loops is not the best, but it must be done... maybe an array would be better
	foundCount := 0
	for key, value := range target {
		prev := len(target[key])
		for _, pat := range value { // [.jpg, .jpeg, ...]
			if len(pat) <= 0 {
				continue
			}

			tempFiles := make([]string, len(files))
			copy(tempFiles, files)
			for i, file := range tempFiles { // [pic.jpg, other.jpeg, cat.jpg, ...]
				matched, err := regexp.MatchString(pat, file)
				if err != nil {
					ppt.Errorln("Failed to match", pat, "with", file)
					ppt.Errorln("\t", err.Error())
				}
				ppt.Infoln("Checking: ", file)
				if matched {
					target[key] = append(target[key], file)
					deleteElement(files, i-foundCount) // untested
					foundCount++
					ppt.Infoln("Found File: ", file)
					ppt.Infoln("List: ", files)
					ppt.Infoln("List: ", tempFiles)
				}
			}
		}

		target[key] = target[key][prev:]
	}

	ppt.Infoln("Matched", foundCount, "files...")
}

func deleteElement(arr []string, i int) []string {
	copy(arr[i:], arr[i+1:]) // Shift a[i+1:] left one index.
	arr[len(arr)-1] = ""     // Erase last element
	arr = arr[:len(arr)-1]
	return arr
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
	if ok || fileExists(path) {
		ppt.Warnln("Found duplicate file", path, "would you like to overwrite? (y/n):")
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
