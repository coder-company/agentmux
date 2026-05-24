package commands

import "errors"

var errTmuxMissing = errors.New("tmux is not installed or not in PATH\n\nInstall tmux:\n  macOS:  brew install tmux\n  Debian: sudo apt install tmux\n  Arch:   sudo pacman -S tmux")
