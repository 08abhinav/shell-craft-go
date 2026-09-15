# Building My Own Shell Using Go

A shell is an interface between the user and the operating system. It takes commands entered by the user, interprets them, and then executes them.

For example, when we type:

```bash
ls
```

the shell needs to figure out what `ls` means, find the corresponding executable, and start it.

There are generally two types of commands that a shell handles:

### Built-in commands

Some commands are implemented directly by the shell itself. These are called **built-in commands**.

Examples include:

```bash
cd
pwd
exit
```

These commands are handled by the shell rather than starting a separate executable process.

For example, `cd` needs to change the current working directory of the shell process itself. If it were executed as a separate child process, the directory change would disappear when that process exited.

### Executable commands

Other commands are provided by executable files installed on the system.

Examples include:

```bash
cat
ls
grep
ps
chmod
chown
```

When the shell receives one of these commands, it searches the directories listed in the `PATH` environment variable to find the corresponding executable. Once found, the shell starts a new process to execute it.

For example:

```text
ls
 ↓
Search PATH
 ↓
Find ls executable
 ↓
Start process
 ↓
Execute ls
 ↓
Display output
```

## About This Project

In this project, I am building a small shell from scratch using **Go** to understand what happens behind the scenes when we type commands into a terminal.

The project covers things such as:

* Handling built-in commands
* Finding executables using `PATH`
* Passing arguments to commands
* Creating and running processes using Go's `os/exec` package
* Reading command output
* Handling errors
* Working with directories using the `os` package
* Understanding how commands are parsed and executed

The goal isn't to build a production-ready shell. The main goal is to understand the basic building blocks behind a shell and how it interacts with the operating system.
