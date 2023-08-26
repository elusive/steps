package steps

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
    "syscall"
)

const (
    CREATE_NEW_CONSOLE = 0x10
)

const (
	FileNotFoundError   = "The %s file was not found"
	StepDoesNotExist    = "Step %d does not exist."
	ExecuteBatError     = "Error %v executing BAT step: %s"
	ExecuteCmdError     = "Error executing CMD step: %s"
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

var (
    StepFile string
    StepCount int
)

type List []step


// private definition of step type
type step struct {
	Type   StepType
	Result StepResult
	Text   []string
}

func (s *step) ToString() string {
	serialized := fmt.Sprintf("%s,%s,%s", s.Type, s.Result, strings.Join(s.Text, " "))
	return serialized
}





/**
 * PUBLIC methods for the list
 */

// Add steps from []string
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

	s.Text = record[2:]
	*l = append(*l, s)
	return nil
}

// Count of steps loaded into list
func (l *List) Count() int {
    lst := *l
    return len(lst)
}

// Execute step at index provided
func (l *List) Execute(i int) error {
	lst := *l

	if i < 0 || i >= len(lst) {
		return fmt.Errorf(StepDoesNotExist, i)
	}

	step := lst[i]

	if step.Type == BAT {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("BAT file execution not available on non-windows system.")
		}

        fpath, _ := filepath.Abs(step.Text[0])
		cmd := exec.Cmd{
            Path: fpath,
            SysProcAttr: &syscall.SysProcAttr{
                CreationFlags:    CREATE_NEW_CONSOLE,
                NoInheritHandles: true,
            },
        }

		if err := cmd.Run(); err != nil {
			return fmt.Errorf(ExecuteBatError, err, step.Text)
		}

		return nil
	}

	if step.Type == CMD {
	    cmd := exec.Command("cmd", step.Text[:]...)
		cmd.Stdout = os.Stdout
        if err := cmd.Run(); err != nil {
			return fmt.Errorf(ExecuteCmdError, err)
		}

		return nil
	}

    if step.Type == EXE {
        return fmt.Errorf(UnsupportedStepType, step.Type)
    }

	return fmt.Errorf(UnsupportedStepType, step.Type)
}

// Load list from *.steps file
func (l *List) Load(filename string) error {
	if filename == "" {
        return fmt.Errorf("Empty steps file parameter.")
	} else {
		StepFile = filename
	}

	f, err := os.Open(StepFile)
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
		l.Add(stepRecord)
	}

	return nil
}



/**
 *   PUBLIC methods (not part of list)
 */

// GetStepFile returns the first steps file found. 


func ParseStepType(str string) (StepType, bool) {
	t, ok := stepTypeMap[strings.ToLower(str)]
	return t, ok
}

func ParseStepResult(str string) (StepResult, bool) {
	r, ok := stepResultMap[strings.ToLower(str)]
	return r, ok
}



/**
 * PRIVATE 
 */
var (
	stepTypeMap = map[string]StepType{
		"bat": BAT,
		"cmd": CMD,
		"exe": EXE,
	}
)


var (
	stepResultMap = map[string]StepResult{
		"required": Required,
		"optional": Optional,
	}
)

