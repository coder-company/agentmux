package tmux

import (
	"strconv"
	"strings"
	"time"

	"agentmux/internal/core"
)

// fieldSep is the delimiter used in tmux format strings.
//
// We use a horizontal TAB (ASCII 9): tmux 3.4+ escapes any control byte under
// 0x20 (except TAB) into a literal 4-character "\NNN" sequence in format output,
// so the previously-used Unit Separator (0x1f) was rendered as the bytes "\037"
// rather than the actual byte 0x1f, breaking the parser. TABs cannot appear in
// tmux session names and are virtually never present in file paths.
const fieldSep = "\t"

// Format string for tmux list-sessions -F
var sessionFormat = "#{session_name}" + fieldSep + "#{session_windows}" + fieldSep + "#{session_created}" + fieldSep + "#{session_attached}" + fieldSep + "#{pane_current_path}"

// ParseSessions parses tmux list-sessions output using the sessionFormat.
func ParseSessions(output string) ([]core.Session, error) {
	var sessions []core.Session
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		s, err := parseLine(line)
		if err != nil {
			continue // skip malformed lines
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func parseLine(line string) (core.Session, error) {
	parts := strings.SplitN(line, fieldSep, 5)
	if len(parts) < 4 {
		return core.Session{}, errMalformed
	}

	windows, _ := strconv.Atoi(parts[1])
	created := parseUnixTime(parts[2])
	attached := parts[3] == "1"

	dir := ""
	if len(parts) == 5 {
		dir = parts[4]
	}

	return core.Session{
		Name:      parts[0],
		Windows:   windows,
		Created:   created,
		Attached:  attached,
		Directory: dir,
	}, nil
}

func parseUnixTime(s string) time.Time {
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(ts, 0)
}

type parseError string

func (e parseError) Error() string { return string(e) }

const errMalformed = parseError("malformed session line")
