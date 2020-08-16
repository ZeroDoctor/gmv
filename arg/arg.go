package arg

import (
	"os"

	ppt "github.com/zerodoctor/goprettyprinter"
)

// ExecuteArgs :
type ExecuteArgs struct {
	Key     string
	Option  string
	Value   string
	Execute func(...interface{}) interface{}
}

var (
	totalOptions = make(map[string]*ExecuteArgs)
)

// Empty :
func Empty(...interface{}) interface{} {
	return false
}

// CreateEmptyArgs :
func CreateEmptyArgs() *ExecuteArgs {
	return &ExecuteArgs{
		Value:   "",
		Execute: Empty,
	}
}

// HandleArgs :
func HandleArgs(options []string, exec map[string]*ExecuteArgs) map[string]*ExecuteArgs {
	result := make(map[string]*ExecuteArgs)
	command := []string{"src", "dst"} // this could be in main.go
	ccount := 0                       // command count
	for i := 0; i < len(options); i++ {
		current := i
		arg, ok := exec[options[i]]
		if !ok {
			// could not find argument
			if options[i][0] == '-' || ccount > 1 {
				ppt.Errorln("Incorrect format", options[i])
				help, ok := exec["-h"]
				if ok {
					help.Execute([]string{})
					os.Exit(1)
				}
				ppt.Warnln("No usage available...") // could not find usage
				os.Exit(1)
			}
			arg, _ = exec[command[ccount]] // could check if obtaining command was successful
		}

		switch arg.Value {
		case "$path":
			i++
			if options[i%len(options)][0] == '-' {
				i--
				arg.Value = ""
			} else {
				arg.Value = options[i%len(options)] + "/"
			}
		case "$user":
			i++
			arg.Value = options[i%len(options)]
		case "$command":

			arg.Value = options[i]
			result[command[ccount]] = arg
			delete(totalOptions, command[ccount])
			ccount++
			continue
		}

		index := i + (current - i)
		delete(totalOptions, options[index])
		result[options[index]] = arg
	}

	for s := range totalOptions {
		result[s] = CreateEmptyArgs() // empty unused options/commands
	}

	return result
}

// CreateArg :
func CreateArg(option string, value string, execute func(...interface{}) interface{}) *ExecuteArgs {
	exec := &ExecuteArgs{
		Option:  option,
		Value:   value,
		Execute: execute,
	}

	totalOptions[option] = exec

	return exec
}
