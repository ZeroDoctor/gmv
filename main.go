/**
 * @file main.go
 *
 * @brief moves or rename files to a specificed folder
 *
 * @detail it uses a config file or command line args to
 *	move files to folders. Also, It uses regular experssion
 *	to match files to the specific folder(s). Available
 *	options can be viewed doing 'gmv -h'
 *
 * @brief design
 *		1. maps used options/commands to their desired function
 *		2. process pattern in config into a map against their desire folder output
 *		3. updates pattern if used -t options is used
 *		4. find desired files and update folder and file locations
 *		5. start creating folders if -g is used and move files to folders
 *
 * @limits
 *	1. when using '.' to match files use '\'
 *		i.e. \.docx
 *
 *	2. the program won't check if options are
 *		used correctly. This means that sometimes
 *		nothing will happen
 *
 * @author Daniel Castro
 *
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zerodoctor/gmv/ini"

	"github.com/zerodoctor/gmv/arg"

	ppt "github.com/zerodoctor/goprettyprinter"
)

func main() {

	arr := os.Args[1:] // array of string arguments
	if len(arr) <= 0 {
		arr = append(arr, "-g")
		arr = append(arr, "-c")
	}

	// setup optional arguments by "mapping" arguments to functions
	config := arg.CreateArg("-c", "$path", Config)
	folder := arg.CreateArg("-f", "$path", Folder)
	target := arg.CreateArg("-t", "$user", Target)
	help := arg.CreateArg("-h", "", Help)
	gen := arg.CreateArg("-g", "", Generate)
	rec := arg.CreateArg("-r", "", Recursive)

	src := arg.CreateArg("src", "$command", SrcFiles)
	dst := arg.CreateArg("dst", "$command", DstFolder)

	// combind args into a map
	argMap := make(map[string]*arg.ExecuteArgs)
	argMap["-c"] = config
	argMap["-f"] = folder
	argMap["-t"] = target
	argMap["-h"] = help
	argMap["-g"] = gen
	argMap["-r"] = rec
	argMap["src"] = src
	argMap["dst"] = dst

	usedOptions := arg.HandleArgs(arr, argMap) // returns a map of actual user used options

	usedOptions["-h"].Execute()          // returns void
	doGen := usedOptions["-g"].Execute() // returns bool
	doRec := usedOptions["-r"].Execute() // returns bool

	srcFiles := usedOptions["src"].Value
	dstFolder := usedOptions["dst"].Value
	ProcessCommands(&srcFiles, &dstFolder)

	configFile := usedOptions["-c"].Value                  // path to config file
	fileFolderMap := usedOptions["-c"].Execute(configFile) // returns map[string]string

	targetPattern := usedOptions["-t"].Value                                    // a regular expression to get matched file names
	newFileFolderMap := usedOptions["-t"].Execute(targetPattern, fileFolderMap) // returns map[string][]string
	fileFolderMap = ProcessTarget(newFileFolderMap, fileFolderMap)

	targetFolder := usedOptions["-f"].Value                                                              // path to folder/files  to be processed
	newFileFolderMap = usedOptions["-f"].Execute(doRec, srcFiles, dstFolder+targetFolder, fileFolderMap) // returns []string
	fileFolderMap = ProcessFolder(newFileFolderMap, fileFolderMap, doRec, srcFiles, dstFolder)

	_, ok := fileFolderMap.(bool)
	if ok {
		ppt.Warnln("Didn't create map...")
	}
	MoveFiles(doGen.(bool), srcFiles, fileFolderMap) // actually move the files
	ppt.Infoln("Done!")
}

// Config : parses config file to map target files to destination folders
func Config(inter ...interface{}) interface{} {
	configFile := inter[0].(string)

	// check if config file exist
	if configFile == "" {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}

		// find first *.ini file
		for _, f := range files {
			if strings.Contains(f.Name(), ".ini") {
				configFile = f.Name()
				ppt.Infoln("Found first config file named:", f.Name(), "...")
			}
		}
	}

	// parse config
	targetMap, err := ini.ParseFile(configFile)

	if err != nil {
		os.Exit(2)
	}

	// return map of directories and files
	return targetMap
}

// Target : appends or creates pattern to mapped folders and files
func Target(inter ...interface{}) interface{} {
	targetPattern := inter[0].(string)
	fileFolderMap, ok := inter[1].(map[string]string)

	result := make(map[string][]string)

	if ok {
		for key := range fileFolderMap {
			arrStr := strings.Split(fileFolderMap[key], ",")
			arrStr = arrStr[:len(arrStr)-1]

			result[key] = append(result[key], arrStr...)
			if len(targetPattern) > 0 {
				result[key] = append(result[key], "$"+targetPattern)
			}
		}

		return result
	}

	result["$TEMP"] = []string{targetPattern}

	return result
}

// Folder : creates folders and finds matching files
func Folder(inter ...interface{}) interface{} {
	doRec := inter[0].(bool)
	srcFiles := inter[1].(string)
	dstFolder := inter[2].(string)
	fileFolderMap, ok := inter[3].(map[string][]string)

	if ok { // might not need this conditional statement
		tempMap := make(map[string][]string)
		for key := range fileFolderMap {
			path := dstFolder + key
			if key == "$TEMP" {
				path = dstFolder
			}
			tempMap[path] = fileFolderMap[key]
		}
		fileFolderMap = tempMap
	}

	srcFiles, err := filepath.Abs(srcFiles)
	if err != nil {
		ppt.Errorln("Couldn't not find relative path of", srcFiles)
		panic(err)
	}

	count := 0
	var files []string // get current list of files
	err = filepath.Walk(srcFiles, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			ppt.Errorln("Couldn't find directory named:", srcFiles)
			os.Exit(2)
		}

		if info.IsDir() {
			return nil
		}

		path = path[len(srcFiles)+1:]

		if !doRec && strings.Contains(path, "/") {
			return nil
		}

		count++
		files = append(files, path)
		return nil
	})

	ppt.Infoln("Found", count, "files...")

	if err != nil {
		ppt.Errorln("Couldn't access folder")
		panic(err)
	}

	GetFilesThatMatch(files, fileFolderMap) // actually matches files using the desired pattern

	return fileFolderMap
}

// Help : prints usage to console
func Help(inter ...interface{}) interface{} {
	var result interface{}

	fmt.Println(`
usage: gmv [options] [src dst]

	src - path to source files
	dst - moves files to specified destination folder

options:
	-c		uses a config file to move file to a folder, if left blank 
			it will assume config is in the same directory
	
	-f		moves files to a specific folder
	-t		gets a list of files that matches user specific string *untested*
	-h		prints out usages of this program
	-g		generates folders if they do not exisit
	-r		recursively find files in subdirectories
	`)

	os.Exit(0)

	return result
}

// Generate : if used then it will generate needed folders
func Generate(inter ...interface{}) interface{} {
	return true
}

// Recursive : if used then will match files from subdirectories
func Recursive(inter ...interface{}) interface{} {
	return true
}

// SrcFiles : not used but needed to pass into CreateArgs
func SrcFiles(inter ...interface{}) interface{} {
	var result interface{}
	return result
}

// DstFolder : not used but needed to pass into CreateArgs
func DstFolder(inter ...interface{}) interface{} {
	var result interface{}
	return result
}

// MoveFiles : move mapped files to folder
func MoveFiles(doGen bool, srcFiles string, folderFiles interface{}) {

	ffmap, _ := folderFiles.(map[string][]string)
	// generate folders if -g was used TODO: move this else where
	for key, value := range ffmap {
		if doGen && len(value) > 0 {
			ppt.Infoln("Creating folder", key)
			err := os.MkdirAll(key, os.ModePerm)
			if err != nil {
				ppt.Errorln("Couldn't create folder:", key)
				panic(err)
			}
		}

		for _, file := range value {
			newPath := key + "/" + file
			err := os.Rename(srcFiles+"/"+file, newPath)
			if err != nil {
				ppt.Errorln("Couldn't move file:", file, "to:", newPath)
				ppt.Errorln("\t", err.Error())
			} else {
				ppt.Infoln("moved:", file, "to", newPath, "...")
			}
		}
	}
}
