# Steps Utility
Simple CLI designed to ease the execution of a series of cmd line, 
bat file, scripts, etc. Each step to be completed is sourced from
a text manifest that serves as a direct guide to the order of 
execution and the steps to be executed. 

## Steps Manifest 
The steps manifest file allows users to easily edit what is being
executed. Each line in the steps manifest file is one of a several
optional strings that describe a step in the execution.  Each step
in the manifest can also be configured as required, or optional to
support either making sure a step is successful before going on to
the next one, or ignoring the result of a step and continuing to
the next step regardless of the result.

## What is a Step?
A step is a command line task that must return an integer value that
indicates it was successful or not.  A non-zero integer result tells
the steps program that the result was unsuccessful. 

Each step has is one of a supported list of types.

###### Supported Types

| Type | Description |
| --- | --- |
| BAT | Windows batch file |
| CMD | Windows cmd string |
| EXE | Executable file |

upcoming types are VBS, JS, PS... as needed

## Elevated Execution
The steps program will run in an elevated process and therefore 
require a UAC prompt but this single UAC prompt will cover the
entire manifest of steps.


