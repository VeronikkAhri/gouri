# Gouri

Gouri is a simple terminal assistant written in Go. It automates common shell
tasks such as updating packages, upgrading the system, creating aliases and
managing files. The commands are intentionally short and easy to remember.
It works on Linux and Windows PowerShell, abstracting away platform
differences.

## Build

```bash
# compile the binary
go build
```

## Usage

```
./gouri update             # update system packages
./gouri upgrade            # upgrade installed packages
./gouri alias add NAME CMD # append alias to your shell config
./gouri alias remove NAME  # remove alias from your shell config
./gouri view FILE          # print file contents
./gouri remove FILE        # delete a file
./gouri list DIR           # list directory contents
./gouri copy SRC DST       # copy a file
./gouri move SRC DST       # move or rename a file
./gouri search FILE TERM   # search for TERM inside FILE
./gouri disk               # show disk usage
./gouri ping HOST          # ping a host
./gouri tree DIR           # show directory tree
./gouri create FILE        # create an empty file
./gouri lines FILE         # count lines in a file
./gouri head FILE [N]      # show first N lines (default 10)
./gouri tail FILE [N]      # show last N lines (default 10)
./gouri wc FILE            # count lines, words and bytes
./gouri alias list         # list existing aliases
./gouri mkdir DIR          # create a directory
./gouri uptime             # show system uptime
./gouri edit FILE          # open file in $EDITOR
./gouri env KEY            # print environment variable
./gouri env set KEY VAL    # persist environment variable
./gouri free               # show memory usage
./gouri ps                 # list running processes
./gouri kill PID           # terminate a process
./gouri echo TEXT          # print text
./gouri cpuinfo            # show CPU information
./gouri compress OUT FILES # create a tar.gz archive
./gouri extract ARCH DIR   # extract a tar.gz archive
./gouri whoami             # show current user
./gouri date               # show date and time
./gouri net                # show network interfaces
./gouri hostname           # print host name
./gouri calc EXPR          # evaluate expression
./gouri open PATH          # open file or directory
./gouri download URL FILE  # download URL to FILE
./gouri serve DIR PORT     # start http server
./gouri uuid               # generate a UUID
./gouri checksum FILE      # SHA256 of FILE
./gouri sysinfo            # show OS and arch
./gouri clear              # clear the screen
./gouri config get KEY     # show saved config value
./gouri config set KEY VAL # set config value
./gouri config path        # print config file location
./gouri pwd                # print working directory
./gouri history            # show shell history
./gouri manual             # show full manual
```

Run `./gouri manual` to view extended documentation for all commands.

The shell configuration file is detected based on the `SHELL` environment
variable. For `zsh` it uses `~/.zshrc` and defaults to `~/.bashrc`.

Gouri stores customization options in `~/.gouri.json`. You can override the
default editor or shell configuration file path using the config commands.
Other custom keys may also be saved and retrieved for scripts.

Example:

```bash
# set nano as the default editor
./gouri config set editor nano
# show the saved editor
./gouri config get editor
```

## Notes

Administrative commands such as update and upgrade require `sudo` privileges.
Ensure your user has the appropriate permissions.
