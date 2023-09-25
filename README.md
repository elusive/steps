# Steps

Windows command runner in a simple CLI designed to ease execution 
of a series of cmd line commands, bat files, scripts, and executable
program. Each step to be completed is configured using a text file
in CSV format where a record represents a single step to be executed.

## Origination

This project was created for [Chipster.org](https://www.chipster.org) for use in configurating and starting their software components all in a single click for their users.  The goals were to support executing BAT files and EXE files as needed to startup the Chipster Software environment and prepare it for use.

## Step Records

A step is a comma separated set of values on a single line in the steps file. Each step record consists of 3 values:

> 1. Step Type:  `BAT`, `CMD`, or `EXE` - step type determines how the step Text is handled when the step is run.
> 2. Step Result: `required`, or `optional` - step result controls whether or not an error during a step will stop the rest of the run.
> 3. Step Text: `test\\test.bat`, or `echo hello world`, or `notepad.exe`

### Steps File

The steps file allows users to easily edit the steps that are being
executed. Each line in the steps file is a single step and all steps 
are loaded when the steps file is resolved. The steps file itself 
can be located in a few places or passed as a parameter to the 
steps program.

#### Steps File Location

The logic used to load the steps file is applied to each possible 
location in sequence. The order and description of the locations
is as follows:

1. First there is a default steps file constant value of `.steps` which is expected to be in the same folder as the steps.exe program.
2. Second the environment variable `STEPS_FILENAME` is check for a value which is expected to be relative to the current working directory.
3. Third the user can pass in the relative path and name of the steps file when running the steps.exe program.
4. Finally if a steps file has not been found the steps program will search recursively from the current working directory and find the first file meeting the `*.steps` pattern.

### Supported Step Types

| Type | Description |
| --- | --- |
| BAT | Windows batch file |
| CMD | Windows cmd string |
| EXE | Executable file |

### Supported Step Results

`required` - this value means that if an error occurs in a step, subsequent steps will not be ran.
`optional` - this value means that the steps will continue to be ran even if this one end in an error.


## Elevated Execution

The steps program should be run as an administrator in an elevated process so that the run will 
require a UAC prompt but this single UAC prompt will cover the entire list of steps. Elevated execution
can be achieved using the `-elevated` switch.

### Some Examples

You will find some test batch files and command text use in the .\steps\steps_test.go file. And here are some other examples:

**BAT** Step Type:

```cmd
BAT,required,.\\test\\test.bat
BAT,optional,c:\\Temp\\StartSearch.bat
```

**CMD** Step Type:

```cmd
CMD,optional,cd .\\cmd && dir
CMD,required,clear
```

**EXE** Step Type:

```cmd
EXE,required,notepad.exe .\\README.md
EXE,optional,calc.exe
```

### Switches
The following switches maybe be passed to the steps.exe in order to alter or enhance its execution behavior.

| Switch | Effect |
| ------ | ------ |
| -verbose | Turn on verbose logging |
| -elevated | Executes all steps as elevated user. |

Note: using the `-elevated` switch will result in a single UAC prompt at the beginning of the execution.

## Some Important TIPS

> - The paths used in each Step Text value are resolved using the current working directory of the Steps executable and should therefore be relative to the folder where `steps.exe` is run.
> - Each Step is run in its own process and therefore BAT files should be single purpose. NOTE: cd (change dir) commands are not honored inside a BAT file.
> - If a Step is configured as `required` then it will be executed as a **blocking** step and the next step will not be started until it exits.
> - If a Step is configured as `optional` then it will be exeucted as **non-blocker** and the next step will start w/o it having exited.

Many combinations of Steps can be run together as a large contigous task but using the utility to test each step at a time is the best path forward.


