package steps

import (
    "testing"
)

const (
    ExpectedDoesNotMatch = "expected %v, got %v"
)

func TestAdd(t *testing.T) {
    //lst := List{}

    
}

func TestParseStepType(t *testing.T) {
    tests := []struct{
        value string
        expected StepType
        isOk bool
    }{
        {"BAT", BAT, true},
        {"CMD", CMD, true},
//        {"xyz", nil, false},
//        {"", nil, false},
//        {"whatisthis?", nil, false},
    }

    for _, tc := range tests {
        actual, ok := ParseStepType(tc.value)
        if actual != tc.expected || ok != tc.isOk {
            t.Errorf(ExpectedDoesNotMatch, tc.expected, actual)
        }
    }
}


func TestExecute(t *testing.T) {
    steps := List{}
    record := []string{ string(BAT), string(Required), ".\\test\\test.bat" }
    steps.Add(record)
    
    if len(steps) != 1 {
        t.Fatal("Invalid length of steps list.");
    }
}
