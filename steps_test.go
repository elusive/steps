package steps

import (
	"bufio"
	"os"
    "runtime"
	"testing"
)

const (
	ExpectedDoesNotMatch = "expected %v, got %v"
)

var testStepRecords []string = []string{
	"\"BAT\",\"required\",\"test.bat\"",
	"\"CMD\",\"optional\",\"cd ../ && echo pwd\"",
	"\"EXE\",\"required\",\"notepad.exe\"",
}

func TestAdd(t *testing.T) {
	lst := List{}
	const cmdText string = "&quot;Install.bat&quot;"
	stepValues := []string{string(BAT), string(Required), cmdText}

	lst.Add(stepValues)

	if len(lst) != 1 {
		t.Fatalf("Steps list count %d not expected.", len(lst))
	}

	if lst[0].Type != BAT {
		t.Fatalf("Added step type expected %v got %v", BAT, lst[0].Type)
	}

	if lst[0].Result != Required {
		t.Fatalf("Added step type expected %v got %v", Required, lst[0].Result)
	}

	if lst[0].Text != cmdText {
		t.Fatalf("Added step type expected %v got %v", cmdText, lst[0].Text)
	}
}

func TestParseStepType(t *testing.T) {
	tests := []struct {
		value    string
		expected StepType
		isOk     bool
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

func TestGetStepFile(t *testing.T) {
	l1 := List{}

	tf, err := os.CreateTemp("./", "tmp*.steps")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	l1.Load(tf.Name())

	if StepFile == "" {
		t.Fatal("Step file not set")
	}
}

func TestLoad(t *testing.T) {
	lst := List{}

	tf, err := os.CreateTemp("./", "tmp*.steps")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	// add steps to temp file
	w := bufio.NewWriter(tf)
	for _, rec := range testStepRecords {
		w.WriteString(rec + "\n")
	}
	w.Flush()

	err = lst.Load(tf.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(lst) != 3 {
		t.Fatalf("Steps not loaded expected 3 got: %d", len(lst))
	}
}

func TestExecuteCmd(t *testing.T) {
	steps := List{}
    var commandText string

    if runtime.GOOS == "windows" {
        commandText = "echo hello world"
    } else { 
        commandText = "ls ./"
    }

    record := []string{string(CMD), string(Required), commandText}
	steps.Add(record)

	if len(steps) != 1 {
		t.Fatal("Invalid length of steps list.")
	}

	for i := range steps {
		err := steps.Execute(i)
        if err != nil {
            t.Fatalf("Error occured during execution of step %d: %v", i, err)
        }
	}
}



