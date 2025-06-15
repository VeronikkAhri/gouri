package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Println("Gouri - simple assistant")
	fmt.Println("Usage:")
	fmt.Println("  gouri update             # run system update")
	fmt.Println("  gouri upgrade            # run system upgrade")
	fmt.Println("  gouri alias add name cmd # create an alias")
	fmt.Println("  gouri alias remove name  # remove an alias")
	fmt.Println("  gouri view FILE          # view file contents")
	fmt.Println("  gouri remove FILE        # remove a file")
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func shellConfig() string {
	shell := os.Getenv("SHELL")
	base := filepath.Base(shell)
	switch base {
	case "zsh":
		return filepath.Join(os.Getenv("HOME"), ".zshrc")
	default:
		return filepath.Join(os.Getenv("HOME"), ".bashrc")
	}
}

func addAlias(name, command string) error {
	file := shellConfig()
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "\nalias %s='%s'\n", name, command)
	return err
}

func removeAlias(name string) error {
	file := shellConfig()
	input, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(input), "\n")
	var out []string
	for _, line := range lines {
		if strings.HasPrefix(line, "alias "+name+"=") {
			continue
		}
		out = append(out, line)
	}
	return os.WriteFile(file, []byte(strings.Join(out, "\n")), 0644)
}

func viewFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return scanner.Err()
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "update":
		if err := runCommand("sudo", "apt-get", "update"); err != nil {
			fmt.Println("update error:", err)
		}
	case "upgrade":
		if err := runCommand("sudo", "apt-get", "upgrade", "-y"); err != nil {
			fmt.Println("upgrade error:", err)
		}
	case "alias":
		if len(os.Args) < 4 {
			usage()
			return
		}
		switch os.Args[2] {
		case "add":
			if err := addAlias(os.Args[3], strings.Join(os.Args[4:], " ")); err != nil {
				fmt.Println("alias add error:", err)
			}
		case "remove":
			if err := removeAlias(os.Args[3]); err != nil {
				fmt.Println("alias remove error:", err)
			}
		default:
			usage()
		}
	case "view":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := viewFile(os.Args[2]); err != nil {
			fmt.Println("view error:", err)
		}
	case "remove":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := os.Remove(os.Args[2]); err != nil {
			fmt.Println("remove error:", err)
		}
	default:
		usage()
	}
}
