package gopher

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Dirent struct {
	ItemType ItemType `json:"type"`
	Display  string   `json:"display"`
	Selector string   `json:"selector"`
	Hostname string   `json:"host,omitempty"`
	Port     string   `json:"port,omitempty"`
	Plus     bool     `json:"plus,omitempty"`

	Raw string `json:"-"`
}

var zeroDirent Dirent

func (d *Dirent) write(w *bufio.Writer) error {
	w.WriteByte(byte(d.ItemType))
	w.WriteString(d.Display)
	w.WriteByte('\t')
	w.WriteString(d.Selector)
	w.WriteByte('\t')
	w.WriteString(d.Hostname)
	w.WriteByte('\t')
	w.WriteString(d.Port)
	if d.Plus {
		w.WriteByte('\t')
		w.WriteByte('+')
	}
	_, err := w.Write(crlf)
	return err
}

func (d *Dirent) URL() URL {
	var u = URL{Scheme: "gopher"}
	d.PopulateURL(&u)
	return u
}

func (d *Dirent) PopulateURL(u *URL) {
	u.ItemType = d.ItemType
	u.Selector = d.Selector
	u.Hostname = d.Hostname
	u.Port = d.Port
}

type direntFlag int

const (
	direntHostOptional = 1 << iota
	direntNoValidatePort
)

func parseDirent(txt string, line int, dir *Dirent, flag direntFlag) error {
	if len(txt) == 0 {
		return fmt.Errorf("gopher: empty dirent at line %d", line)
	}

	tsz := len(txt)

	*dir = zeroDirent
	dir.ItemType = ItemType(txt[0])
	dir.Raw = txt

	start := 1
	field := 0

	for i := start; i <= tsz; i++ {
		if i == tsz || txt[i] == '\t' {
			switch field {
			case 0:
				dir.Display = txt[start:i]
				field, start = field+1, i+1
			case 1:
				dir.Selector = txt[start:i]
				field, start = field+1, i+1
			case 2:
				dir.Hostname = txt[start:i]
				field, start = field+1, i+1

			case 3:
				// XXX: Things can get a bit fouled up by bad servers; telefisk.org serves mail
				// archives with bad whitespace in 'i' lines:
				// gopher://telefisk.org/1/mailarchives/gopher/gopher-2014-12.mbox%3F133
				//
				// If we can accept the server's output without doing something
				// unreasonable, we should try, so let's chop whitespace and skip empty
				// strings. Some hosts will fill these fields out with dummy data, so we
				// can't just presume that 'i' means concatenate all fields together and
				// presume that's the line; I think telefisk.org is just serving files up
				// as-is and prepending 'i' to every line regardless of whether that's
				// valid.
				ps := strings.TrimSpace(txt[start:i])
				if ps != "" {
					if flag&direntNoValidatePort == 0 {
						if _, err := strconv.ParseInt(ps, 10, 16); err != nil {
							return fmt.Errorf("gopher: unexpected port %q at line %d: %w", ps, line, err)
						}
					}
					dir.Port = ps
				}
				field, start = field+1, i+1

			case 4:
				ps := txt[start:i]
				if ps == "+" {
					dir.Plus = true
				} else if ps != "" {
					return fmt.Errorf("gopher: unexpected 'plus' field at line %d; expected '+' or '', found %q", line, ps)
				}

			case 5:
				return fmt.Errorf("gopher: extra fields at line %d: %q", line, txt[start:i])
			}
		}
	}

	fieldLimit := 4
	if flag&direntHostOptional != 0 {
		fieldLimit = 2
	}
	if dir.ItemType == Info || dir.ItemType == ItemError {
		// XXX: Lots of servers don't fill out the extra fields for 'i' lines.
		// Some don't fill it out for '3' lines, either.
		fieldLimit = 1
	}

	if field < fieldLimit {
		return fmt.Errorf("gopher: missing fields at line %d: %q", line, txt)
	}

	return nil
}
