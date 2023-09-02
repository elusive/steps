package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/elusive/steps/steps"
	"github.com/elusive/steps/util"
)

// variables to hold automated version and build 
var (
    Version string
    Build string
)

const (
	stepsFileEnvVarName string = "STEPS_FILENAME"

	// exit codes
	WORKING_PATH_ERROR    = 2
	NO_STEPS_FILE         = 3
	STEPS_FILE_LOAD_ERROR = 4
	EXECUTION_ERROR       = 5
)

// filename used if none specified
var (
	currentPath   string = ""
	stepsFileName string = ".steps"
	verbose              = false
)

func main() {

	// vars
	l := &steps.List{}

	flag.BoolVar(&verbose, "verbose", false, "Verbose output for debugging.")
	flag.Parse()

	out("Steps %s, %s", Version, Build)

	setEnvStepsFileIfFound()
	setStepsFileFromArgsIfFound()
	setCurrentPath()

	/**
     *  Initial logic for steps filename is:
     *    - default steps file
     *   - env var STEPS_FILENAME
     *   - search for *.steps
     */

	var resolved bool

	// use stepsFileName value and resolve to cwd
	stepsFileName, absErr := filepath.Abs(stepsFileName)
	resolved = (absErr == nil) && util.FileExists(stepsFileName)
	if resolved {
		out("Resolved steps filename: %s", stepsFileName)
	}

	if resolved {
		// try loading current steps file
		loadErr := l.Load(stepsFileName)
		if loadErr != nil {
			out("Error returned from loading steps file: %s", stepsFileName)
			out("%v", loadErr)
			os.Exit(STEPS_FILE_LOAD_ERROR)
		}
	} else {
		// search for steps file, if found try to load
		foundStepsFile, foundErr := FindStepsFile()
		if foundErr != nil {
			out("Error returned from find steps file: %v", foundErr)
			os.Exit(NO_STEPS_FILE)
		}
		out("Found steps file: %s", foundStepsFile)
		stepsFileName = foundStepsFile

		loadFoundErr := l.Load(stepsFileName)
		if loadFoundErr != nil {
			out("Error returned from loading found file: %v", loadFoundErr)
			os.Exit(STEPS_FILE_LOAD_ERROR)
		}
		out("Found steps file loaded")
	}

    // execute steps
	for i := 0; i < l.Count(); i++ {
        out("Executing Step %d: %v", i, l.At(i).ToString())
		exeErr := l.Execute(i)
         out("Executing STEP %d", i)
		if exeErr != nil {
			out("Error during exec of step %d: %v\n", i, exeErr)
            os.Exit(EXECUTION_ERROR)
		}

         out("Completed STEP %d", i)
	}

	out("...Exiting Steps Utility")
}

/**
 *     PRIVATE
 */

// Check if user-defined steps filename exists,
// this path should be relative to working dir.
func setEnvStepsFileIfFound() {
	envFn := os.Getenv(stepsFileEnvVarName)
	if len(envFn) > -1 {
		stepsFileName = envFn
		out("Steps filename read from env var: %s", envFn)
	}
}

func setCurrentPath() {
	var err error
	currentPath, err = os.Getwd()
	if err != nil {
		os.Exit(WORKING_PATH_ERROR)
	}
	out("Current working dir: %s", currentPath)
}

func setStepsFileFromArgsIfFound() {
	if len(flag.Args()) == 1 {
		fn := flag.Args()[0]
		stepsFileName = fn
	}
}

func out(msg string, values ...any) {
	if verbose {
		s := fmt.Sprintf(msg, values...)
		fmt.Println(s)
	}
}

func FindStepsFile() (string, error) {
	var stepsFile string

	// get current directory
	path, err := os.Getwd()
	if err != nil {
		return stepsFile, fmt.Errorf("error getting current directory: %v", err)
	}

	// search for file by extension
	for _, sf := range util.Find(path, ".steps") {
		stepsFile = sf
		break
	}

	// if nothing found return error
	if stepsFile == "" {
		return stepsFile, fmt.Errorf("no steps file found in %s", path)
	}

	// else return result
	return stepsFile, nil
}
