package steps_test

import (
    "fmt"
    "testing"

    steps "github.com/elusive/steps"
)

const (
    ExpectedDoesNotMatch = "expected %v, got %v"
)

func TestParseStepType(t *testing.T) {
    tests := []struct{
        value string
        expected steps.StepType
        isOk bool
    }{
        {"BAT", steps.BAT, true},
        {"CMD", steps.CMD, true},
        {"xyz", nil, false},
        {"", nil, false},
        {"whatisthis?", nil, false},
    }

    for _, tc := range tests {
        actual, ok := steps.ParseStepType(tc.value)
        if actual != tc.expected || ok != tc.isOk {
            t.Fatal(ExpectedDoesNotMatch, tc.expected, actual)
        }
    }
}


func TestExecute(t *testing.T) {
    steps := steps.List{}
    record := []string{ steps.BAT, steps.StepResult.Required, ".\\test\\test.bat" }
    steps.Add(record)
    
    if len(steps) != 1 {
        t.Fatal("Invalid length of steps list.");
    }
}
