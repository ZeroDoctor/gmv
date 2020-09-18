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
	ppt.Init()

	// setup optional arguments by "mapping" arguments to their function
	config := arg.CreateArg("-c", "$path", Config)
	folder := arg.CreateArg("-f", "$path", Folder)
	target := arg.CreateArg("-t", "$user", Target)
	help := arg.CreateArg("-h", "", Help)
	version := arg.CreateArg("-v", "", Version)
	gen := arg.CreateArg("-g", "", Generate)
	rec := arg.CreateArg("-r", "", Recursive)
	ignore := arg.CreateArg("-i", "$multipath", Ignore)

	src := arg.CreateArg("src", "$command", SrcFiles)
	dst := arg.CreateArg("dst", "$command", DstFolder)

	// combind args into a map
	argMap := make(map[string]*arg.ExecuteArgs)
	argMap["-c"] = config
	argMap["-f"] = folder
	argMap["-t"] = target
	argMap["-h"] = help
	argMap["-v"] = version
	argMap["-g"] = gen
	argMap["-r"] = rec
	argMap["-i"] = ignore
	argMap["src"] = src
	argMap["dst"] = dst

	usedOptions := arg.HandleArgs(arr, argMap) // returns a map of actual user used options

	usedOptions["-h"].Execute() // returns void from Help()
	usedOptions["-v"].Execute()
	doGen := usedOptions["-g"].Execute() // returns bool from Generate()
	doRec := usedOptions["-r"].Execute() // returns bool from Recursive()

	srcFiles := usedOptions["src"].Value
	dstFolder := usedOptions["dst"].Value
	ProcessCommands(&srcFiles, &dstFolder)

	configFile := usedOptions["-c"].Value                  // path to config file
	fileFolderMap := usedOptions["-c"].Execute(configFile) // returns map[string]string from Config()

	targetPattern := usedOptions["-t"].Value                                    // a regular expression to get matched file names
	newFileFolderMap := usedOptions["-t"].Execute(targetPattern, fileFolderMap) // returns map[string][]string from Target()
	fileFolderMap = ProcessTarget(newFileFolderMap, fileFolderMap)

	ignoreFolder := usedOptions["-i"].ValueArr
	targetFolder := usedOptions["-f"].Value                                                    // path to folder/files  to be processed
	newFileFolderMap = usedOptions["-f"].Execute(doRec, srcFiles, targetFolder, fileFolderMap) // returns []string from Folder()
	fileFolderMap = ProcessFolder(newFileFolderMap, fileFolderMap, doRec, srcFiles, dstFolder, ignoreFolder)

	fmt.Println()
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
	ignoreFolder := inter[4].([]string)

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
	ppt.Infoln("Searching in", srcFiles, "...")
	files, err := myWalk(srcFiles, "", ignoreFolder, &count, doRec)
	if err != nil {
		ppt.Errorln("Couldn't access folder", err.Error())
		panic(err)
	}

	ppt.Infoln("Found", count, "files...")

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

// Version :
func Version(inter ...interface{}) interface{} {
	var result interface{}

	fmt.Println(`
gmv version: v1.1.0
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

// Ignore :
func Ignore(inter ...interface{}) interface{} {
	var result interface{}
	return result
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

	dupMap := make(map[string]bool)
	ffmap, _ := folderFiles.(map[string][]string)

	for key, value := range ffmap { // map[key = folder it needs to be in] value = file that needs to be moved
		if doGen && len(value) > 0 { // generate folders if -g was used
			ppt.Infoln("Creating folder", key, "...")
			err := os.MkdirAll(key, os.ModePerm)
			if err != nil {
				ppt.Errorln("Couldn't create folder:", key)
				panic(err)
			}
		}

		for _, file := range value {
			newPath, same := findNewPath(file, srcFiles, key)
			if same {
				continue // its already in the target folder
			}
			okay := checkDups(newPath, dupMap)
			if okay {
				err := os.Rename(srcFiles+"/"+file, newPath)
				if err != nil {
					ppt.Errorln("Couldn't move file:", file, "to:", key)
					ppt.Errorln("\t", err.Error())
				}

				ppt.Infoln("moved:", file, "to", key, "...")
			}
		}
	}
}
