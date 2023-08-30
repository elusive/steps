# Steps Utility
Simple CLI designed to ease the execution of a series of cmd line, 
bat file, scripts, etc. Each step to be completed is sourced from
a text manifest that serves as a direct guide to the order of 
execution and the steps to be executed. 

## Steps File
The steps file allows users to easily edit the steps that are being
executed. Each line in the steps file is a CSV record that specifies
the strings that describe a step in the series of steps that make up
a steps program run. 

#### Steps File Location
There are a few possible locations for the steps file and they are 
tried in order using the following logic:
> 1. Default Steps File:  the default steps file is 

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


