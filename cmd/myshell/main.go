package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Define the constants for the command strings
const (
	echo_s = "echo"
	exit_s = "exit"
	type_s = "type"
	pwd_s  = "pwd"
	cd_s   = "cd"
)

// Define the map at the package level
var prefixFuncMap = map[string]func([]string){
	echo_s: echoFunc,
	exit_s: exitFunc,
	type_s: typeFunc,
	pwd_s:  pwdFunc,
	cd_s:   cdFunc,
}

var prexixTab = []string{
	echo_s,
	exit_s,
	type_s,
	pwd_s,
	cd_s,
}

func main() {
	for {
		// Print the prompt
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		in, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stdout, "Error during reading command occurred! ", err.Error())
			os.Exit(1)
		}

		inputs := strings.Split(strings.TrimSpace(in), " ")
		cmd := inputs[0]
		args := inputs[1:]

		handled := false
		for prefix, function := range prefixFuncMap {
			if cmd == prefix {
				function(args)
				handled = true
				break
			}
		}
		if !handled {
			command := exec.Command(cmd, args...)
			command.Stderr = os.Stderr
			command.Stdout = os.Stdout

			err = command.Run()
			if err != nil {
				nonexistentFunc(cmd, true)
			}
		}
	}
}

func echoFunc(args []string) {
	fmt.Fprintln(os.Stdout, strings.Join(args, " "))
}

func exitFunc(args []string) {
	if len(args) == 0 {
		os.Exit(1)
	}
	if code, err := strconv.Atoi(args[0]); err == nil {
		os.Exit(code)
	}
}

func typeFunc(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stdout)
		return
	}

	for prefix := range prexixTab {
		if args[0] == prexixTab[prefix] {
			fmt.Fprint(os.Stdout, prexixTab[prefix], " is a shell builtin\n")
			return
		}
	}

	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, path := range paths {
		fp := filepath.Join(path, args[0])
		if _, err := os.Stat(fp); err == nil {
			fmt.Println(fp)
			return
		}
	}

	nonexistentFunc(args[0], false)
}

func pwdFunc(args []string) {
	if len(args) != 0 {
		fmt.Fprintln(os.Stdout)
		return
	}

	pwd, _ := os.Getwd()

	fmt.Fprintln(os.Stdout, pwd)
}

func cdFunc(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stdout, "Usage: cd <directory>")
		return
	}

	targetDir := args[0]

	var err error

	if targetDir == "~" {
		targetDir, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stdout, "cd: could not find home directory")
			return
		}
	}

	if filepath.IsAbs(targetDir) {
		err = os.Chdir(targetDir)
	} else {
		err = os.Chdir(filepath.Clean(filepath.Join(".", targetDir)))
	}

	if err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", targetDir)
	}
}

func nonexistentFunc(cmd string, command bool) {
	if command {
		fmt.Fprint(os.Stdout, cmd, ": command not found\n")
	} else {
		fmt.Fprint(os.Stdout, cmd, ": not found\n")
	}
}
