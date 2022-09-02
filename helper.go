package main

import (
	"bufio"
	"io/ioutil"
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
			files = make([]string, 0, len(tempFiles))
			for _, file := range tempFiles { // [pic.jpg, other.jpeg, cat.jpg, ...]
				matched, err := regexp.MatchString(pat, file)
				if err != nil {
					ppt.Errorln("Failed to match", pat, "with", file)
					ppt.Errorln("\t", err.Error())
				}

				if matched {
					target[key] = append(target[key], file)
					files = append(files, file)
					foundCount++
				}
			}
		}

		target[key] = target[key][prev:]
	}

	ppt.Infoln("Matched", foundCount, "files...")
}

func findNewPath(file, srcFiles, dstFolder string) (string, bool) {
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

	return newPath, (dstFolder+"/" == "./"+file[:index])
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

func myWalk(src, path string, ignoreFolder []string, count *int, doRec bool) ([]string, error) {
	info, err := ioutil.ReadDir(src + "/" + path)
	if err != nil {
		return nil, err
	}

	var files []string
out:
	for _, file := range info {
		if !file.IsDir() {
			*count++
			files = append(files, path+file.Name())
		}

		if doRec && file.IsDir() {
			for _, i := range ignoreFolder {
				if strings.Contains(file.Name(), i) || strings.Contains(path, i) {
					continue out
				}
			}
			ppt.Infoln("Searching in ", path+file.Name()+"/")
			fileList, err := myWalk(src, path+file.Name()+"/", ignoreFolder, count, doRec)
			if err != nil {
				return nil, err
			}
			files = append(files, fileList...)
		}
	}

	return files, nil
}
