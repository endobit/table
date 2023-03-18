package table

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/endobit/table.git/sgr"
)

type cell struct {
	Text  string
	Value reflect.Value
}

type columnInfo struct {
	Label     string
	Width     int
	OmitEmpty bool
	IsZero    bool
}

type wrapper interface {
	Wrap() sgr.Wrapped
}

// FlushText flushes the Table data to its io.Writer as column aligned ANSI
// styled text. If the io.Writer is not a terminal no ANSI styles will be
// applied.
func (t *Table) FlushText() {
	var (
		prevType reflect.Type
		column   []columnInfo
	)

	rows := make([][]cell, 0, len(t.rows))

	for i := range t.rows {
		val := reflect.ValueOf(t.rows[i])
		numFields := val.NumField()
		currType := reflect.TypeOf(t.rows[i])

		if currType != prevType { // start a new table
			prevType = currType

			if column != nil {
				t.flush(column, rows)
				rows = nil // flush and reset for next table
			}

			column = make([]columnInfo, numFields)

			for i := 0; i < numFields; i++ {
				f := currType.Field(i)
				column[i] = columnInfo{Label: f.Name, Width: len(f.Name)}

				if tag := f.Tag.Get("table"); tag != "" {
					label, opts := parseTag(tag)

					column[i].Label = label
					column[i].Width = len(label)
					column[i].OmitEmpty = opts.Contains("omitempty")
				}
			}
		}

		row := make([]cell, numFields)

		for j := 0; j < numFields; j++ {
			v := val.Field(j)
			cell := cell{
				Text:  v.String(), // cache it
				Value: v,
			}

			length := len(cell.Text)

			if v.CanInterface() {
				if a, ok := v.Interface().(wrapper); ok {
					w := a.Wrap()
					length = len(w.Text)
					if t.noColor {
						cell.Text = w.Text
					}
				}
			}

			if length > column[j].Width {
				column[j].Width = length
			}

			row[j] = cell
		}

		rows = append(rows, row)
	}

	t.flush(column, rows)
}

func (t *Table) flush(c []columnInfo, r [][]cell) {
	t.flushHeader(c, r)

	for i := range r {
		var reps []bool

		if i > 0 {
			reps = repeats(r[i-1], r[i])
		} else {
			reps = repeats(nil, r[i])
		}

		for j := range r[i] {
			if c[j].IsZero {
				continue
			}

			isRepeat := reps[j]

			rowColor := t.colors.EvenRow
			if i%2 != 0 {
				rowColor = t.colors.OddRow
			}

			cell := r[i][j]
			pad := c[j].Width - len(cell.Text)
			s := cell.Text

			if !t.noColor && cell.Value.CanInterface() {
				if a, ok := cell.Value.Interface().(wrapper); ok {
					s = a.Wrap().String()
				}
			}

			switch {
			case cell.Text == "":
				s = sgr.Wrap(t.colors.Empty, strings.Repeat("-", len(c[j].Label))).String()
				pad = c[j].Width - len(c[j].Label)
			case isRepeat:
				s = sgr.Wrap(t.colors.Repeat, s).String()
			}

			spaces := ""
			if pad > 0 {
				spaces = strings.Repeat(" ", pad)
			}

			fmt.Fprint(t.writer, sgr.Wrap(rowColor, s, spaces))

			if j != len(r[i])-1 {
				fmt.Fprint(t.writer, " ")
			}
		}

		fmt.Fprintln(t.writer)
	}
}

func (t *Table) flushHeader(c []columnInfo, r [][]cell) {
	for i := range c { // header
		if c[i].OmitEmpty && isColumnZero(i, r) {
			c[i].IsZero = true
			continue
		}

		fmt.Fprint(t.writer, sgr.Wrapf(t.colors.Header, "%-*s", c[i].Width, c[i].Label))
		if i != len(c)-1 {
			fmt.Fprint(t.writer, " ")
		}
	}

	fmt.Fprintln(t.writer)
}

func repeats(top, bottom []cell) []bool {
	r := make([]bool, len(bottom))

	if top == nil || len(top) != len(bottom) {
		return r
	}

	for i := range bottom {
		if top[i].Text != bottom[i].Text {
			break
		}

		r[i] = true
	}

	return r
}

func isColumnZero(n int, r [][]cell) bool {
	for i := range r {
		if !r[i][n].Value.IsZero() {
			return false
		}
	}

	return true
}
