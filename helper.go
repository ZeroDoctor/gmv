package main

import (
	"regexp"

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
