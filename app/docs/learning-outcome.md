# Learning Go: Shell and OS Concepts

While learning Go and building a small shell-like program, I came across a few concepts that were new to me. I am documenting what I understood along the way.

## 1. Taking Input from the User

One of the first things I learned was the difference between `Scan()`, `Scanln()`, `Scanf()` and `ReadString()`.

### The issue with `Scan()`, `Scanln()` and `Scanf()`

The main thing to understand is that these functions treat **whitespace as a separator** when reading values.

For example:

```text
John Doe
```

If I read this as a single string using `Scan()`, it won't treat `"John Doe"` as one value. It reads `John` and considers the whitespace as a separator.

This works fine when I want to read individual values such as:

```text
10
25
hello
```

But it becomes a problem when the input itself can contain spaces, such as:

```text
John Doe
hello world
echo hello world
```

Since I am building a shell, commands can contain multiple words, so I need a way to read the complete line.

### Using `ReadString()`

For this, I can use `bufio.Reader`:

```go
reader := bufio.NewReader(os.Stdin)

input, err := reader.ReadString('\n')
```

`ReadString('\n')` keeps reading until it finds the specified delimiter.

A **delimiter** is a character that tells the program where one piece of input ends.

Here:

```go
'\n'
```

is the delimiter, so the reader reads the complete line until the user presses Enter.

For example:

```text
echo hello world
```

can be read as one complete input instead of splitting it at every whitespace.

This is useful for a shell because I first need to receive the complete command and then parse it.

---

## 2. REPL — Read, Eval, Print, Loop

While building the shell, I also learned about the **REPL** concept.

REPL stands for:

* **Read** — Read the user's input
* **Eval** — Evaluate/parse and execute the command
* **Print** — Display the output or error
* **Loop** — Go back and wait for the next command

A shell basically follows this cycle:

```text
        ┌──────────────┐
        │     Read     │
        └──────┬───────┘
               ↓
        ┌──────────────┐
        │     Eval     │
        └──────┬───────┘
               ↓
        ┌──────────────┐
        │    Print     │
        └──────┬───────┘
               ↓
        ┌──────────────┐
        │     Loop     │
        └──────┬───────┘
               │
               └──────────→ Read again
```

This cycle continues until the shell exits or the process is terminated.

For my shell, the basic flow is:

```text
$ echo hello
hello

$ pwd
/current/directory

$ cd /some/directory

$ pwd
/some/directory
```

---

## 3. Locating Executable Files

Another interesting thing I learned was how commands are actually located.

When I type:

```bash
whoami
```

in Git Bash, the shell searches through directories in its `PATH` and finds the executable.

For example, I can use:

```bash
type whoami
```

and get something similar to:

```text
C:\Program Files\Git\usr\bin\whoami.exe
```

### Git Bash vs Windows

This became confusing when I tried:

```go
os.Chdir("/usr/bin")
```

from my Go program.

It failed with:

```text
The system cannot find the path specified.
```

But this worked perfectly in Git Bash:

```bash
cd /usr/bin
```

The reason is that **Git Bash provides a Unix-like environment on top of Windows**.

Git Bash understands paths such as:

```text
/usr/bin
```

and maps them to the corresponding Windows filesystem location.

My Go program, however, was running as a Windows program. Therefore, when I used:

```go
os.Chdir("/usr/bin")
```

the Windows filesystem APIs were being used, and `/usr/bin` was not a normal Windows path.

I could find the corresponding Windows path using:

```bash
cygpath -w /usr/bin
```

This converts a Unix-style path into a Windows-style path.

This helped me understand that the path I see in a shell is not necessarily the same path representation that the underlying operating system uses.

---

## 4. Finding and Running Executable Commands from Go

Go provides `os/exec` for working with external programs.

For example:

```go
path, err := exec.LookPath("grep")
```

`LookPath()` searches the directories in `PATH` and returns the executable's path if it can find it.

Once the executable is found, I can run it using:

```go
cmd := exec.Command("grep", "hello", "file.txt")
```

The important thing I learned here is how command arguments are represented.

If the user enters:

```text
grep hello file.txt
```

I can split the input into:

```text
["grep", "hello", "file.txt"]
```

Then:

```go
command := parts[0]
args := parts[1:]
```

gives:

```text
command = "grep"

args = ["hello", "file.txt"]
```

And:

```go
exec.Command(command, args...)
```

effectively becomes:

```go
exec.Command("grep", "hello", "file.txt")
```

The `...` is used to pass the slice as individual arguments.

One thing I also realized is that simply using:

```go
strings.Fields(input)
```

is not enough to implement a complete shell parser because it does not understand things like quoted strings.

For example:

```text
echo "hello world"
```

A real shell should treat `"hello world"` as one argument, while `strings.Fields()` would simply split it based on whitespace.

---

## 5. Built-in Commands vs External Commands

Another important concept I came across is the difference between **built-in commands** and **external commands**.

### Built-in commands

Built-in commands are commands that are implemented by the shell itself rather than being separate executable programs.

For example:

```text
cd
pwd
```

can be implemented directly inside the shell.

For my Go shell, I can handle them myself:

```go
switch command {
case "pwd":
    // use os.Getwd()

case "cd":
    // use os.Chdir()
}
```

One important example is `cd`.

`cd` needs to change the **current working directory of the shell itself**.

If I created a separate child process just to execute `cd`, that child process could change its own working directory, but once it exits, my shell's working directory would remain unchanged.

That is why `cd` needs to be handled by the shell process itself.

### External commands

External commands are separate executable programs.

Examples include:

```text
grep
ls
chmod
chown
ps
```

On Windows, many of these may exist as `.exe` files, especially when working with environments such as Git for Windows.

The general flow is:

```text
User enters command
       ↓
Shell searches PATH
       ↓
Find executable
       ↓
Create/start process
       ↓
Execute program
       ↓
Read/display output
```

For example:

```text
ls
 ↓
find ls executable
 ↓
start ls process
 ↓
ls executes
 ↓
output is returned to the shell
```

---

## What I Understand So Far

Putting everything together, the basic idea of my shell currently looks like this:

```text
User
 ↓
Enter command
 ↓
Read complete line
 ↓
Parse command + arguments
 ↓
Is it a built-in?
 ├── Yes → Execute inside shell
 │          ├── cd  → os.Chdir()
 │          └── pwd → os.Getwd()
 │
 └── No  → Find external executable
            ↓
          Execute process
            ↓
          Display output
            ↓
          Read next command
```

This helped me understand that a shell is not just a program that "runs commands". It is responsible for reading input, parsing it, deciding what the command means, handling built-ins, finding external programs, creating processes and connecting their input/output.

There is still a lot more to understand, especially command parsing, quoting, redirection, pipes, environment variables and process handling. But this is what I have understood so far while building my shell in Go.
