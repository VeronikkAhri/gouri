# Gouri

Gouri is a simple terminal assistant written in Go. It automates common shell
tasks such as updating packages, upgrading the system, creating aliases and
managing files. The commands are intentionally short and easy to remember.

## Build

```bash
# compile the binary
go build
```

## Usage

```
./gouri update             # run "sudo apt-get update"
./gouri upgrade            # run "sudo apt-get upgrade -y"
./gouri alias add NAME CMD # append alias to your shell config
./gouri alias remove NAME  # remove alias from your shell config
./gouri view FILE          # print file contents
./gouri remove FILE        # delete a file
```

The shell configuration file is detected based on the `SHELL` environment
variable. For `zsh` it uses `~/.zshrc` and defaults to `~/.bashrc`.

## Notes

Administrative commands such as update and upgrade require `sudo` privileges.
Ensure your user has the appropriate permissions.
