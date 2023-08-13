package steps

import (
    "encoding/csv"
    "errors"
    "fmt"
    "os"
    "strings"
)


const (
    FileNotFoundError = "The %s file was not found"
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
    serialized, err := fmt.Sprintf("%s,%s,%s", s.Type, s.Result, s.Text)
    if err != nil {
        return nil
    }

}


