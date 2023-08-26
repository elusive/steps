package main

import (
    "flag"
	"fmt"
    "path/filepath"
	"os"

	"github.com/elusive/steps/steps"
	"github.com/elusive/steps/util"
)

const (
    stepsFileEnvVarName string = "STEPS_FILENAME"

    // exit codes
    WORKING_PATH_ERROR = 2
    NO_STEPS_FILE = 3
    STEPS_FILE_LOAD_ERROR = 4
    EXECUTION_ERROR = 5
)

// filename used if none specified
var ( 
    currentPath string = ""
    stepsFileName string = ".steps"
)


func main() {

    // vars 
    l := &steps.List{}
    fn := stepsFileName

    verbose := flag.Bool("verbose", false, "Verbose output for debugging.")
    flag.Parse()

    out("Starting Steps Utility...", *verbose)


    // check if user-defined steps filename exists,
    // this path should be relative to working dir
    envFn := os.Getenv(stepsFileEnvVarName)
    if envFn != ""
        
        stepsFileName = envFn
    }


    // HANDLE ARGS...
    
    /**
     *  Initial logic is that if no args are provided
     *  then we will load the default filename of steps.
     *  And we will begin execution of the steps found.
     *    
     *  If the default file does not exist then we can
     *  perform a recursive file search for the first
     *  file we find that has the "*.steps" extension.
     * 
     *  If a single arg is passed it is treated as the
     *  filename (for now) and steps loaded from it. And
     *  then we will begin execution of those steps.
     */

    switch {

        // first arg always steps filename
        case len(flag.Args()) == 1: 
            fn = os.Args[0]
            loadErr := l.Load(fn)
            if loadErr != nil {
                os.Exit(NO_STEPS_FILE)
            }  
      
        // no args is default
        default: 
            // try to resolve default or user env steps file
            var resolved bool
            stepsFileName, absErr := filepath.Abs(stepsFileName)
            resolved = (absErr == nil) 
            out(fmt.Sprintf("Resolved steps filename: %s", stepsFileName), *verbose)

            if resolved {
                // try loading user env steps file
                loadEnvSpecifiedErr := l.Load(stepsFileName)
                if loadEnvSpecifiedErr != nil {
                    out(fmt.Sprintf("Error returned from loading steps file: %s", stepsFileName), *verbose)
                    out(fmt.Sprintf("%v", loadEnvSpecifiedErr), *verbose)
                    os.Exit(STEPS_FILE_LOAD_ERROR)
                }
            } else {
                // search for steps file, if found try to load
                foundStepsFile, foundErr := FindStepsFile()
                if foundErr != nil {
                    out(fmt.Sprintf("Error returned from find steps file: %v", foundErr), *verbose)
                    os.Exit(NO_STEPS_FILE)
                }
                out(fmt.Sprintf("Found steps file: %s", foundStepsFile), *verbose)

                loadFoundErr := l.Load(foundStepsFile)
                if loadFoundErr != nil {
                    out("Error returned from loading found file.", *verbose)                
                    os.Exit(STEPS_FILE_LOAD_ERROR)
                }
                out("Found steps file loaded", *verbose)
            }
    }

    // output some feedback (TODO:  remove or update this)
    out(fmt.Sprintf("%d steps loaded from .steps file: %s\n", l.Count(), fn), *verbose)     
    out("...Exiting Steps Utility", *verbose)
}   





/**
 *     PRIVATE
 */

func out(msg string, verbose bool) {
    if verbose {
        fmt.Println(msg)
    }
}

func FindStepsFile() (string, error) {
    var stepsFile string
	
    // get current directory
	path, err := os.Getwd()
	if err != nil {
		return stepsFile, fmt.Errorf("Error getting current directory: %v", err)
	}

    // search for file by extension
	for _, sf := range util.Find(path, ".steps") {
		stepsFile = sf
		break
	}

    // if nothing found return error
	if stepsFile == "" {
		return stepsFile, fmt.Errorf("No Steps file found in %s", path)
	}

    // else return result
	return stepsFile, nil
}


