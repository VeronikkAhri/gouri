package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Config struct {
	Editor      string            `json:"editor,omitempty"`
	ShellConfig string            `json:"shell_config,omitempty"`
	Custom      map[string]string `json:"custom,omitempty"`
}

func configPath() string {
	home := os.Getenv("HOME")
	if home == "" {
		if h, err := os.UserHomeDir(); err == nil {
			home = h
		}
	}
	return filepath.Join(home, ".gouri.json")
}

func loadConfig() (*Config, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}

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
	fmt.Println("  gouri head FILE [N]      # show first N lines")
	fmt.Println("  gouri tail FILE [N]      # show last N lines")
	fmt.Println("  gouri wc FILE            # count lines, words and bytes")
	fmt.Println("  gouri alias list         # list defined aliases")
	fmt.Println("  gouri mkdir DIR          # create a directory")
	fmt.Println("  gouri uptime             # show system uptime")
	fmt.Println("  gouri edit FILE          # open file in $EDITOR")
	fmt.Println("  gouri env KEY            # print environment variable")
	fmt.Println("  gouri env set KEY VAL    # persist environment variable")
	fmt.Println("  gouri free               # show memory usage")
	fmt.Println("  gouri ps                 # list running processes")
	fmt.Println("  gouri kill PID           # terminate a process")
	fmt.Println("  gouri echo TEXT          # print text")
	fmt.Println("  gouri cpuinfo            # show CPU information")
	fmt.Println("  gouri compress OUT FILES # create a tar.gz archive")
	fmt.Println("  gouri extract ARCH DIR   # extract a tar.gz archive")
	fmt.Println("  gouri whoami             # show current user")
	fmt.Println("  gouri date               # show date and time")
	fmt.Println("  gouri net                # show network interfaces")
	fmt.Println("  gouri hostname           # print host name")
	fmt.Println("  gouri calc EXPR          # evaluate expression")
	fmt.Println("  gouri open PATH          # open file or directory")
	fmt.Println("  gouri download URL FILE  # download URL to FILE")
	fmt.Println("  gouri serve DIR PORT     # start http server")
	fmt.Println("  gouri uuid               # generate a UUID")
	fmt.Println("  gouri checksum FILE      # SHA256 of FILE")
	fmt.Println("  gouri encrypt IN OUT PASS  # encrypt file with password")
	fmt.Println("  gouri decrypt IN OUT PASS  # decrypt file with password")
	fmt.Println("  gouri sysinfo            # show OS and arch")
	fmt.Println("  gouri clear              # clear the screen")
	fmt.Println("  gouri config get KEY     # show config value")
	fmt.Println("  gouri config set KEY VAL # set config value")
	fmt.Println("  gouri config path        # print config location")
	fmt.Println("  gouri manual             # show the full manual")
}

func runCommand(name string, args ...string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		psArgs := append([]string{"-Command", name}, args...)
		cmd = exec.Command("powershell", psArgs...)
	} else {
		cmd = exec.Command(name, args...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func shellConfig() string {
	if cfg, err := loadConfig(); err == nil && cfg.ShellConfig != "" {
		return cfg.ShellConfig
	}
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

func headFile(path string, n int) error {
	if n <= 0 {
		n = 10
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		count++
		if count >= n {
			break
		}
	}
	return scanner.Err()
}

func tailFile(path string, n int) error {
	if n <= 0 {
		n = 10
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if n > len(lines) {
		n = len(lines)
	}
	for _, line := range lines[len(lines)-n:] {
		fmt.Println(line)
	}
	return nil
}

func wordCount(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	text := string(data)
	lines := strings.Count(text, "\n")
	words := len(strings.Fields(text))
	fmt.Printf("%d %d %d\n", lines, words, len(data))
	return nil
}

func encryptFile(in, out, password string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	enc := make([]byte, len(data))
	stream.XORKeyStream(enc, data)
	outData := append(iv, enc...)
	return os.WriteFile(out, outData, 0644)
}

func decryptFile(in, out, password string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	if len(data) < aes.BlockSize {
		return fmt.Errorf("ciphertext too short")
	}
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}
	iv := data[:aes.BlockSize]
	enc := data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	dec := make([]byte, len(enc))
	stream.XORKeyStream(dec, enc)
	return os.WriteFile(out, dec, 0644)
}

func listAliases() error {
	data, err := os.ReadFile(shellConfig())
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "alias ") {
			fmt.Println(line)
		}
	}
	return scanner.Err()
}

func makeDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func showUptime() error {
	if runtime.GOOS == "windows" {
		return runCommand("Get-Uptime")
	}
	return runCommand("uptime")
}

func editFile(path string) error {
	editor := os.Getenv("EDITOR")
	if cfg, err := loadConfig(); err == nil && cfg.Editor != "" {
		editor = cfg.Editor
	}
	if editor == "" {
		editor = "nano"
	}
	return runCommand(editor, path)
}

func showEnv(key string) {
	fmt.Println(os.Getenv(key))
}

func setEnv(key, value string) error {
	file := shellConfig()
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "\nexport %s=%q\n", key, value)
	return err
}

func showMemory() error {
	if runtime.GOOS == "windows" {
		return runCommand("Get-CimInstance", "Win32_OperatingSystem")
	}
	return runCommand("free", "-h")
}

func listProcesses() error {
	if runtime.GOOS == "windows" {
		return runCommand("Get-Process")
	}
	return runCommand("ps", "aux")
}

func compressFiles(out string, files []string) error {
	args := append([]string{"-czf", out}, files...)
	return runCommand("tar", args...)
}

func extractArchive(archive, dir string) error {
	return runCommand("tar", "-xzf", archive, "-C", dir)
}

func showUser() error {
	return runCommand("whoami")
}

func showDate() error {
	return runCommand("date")
}

func showNetwork() error {
	if runtime.GOOS == "windows" {
		return runCommand("Get-NetIPAddress")
	}
	return runCommand("ip", "addr")
}

func showCPUInfo() error {
	if runtime.GOOS == "windows" {
		return runCommand("Get-CimInstance", "Win32_Processor")
	}
	if _, err := os.Stat("/proc/cpuinfo"); err == nil {
		return runCommand("cat", "/proc/cpuinfo")
	}
	return runCommand("lscpu")
}

func killProcess(pid string) error {
	if runtime.GOOS == "windows" {
		return runCommand("Stop-Process", "-Id", pid, "-Force")
	}
	return runCommand("kill", "-9", pid)
}

func echoText(text string) {
	fmt.Println(text)
}

func showHostname() error {
	name, err := os.Hostname()
	if err != nil {
		return err
	}
	fmt.Println(name)
	return nil
}

func calcExpr(expr string) error {
	if runtime.GOOS == "windows" {
		return runCommand("Write-Output", fmt.Sprintf("(%s)", expr))
	}
	return runCommand("bash", "-c", fmt.Sprintf("echo '%s' | bc -l", expr))
}

func openPath(path string) error {
	switch runtime.GOOS {
	case "windows":
		return runCommand("Start-Process", path)
	case "darwin":
		return runCommand("open", path)
	default:
		return runCommand("xdg-open", path)
	}
}

func downloadFile(url, out string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func serveDir(dir, port string) error {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)
	return http.ListenAndServe(":"+port, nil)
}

func newUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func checksumFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	sum := h.Sum(nil)
	fmt.Println(hex.EncodeToString(sum))
	return nil
}

func showSysInfo() {
	fmt.Printf("%s %s\n", runtime.GOOS, runtime.GOARCH)
}

func showPwd() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

func showHistory() error {
	if runtime.GOOS == "windows" {
		return runCommand("Get-History")
	}
	file := os.Getenv("HISTFILE")
	if file == "" {
		file = filepath.Join(os.Getenv("HOME"), ".bash_history")
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}

func clearScreen() error {
	if runtime.GOOS == "windows" {
		return runCommand("cls")
	}
	return runCommand("clear")
}

const manualText = `Gouri Manual

Gouri is a simple terminal assistant that automates common shell tasks.

Available commands:
  update             run system update
  upgrade            run system upgrade
  alias add          create an alias
  alias remove       remove an alias
  alias list         list defined aliases
  view               print file contents
  remove             delete a file
  list               list directory contents
  copy               copy a file
  move               move or rename a file
  search             search term in file
  disk               show disk usage
  ping               ping a network host
  tree               show directory tree
  create             create an empty file
  lines              count lines in file
  head               show first lines of file
  tail               show last lines of file
  wc                 count lines, words and bytes
  mkdir              create a directory
  uptime             show system uptime
  edit               open file in $EDITOR
  env                get or set environment variables
  free               show memory usage
  ps                 list running processes
  kill               terminate a process
  echo               print text
  cpuinfo            show CPU information
  compress           create a tar.gz archive
  extract            extract a tar.gz archive
  whoami             show current user
  date               show date and time
  net                show network interfaces
  hostname           print host name
  calc               evaluate expression
  open               open file or directory
  download           download URL to file
  serve              start http server
  uuid               generate a UUID
  checksum           SHA256 of a file
  encrypt            encrypt a file with password
  decrypt            decrypt a file with password
  sysinfo            show OS and architecture
  clear              clear the screen
  config             manage configuration values
  pwd                print working directory
  history            show shell history
  manual             show this manual`

func showManual() {
	fmt.Println(manualText)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "update":
		if runtime.GOOS == "windows" {
			if err := runCommand("winget", "upgrade"); err != nil {
				fmt.Println("update error:", err)
			}
		} else {
			if err := runCommand("sudo", "apt-get", "update"); err != nil {
				fmt.Println("update error:", err)
			}
		}
	case "upgrade":
		if runtime.GOOS == "windows" {
			if err := runCommand("winget", "upgrade", "--all"); err != nil {
				fmt.Println("upgrade error:", err)
			}
		} else {
			if err := runCommand("sudo", "apt-get", "upgrade", "-y"); err != nil {
				fmt.Println("upgrade error:", err)
			}
		}
	case "alias":
		if len(os.Args) < 3 {
			usage()
			return
		}
		switch os.Args[2] {
		case "add":
			if len(os.Args) < 5 {
				usage()
				return
			}
			if err := addAlias(os.Args[3], strings.Join(os.Args[4:], " ")); err != nil {
				fmt.Println("alias add error:", err)
			}
		case "remove":
			if len(os.Args) < 4 {
				usage()
				return
			}
			if err := removeAlias(os.Args[3]); err != nil {
				fmt.Println("alias remove error:", err)
			}
		case "list":
			if err := listAliases(); err != nil {
				fmt.Println("alias list error:", err)
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
		if runtime.GOOS == "windows" {
			if err := runCommand("Get-PSDrive"); err != nil {
				fmt.Println("disk error:", err)
			}
		} else {
			if err := runCommand("df", "-h"); err != nil {
				fmt.Println("disk error:", err)
			}
		}
	case "ping":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if runtime.GOOS == "windows" {
			if err := runCommand("ping", "-n", "4", os.Args[2]); err != nil {
				fmt.Println("ping error:", err)
			}
		} else if err := runCommand("ping", "-c", "4", os.Args[2]); err != nil {
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
	case "head":
		if len(os.Args) < 3 {
			usage()
			return
		}
		n := 10
		if len(os.Args) >= 4 {
			if v, err := strconv.Atoi(os.Args[3]); err == nil {
				n = v
			}
		}
		if err := headFile(os.Args[2], n); err != nil {
			fmt.Println("head error:", err)
		}
	case "tail":
		if len(os.Args) < 3 {
			usage()
			return
		}
		n := 10
		if len(os.Args) >= 4 {
			if v, err := strconv.Atoi(os.Args[3]); err == nil {
				n = v
			}
		}
		if err := tailFile(os.Args[2], n); err != nil {
			fmt.Println("tail error:", err)
		}
	case "wc":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := wordCount(os.Args[2]); err != nil {
			fmt.Println("wc error:", err)
		}
	case "mkdir":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := makeDir(os.Args[2]); err != nil {
			fmt.Println("mkdir error:", err)
		}
	case "uptime":
		if err := showUptime(); err != nil {
			fmt.Println("uptime error:", err)
		}
	case "edit":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := editFile(os.Args[2]); err != nil {
			fmt.Println("edit error:", err)
		}
	case "env":
		if len(os.Args) < 3 {
			usage()
			return
		}
		switch os.Args[2] {
		case "set":
			if len(os.Args) < 5 {
				usage()
				return
			}
			if err := setEnv(os.Args[3], os.Args[4]); err != nil {
				fmt.Println("env set error:", err)
			}
		default:
			showEnv(os.Args[2])
		}
	case "free":
		if err := showMemory(); err != nil {
			fmt.Println("free error:", err)
		}
	case "ps":
		if err := listProcesses(); err != nil {
			fmt.Println("ps error:", err)
		}
	case "kill":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := killProcess(os.Args[2]); err != nil {
			fmt.Println("kill error:", err)
		}
	case "echo":
		if len(os.Args) < 3 {
			usage()
			return
		}
		echoText(strings.Join(os.Args[2:], " "))
	case "cpuinfo":
		if err := showCPUInfo(); err != nil {
			fmt.Println("cpuinfo error:", err)
		}
	case "compress":
		if len(os.Args) < 4 {
			usage()
			return
		}
		if err := compressFiles(os.Args[2], os.Args[3:]); err != nil {
			fmt.Println("compress error:", err)
		}
	case "extract":
		if len(os.Args) < 4 {
			usage()
			return
		}
		if err := extractArchive(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("extract error:", err)
		}
	case "whoami":
		if err := showUser(); err != nil {
			fmt.Println("whoami error:", err)
		}
	case "date":
		if err := showDate(); err != nil {
			fmt.Println("date error:", err)
		}
	case "net":
		if err := showNetwork(); err != nil {
			fmt.Println("net error:", err)
		}
	case "hostname":
		if err := showHostname(); err != nil {
			fmt.Println("hostname error:", err)
		}
	case "calc":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := calcExpr(strings.Join(os.Args[2:], " ")); err != nil {
			fmt.Println("calc error:", err)
		}
	case "open":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := openPath(os.Args[2]); err != nil {
			fmt.Println("open error:", err)
		}
	case "download":
		if len(os.Args) < 4 {
			usage()
			return
		}
		if err := downloadFile(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("download error:", err)
		}
	case "serve":
		dir := "."
		port := "8080"
		if len(os.Args) >= 3 {
			dir = os.Args[2]
		}
		if len(os.Args) >= 4 {
			port = os.Args[3]
		}
		if err := serveDir(dir, port); err != nil {
			fmt.Println("serve error:", err)
		}
	case "uuid":
		id, err := newUUID()
		if err != nil {
			fmt.Println("uuid error:", err)
			break
		}
		fmt.Println(id)
	case "checksum":
		if len(os.Args) < 3 {
			usage()
			return
		}
		if err := checksumFile(os.Args[2]); err != nil {
			fmt.Println("checksum error:", err)
		}
	case "encrypt":
		if len(os.Args) < 5 {
			usage()
			return
		}
		if err := encryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
			fmt.Println("encrypt error:", err)
		}
	case "decrypt":
		if len(os.Args) < 5 {
			usage()
			return
		}
		if err := decryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
			fmt.Println("decrypt error:", err)
		}
	case "sysinfo":
		showSysInfo()
	case "pwd":
		if err := showPwd(); err != nil {
			fmt.Println("pwd error:", err)
		}
	case "history":
		if err := showHistory(); err != nil {
			fmt.Println("history error:", err)
		}
	case "clear":
		if err := clearScreen(); err != nil {
			fmt.Println("clear error:", err)
		}
	case "config":
		if len(os.Args) < 3 {
			usage()
			return
		}
		cfg, err := loadConfig()
		if err != nil {
			fmt.Println("config error:", err)
			return
		}
		switch os.Args[2] {
		case "get":
			if len(os.Args) < 4 {
				usage()
				return
			}
			key := os.Args[3]
			switch key {
			case "editor":
				fmt.Println(cfg.Editor)
			case "shell_config":
				fmt.Println(cfg.ShellConfig)
			default:
				fmt.Println(cfg.Custom[key])
			}
		case "set":
			if len(os.Args) < 5 {
				usage()
				return
			}
			key := os.Args[3]
			val := os.Args[4]
			switch key {
			case "editor":
				cfg.Editor = val
			case "shell_config":
				cfg.ShellConfig = val
			default:
				if cfg.Custom == nil {
					cfg.Custom = make(map[string]string)
				}
				cfg.Custom[key] = val
			}
			if err := saveConfig(cfg); err != nil {
				fmt.Println("config set error:", err)
			}
		case "path":
			fmt.Println(configPath())
		default:
			usage()
		}
	case "manual":
		showManual()
	default:
		usage()
	}
}
