package steps

import (
    "encoding/csv"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "strings"
)


const (
    FileNotFoundError = "The %s file was not found"
    StepDoesNotExist = "Step %d does not exist."
    ExecuteBatError = "Error executing BAT step: %s"
    ExecuteCmdError = "Error executing CMD step: %s"
    UnsupportedStepType = "Unknown or unsupported step type: %s"
)



type StepType *string

const (
    BAT StepType = "BAT"
    CMD StepType = "CMD"
    EXE StepType = "EXE"
)



type StepResult *string

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

func (s *step) Execute() error {
   }


var (
    capabilitiesMap = map[string]Capability{
        "read":   Read,
        "create": Create,
        "update": Update,
        "delete": Delete,
        "list":   List,
    }
)
func ParseString(str string) (Capability, bool) {
    c, ok := capabilitiesMap[strings.ToLower(str)]
    return c, ok
}
type List []step

func (l *List) Add(record []string) error {
    lst := *l
    s := step{
        Type: ParseStepType(record[0]),
        Result: ParseStepResult(record[1]),

    } 
}

func (l *List) Execute(i int) error {
    lst := *l
    if i <= 0 || i > len(lst) {
        return fmt.Errorf(StepDoesNotExist, i)
    }

    step := lst[i-1] 

    if step.Type == BAT {
        path := filepath.Abs(step.Text)
        _, err := exec.Command("CMD", "/C", path).CombinedOutput()
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
