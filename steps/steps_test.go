package steps

import (
	"bufio"
	"os"
	"runtime"
	"testing"

	"github.com/elusive/steps/util"
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

func TestCount(t *testing.T) {
	lst := List{}
	expected := 1

	const cmdText string = "&quot;Install.bat&quot;"
	stepValues := []string{string(BAT), string(Required), cmdText}
	lst.Add(stepValues)

	// verify we have added a single step
	if len(lst) != 1 {
		t.Fatalf("Steps list count %d not expected.", len(lst))
	}

	// now verify using Count() method
	actual := lst.Count()
	if actual != expected {
		t.Fatalf("Expected count of %d, got %d", expected, actual)
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

	// arrange
	l1 := List{}
	tf, err := os.CreateTemp("./", "tmp*.steps")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
    defer tf.Close()

	// act
	l1.Load(tf.Name())

	// assert
	if StepFile == "" {
		t.Fatal("Step file not set")
	}

	// cleanup
	t.Cleanup(func() {
		os.Remove(tf.Name())
	})
}

func TestLoad(t *testing.T) {

	// arrange
	lst := List{}
	tf, err := os.CreateTemp("./", "tmp*.steps")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
    defer tf.Close()

	// act
	w := bufio.NewWriter(tf)
	for _, rec := range testStepRecords {
		w.WriteString(rec + "\n")
	}
	w.Flush()
	err = lst.Load(tf.Name())
	if err != nil {
		t.Fatal(err)
	}

	// assert
	if len(lst) != 3 {
		t.Fatalf("Steps not loaded expected 3 got: %d", len(lst))
	}

	// cleanup
	t.Cleanup(func() {
		os.Remove(tf.Name())
	})
}

func TestExecuteCmd(t *testing.T) {

	// arrange
	steps := List{}
	var commandText string
	if runtime.GOOS == "windows" {
		commandText = "pwd"
	} else {
		commandText = "ls"
	}
	record := []string{string(CMD), string(Required), commandText}

	// act
	steps.Add(record)

	// assert
	if len(steps) != 1 {
		t.Fatal("Invalid length of steps list.")
	}

	// act
	for i := range steps {
		err := steps.Execute(i)

		// assert
		if err != nil {
			t.Fatalf("Error occured during execution of step %d: %v", i, err)
		}
	}
}

func TestExecuteBat(t *testing.T) {

	// arrange
	steps := List{}
	var record []string
	if runtime.GOOS == "windows" {
		record = []string{string(BAT), string(Required), "..\\test\\test.bat"}
	} else {
		record = []string{string(BAT), string(Required), "../test/test.sh"}
	}

	// act
	steps.Add(record)

	// assert
	if len(steps) != 1 {
		t.Fatal("Invalid length of steps list.")
	}

	// act
	for i := range steps {
		err := steps.Execute(i)

		// assert
		if err != nil {
			t.Fatalf("Error ocurred during execution of step %d: %v", i, err)
		}
	}
}

func TestExecuteExe(t *testing.T) {
	// arrange
	const processName = "notepad.exe"
	steps := List{}
	record := []string{string(EXE), string(Optional), processName}

	// act
	steps.Add(record)
	for i := range steps {
		util.KillProcessAsync(processName)
		err := steps.Execute(i)

		// assert (if on windows, linux does not have notepad)
		if err != nil && runtime.GOOS == "windows" {
			t.Fatalf("Error occured during exec of step %d: %v", i, err)
		}
	}
}
