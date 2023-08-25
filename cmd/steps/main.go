package main 

import (
    "fmt"
    "os"
//    "strings"

    "github.com/elusive/steps"
)

// filename used if none specified
const stepsFileName = ".steps"


func main() {

    // vars 
    l := &steps.List{}
    fn := stepsFileName


    // HANDLE ARGS...
    
    /**
     *  Initial logic is that if no args are provided
     *  then we will load the default filename of steps.
     *  And we will begin execution of the steps found.
     *  
     *  If a single arg is passed it is treated as the
     *  filename (for now) and steps loaded from it. And
     *  then we will begin execution of those steps.
     */

    switch {

        // load steps from default file stepsFileName
        case len(os.Args) == 0 :
            l.Load(fn)

        // first arg always steps filename
        case len(os.Args) == 1: 
              
       
        default: 
            fmt.Println("Help for Steps Utility...")
    }

    // output some feedback (TODO:  remove or update this)
    fmt.Fprintln("%d steps loaded from .steps file: %s", l.Count(), fn)           

    // load from filename
    if err := l.Load(stepsFileName); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}   
