# Manual

The following sections describe each command provided by Gouri.
Commands work on both Linux and Windows PowerShell unless noted.

## update
Run system package update (apt on Linux, winget on Windows).

## upgrade
Upgrade installed packages.

## alias add NAME CMD
Append an alias to your shell configuration file.

## alias remove NAME
Remove an alias from your shell configuration file.

## alias list
List all defined aliases found in the shell configuration.

## view FILE
Print the contents of FILE.

## remove FILE
Delete FILE from the filesystem.

## list DIR
List the contents of DIR.

## copy SRC DST
Copy SRC to DST.

## move SRC DST
Move or rename SRC to DST.

## search FILE TERM
Search TERM inside FILE and print matching lines.

## disk
Display disk usage information.

## ping HOST
Send four ICMP echo requests to HOST.

## tree DIR
Display a recursive directory tree starting at DIR.

## create FILE
Create an empty file called FILE.

## lines FILE
Print the number of lines in FILE.

## head FILE [N]
Show the first N lines of FILE (default 10).

## tail FILE [N]
Show the last N lines of FILE (default 10).

## wc FILE
Display line, word and byte counts of FILE.

## mkdir DIR
Create directory DIR and any missing parents.

## uptime
Show how long the system has been running.

## edit FILE
Open FILE in the configured editor.

## env KEY
Print the value of environment variable KEY.

## env set KEY VAL
Persist KEY=VAL in the shell configuration file.

## free
Show memory usage statistics.

## ps
List running processes.

## kill PID
Terminate the process with the given PID.

## echo TEXT
Print TEXT to the terminal.

## cpuinfo
Display information about the CPU.

## compress OUT FILES
Create a tar.gz archive OUT from FILES.

## extract ARCH DIR
Extract ARCH archive into DIR.

## whoami
Show the current user.

## date
Display the current date and time.

## net
List network interfaces.

## hostname
Print the host name.

## calc EXPR
Evaluate the expression EXPR using bc or PowerShell.

## open PATH
Open PATH with the system default application.

## download URL FILE
Download URL to FILE on disk.

## serve DIR PORT
Serve DIR over HTTP on PORT.

## uuid
Generate and print a UUID.

## checksum FILE
Print the SHA256 checksum of FILE.

## encrypt IN OUT PASS
Encrypt IN to OUT using PASS as the password.

## decrypt IN OUT PASS
Decrypt IN to OUT using PASS as the password.

## sysinfo
Print the operating system and architecture.

## clear
Clear the terminal screen.

## config get KEY
Print a saved configuration value.

## config set KEY VAL
Persist KEY with value VAL in the configuration file.

## config path
Show the location of the configuration file.

## pwd
Print the current working directory.

## history
Show shell command history.

## manual
Display this manual text in the terminal.
