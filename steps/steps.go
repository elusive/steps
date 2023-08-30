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
	FileNotFoundError   = "the %s file was not found"
	StepDoesNotExist    = "step %d does not exist"
	ExecuteBatError     = "error %v executing BAT step: %s"
	ExecuteCmdError     = "error executing CMD step: %s"
	ExecuteExeError     = "error executing EXE step: %s"
	UnsupportedStepType = "unknown or unsupported step type: %s"
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
	StepFile  string
	StepCount int
)

type List []*step

// private definition of step type
type step struct {
	Type   StepType
	Result StepResult
	Text   string
}

func (s *step) ToString() string {
	serialized := fmt.Sprintf("%s,%s,%s", s.Type, s.Result, s.Text)
	return serialized
}

/**
 * PUBLIC methods for the list
 */

// Return step record at index provided
func (l *List) At(index int) *step {
    steps := *l
    for i, step := range steps {
        if i == index {
            return step
        }
    } 

    return nil
}


// Add steps from []string
func (l *List) Add(record []string) error {
	s := step{}
	if t, ok := ParseStepType(record[0]); ok {
		s.Type = t
	} else {
		return fmt.Errorf("invalid value for StepType %s", record[0])
	}

	if r, ok := ParseStepResult(record[1]); ok {
		s.Result = r
	} else {
		return fmt.Errorf("invalid value for StepResult %s", record[1])
	}

	s.Text = record[2]
	*l = append(*l, &s)
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
			return fmt.Errorf("bat file execution not available on non-windows system")
		}

		fpath, _ := filepath.Abs(step.Text)
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
		out, err := exec.Command("cmd", "/C", step.Text).Output()
		if err == nil {
			fmt.Println(string(out))
			return nil
		}

		if step.Result == Required {
			return fmt.Errorf(ExecuteCmdError, err)
		} else {
			return nil
		}
	}

	if step.Type == EXE {
		_, exeErr := exec.Command(step.Text).Output()
		if exeErr != nil {
			return nil
		}

		if step.Result == Required {
			return fmt.Errorf(ExecuteExeError, exeErr)
		} else {
			fmt.Printf(ExecuteExeError+"\n", exeErr)
			return nil
		}
	}

	return fmt.Errorf(UnsupportedStepType, step.Type)
}

// Load list from *.steps file
func (l *List) Load(filename string) error {
	if filename == "" {
		return fmt.Errorf("empty steps file parameter")
	} else {
		StepFile = filename
	}

	f, err := os.Open(StepFile)
	if err != nil {
		return fmt.Errorf("unable to open steps file: %v", err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	steps, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("unable to read steps: %v", err)
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
