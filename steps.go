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

	"github.com/elusive/steps/util"
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

// step struct
type step struct {
	Type   StepType
	Result StepResult
	Text   []string
}

func (s *step) ToString() string {
	serialized := fmt.Sprintf("%s,%s,%s", s.Type, s.Result, strings.Join(s.Text, " "))
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

	s.Text = record[2:]
	*l = append(*l, s)
	return nil
}

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

func (l *List) Load(filename string) error {
	if filename == "" {
		err := GetStepFile()
		if err != nil {
			return fmt.Errorf("No steps file: %v", err)
		}

		filename = StepFile
	} else {
		StepFile = filename
	}

	f, err := os.Open(filename)
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

func GetStepFile() error {
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
		"bat": BAT,
		"cmd": CMD,
		"exe": EXE,
	}
)

func ParseStepType(str string) (StepType, bool) {
	t, ok := stepTypeMap[strings.ToLower(str)]
	return t, ok
}

var (
	stepResultMap = map[string]StepResult{
		"required": Required,
		"optional": Optional,
	}
)

func ParseStepResult(str string) (StepResult, bool) {
	r, ok := stepResultMap[strings.ToLower(str)]
	return r, ok
}
