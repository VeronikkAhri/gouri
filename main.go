package main

import (
	"bufio"
	"fmt"
	"io"
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
	fmt.Println("  gouri list DIR           # list directory contents")
	fmt.Println("  gouri copy SRC DST       # copy a file")
	fmt.Println("  gouri move SRC DST       # move/rename a file")
	fmt.Println("  gouri search FILE TERM   # search term in file")
	fmt.Println("  gouri disk               # show disk usage")
	fmt.Println("  gouri ping HOST          # ping a network host")
	fmt.Println("  gouri tree DIR           # show directory tree")
	fmt.Println("  gouri create FILE        # create an empty file")
	fmt.Println("  gouri lines FILE         # count lines in file")
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

func listDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		fmt.Println(e.Name())
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func moveFile(src, dst string) error {
	return os.Rename(src, dst)
}

func searchInFile(path, term string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, term) {
			fmt.Printf("%d: %s\n", lineNum, line)
		}
		lineNum++
	}
	return scanner.Err()
}

func createFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

func countLines(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Println(count)
	return nil
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
	case "list":
		dir := "."
		if len(os.Args) >= 3 {
			dir = os.Args[2]
		}
		if err := listDir(dir); err != nil {
			fmt.Println("list error:", err)
		}
	case "copy":
		if len(os.Args) < 4 {
			usage()
			return
		}
		if err := copyFile(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("copy error:", err)
		}
	case "move":
		if len(os.Args) < 4 {
			usage()
			return
		}
		if err := moveFile(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("move error:", err)
		}
	case "search":
		if len(os.Args) < 4 {
			usage()
			return
		}
		if err := searchInFile(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("search error:", err)
		}
	case "disk":
		if err := runCommand("df", "-h"); err != nil {
			fmt.Println("disk error:", err)
		}
	case "ping":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := runCommand("ping", "-c", "4", os.Args[2]); err != nil {
			fmt.Println("ping error:", err)
		}
	case "tree":
		dir := "."
		if len(os.Args) >= 3 {
			dir = os.Args[2]
		}
		if err := runCommand("tree", dir); err != nil {
			fmt.Println("tree error:", err)
		}
	case "create":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := createFile(os.Args[2]); err != nil {
			fmt.Println("create error:", err)
		}
	case "lines":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := countLines(os.Args[2]); err != nil {
			fmt.Println("lines error:", err)
		}
	default:
		usage()
	}
}
