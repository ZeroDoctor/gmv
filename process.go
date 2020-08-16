package main

// ProcessCommands :
func ProcessCommands(src *string, dst *string) {
	if *dst == "" {
		*dst = "."
	}

	if *src == "" {
		*src = "."
	}

	*src = *src + "/"
	*dst = *dst + "/"
}

// ProcessTarget :
func ProcessTarget(newFileFolderMap, fileFolderMap interface{}) interface{} {
	var fileMap interface{}
	_, ok := newFileFolderMap.(map[string][]string) // checking if -t option was used
	if ok {
		fileMap = newFileFolderMap // if -t was used then update fileFolderMap
	} else {
		fileMap = Target("", fileFolderMap) // returns map[string][]string
	}

	return fileMap
}

// ProcessFolder :
func ProcessFolder(newFileFolderMap, fileFolderMap, doRec interface{}, srcFiles, dstFolder string) interface{} {
	var fileMap interface{}
	_, ok := newFileFolderMap.(map[string][]string)
	if ok {
		fileMap = newFileFolderMap
	} else {
		fileMap = Folder(doRec, srcFiles, dstFolder, fileFolderMap)
	}

	return fileMap
}
