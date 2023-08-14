package steps

import (
    "fmt"
    "os/exec"
    "path/filepath"
    "strings"
)


const (
    FileNotFoundError = "The %s file was not found"
    StepDoesNotExist = "Step %d does not exist."
    ExecuteBatError = "Error executing BAT step: %s"
    ExecuteCmdError = "Error executing CMD step: %s"
    UnsupportedStepType = "Unknown or unsupported step type: %s"
)



type StepType string

const (
    BAT StepType = "BAT"
    CMD StepType = "CMD"
    EXE StepType = "EXE"
)



type StepResult string

const (
    Required StepResult = "required"
    Optional StepResult = "optional"
)



type step struct {
    Type StepType
    Result StepResult
    Text string
}

func (s *step) ToString() string {
    serialized := fmt.Sprintf("%s,%s,%s", s.Type, s.Result, s.Text)
    return serialized
}



type List []step


func (l *List) Add(record []string) error {
    lst := *l
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
    lst = append(lst, s)
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
