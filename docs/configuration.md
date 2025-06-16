# Configuration

Gouri stores its settings in `~/.gouri.json`. Values are read on start and can
be changed using `gouri config` commands.

Supported keys:

- `editor` – the text editor launched by the `edit` command.
- `shell_config` – the shell configuration file used for alias creation.
- any custom key/value pairs that scripts may need.

Use `gouri config set KEY VALUE` to update a value and `gouri config get KEY`
to read it. The location of the file can be printed with `gouri config path`.
