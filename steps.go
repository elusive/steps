package steps

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/elusive/steps/util"
)


const (
    FileNotFoundError = "The %s file was not found"
    StepDoesNotExist = "Step %d does not exist."
    ExecuteBatError = "Error executing BAT step: %s"
    ExecuteCmdError = "Error executing CMD step: %s"
    UnsupportedStepType = "Unknown or unsupported step type: %s"
)


// StepType enum
type StepType string

const (
    BAT StepType = "BAT"
    CMD StepType = "CMD"
    EXE StepType = "EXE"
)


// StepResult enum
type StepResult string

const (
    Required StepResult = "required"
    Optional StepResult = "optional"
)


// step struct
type step struct {
    Type StepType
    Result StepResult
    Text string
}

func (s *step) ToString() string {
    serialized := fmt.Sprintf("%s,%s,%s", s.Type, s.Result, s.Text)
    return serialized
}

// holds full path to steps file
var StepFile string

// list of steps
type List []step


func (l *List) Add(record []string) error {
    s := step{}
    if t, ok := ParseStepType(record[0]); ok {
        s.Type = t
    } else {
        return fmt.Errorf("Invalid value for StepType %s", record[0])
    }

    if r, ok := ParseStepResult(record[1]); ok {
        s.Result = r 
    } else {
        return fmt.Errorf("Invalid value for StepResult %s", record[1])
    }

    s.Text = record[2] 
    *l = append(*l, s)
    return nil
}

func (l *List) Execute(i int) error {
    lst := *l
    if i <= 0 || i > len(lst) {
        return fmt.Errorf(StepDoesNotExist, i)
    }

    step := lst[i-1] 

    if step.Type == BAT {
        fpath, _ := filepath.Abs(step.Text)
        _, err := exec.Command("CMD", "/C", fpath).CombinedOutput()
        if err != nil {
            return fmt.Errorf(ExecuteBatError, step.Text)
        }

        return nil
    }

    if step.Type == CMD {
        cmd := step.Text
        _, err := exec.Command("CMD", "/C", cmd).CombinedOutput()
        if err != nil {
            return fmt.Errorf(ExecuteCmdError, cmd)
        }

        return nil
    }

    return fmt.Errorf(UnsupportedStepType, step.Type)
}

func (l *List) Load(filename string) error {
    lst := *l

    f, err := os.Open("data.csv")
    if err != nil {
        return fmt.Errorf("Unable to open steps file: %v", err)
    }

    defer f.Close()

    csvReader := csv.NewReader(f)
    steps, err := csvReader.ReadAll()
    if err != nil {
        return fmt.Errorf("Unable to read steps: %v", err)
    }

    for _, stepRecord := range steps {
        lst.Add(stepRecord)
    }

    return nil
}

func (l *List) GetSteps() error {
    // get current directory
    path, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("Error getting current directory: %v", err)
    }

    for _, sf := range util.Find(path, ".steps") {
        StepFile = sf
        break
    }

    if StepFile == "" {
        return fmt.Errorf("No Steps file found in %s", path)
    }

    return nil
} 

/*
 *  PRIVATE
 */ 
var (
    stepTypeMap = map[string]StepType{
        "bat":   BAT,
        "cmd":   CMD,
        "exe":   EXE,
    }
)
func ParseStepType(str string) (StepType, bool) {
    t, ok := stepTypeMap[strings.ToLower(str)]
    return t, ok;
}

var (
    stepResultMap = map[string]StepResult{
        "required": Required,
        "optional": Optional,
    }
)
func ParseStepResult(str string) (StepResult, bool) {
    r, ok := stepResultMap[strings.ToLower(str)]
    return r, ok;
}
