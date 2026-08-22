package table

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"endobit.io/table/sgr"
)

type cell struct {
	Text  string
	Value reflect.Value
}

type columnInfo struct {
	Labels    []string
	Width     int
	OmitEmpty bool
	IsZero    bool
	Skip      bool
}

type wrapper interface {
	Wrap() sgr.Wrapped
}

// FlushText flushes the Table data to its io.Writer as column aligned ANSI
// styled text. If the io.Writer is not a terminal no ANSI styles will be
// applied.
func (t *Table) FlushText() {
	var (
		prevType          reflect.Type
		columns           []columnInfo
		cells             [][]cell
		groupStart        int
		includeGroupStart = true
	)

	// This is the first pass through the table to determine the column widths
	// and cell contents. ANSI formatting is not part of this pass.

	for i := range t.rows {
		val := reflect.ValueOf(t.rows[i])
		currType := reflect.TypeOf(t.rows[i])

		if currType != prevType { // start a new table
			prevType = currType

			if columns != nil { // flush and reset for next table
				t.flush(columns, cells, t.annotationsFor(groupStart, i, includeGroupStart))
				cells = nil
				groupStart = i
				includeGroupStart = false
			}

			columns = t.processHeader(currType)
		}

		numFields := val.NumField()
		fields := make([]cell, numFields)

		for j := range numFields {
			value := val.Field(j)
			cell := cell{
				Text:  valueAsString(value), // cache it
				Value: value,
			}

			length := utf8.RuneCountInString(cell.Text)

			// If the value is a wrapper, use its Wrap() method to get the text
			// and its length.
			if value.CanInterface() {
				if a, ok := value.Interface().(wrapper); ok {
					w := a.Wrap()

					length = utf8.RuneCountInString(w.Text)
					if t.noColor {
						cell.Text = w.Text
					}
				}
			}

			if length > columns[j].Width {
				columns[j].Width = length
			}

			fields[j] = cell
		}

		cells = append(cells, fields)
	}

	t.flush(columns, cells, t.annotationsFor(groupStart, len(t.rows), includeGroupStart))
}

func (t *Table) flush(info []columnInfo, rows [][]cell, annotations []annotation) {
	t.flushHeader(info, rows)

	// This pass applies ANSI styles and prints the table rows.

	for i := range rows {
		annotations = t.flushAnnotations(annotations, i)

		var repeats []bool

		if i == 0 {
			repeats = make([]bool, len(rows[i]))
		} else {
			repeats = findRepeats(rows[i-1], rows[i])
		}

		for j := range rows[i] {
			if !info[j].visible() { // skip empty columns TODO: make this configurable
				continue
			}

			rowColor := t.colors.EvenRow
			if i%2 != 0 {
				rowColor = t.colors.OddRow
			}

			cell := rows[i][j]
			text := cell.Text

			if !t.noColor && cell.Value.CanInterface() {
				if a, ok := cell.Value.Interface().(wrapper); ok {
					text = a.Wrap().String()
				}
			}

			padding := strings.Repeat(" ", info[j].Width-utf8.RuneCountInString(cell.Text))

			switch {
			case cell.Text == "":
				text = sgr.Wrap(t.colors.Empty, strings.Repeat("-", info[j].Width)).String()
				padding = ""
			case repeats[j]:
				text = sgr.Wrap(t.colors.Repeat, text).String()
			}

			// Skip padding for the last column
			if j == lastVisibleColumn(info) {
				fmt.Fprint(t.writer, sgr.Wrap(rowColor, text))
			} else {
				fmt.Fprint(t.writer, sgr.Wrap(rowColor, text, padding), " ")
			}
		}

		fmt.Fprintln(t.writer)
	}

	t.flushAnnotations(annotations, len(rows))
}

func (t *Table) flushAnnotations(annotations []annotation, index int) []annotation {
	for len(annotations) > 0 && annotations[0].index == index {
		fmt.Fprintln(t.writer, sgr.Wrap(t.colors.Annotation, annotations[0].text))
		annotations = annotations[1:]
	}

	return annotations
}

func (t *Table) flushHeader(info []columnInfo, rows [][]cell) {
	var numLines int

	// header can have multiple lines (useful for specifying units)
	// column with the most lines determines the size of the header
	for _, c := range info {
		if len(c.Labels) > numLines {
			numLines = len(c.Labels)
		}
	}

	for j := range info {
		if info[j].OmitEmpty && isColumnZero(j, rows) {
			info[j].IsZero = true
		}
	}

	lastColumn := lastVisibleColumn(info)

	for i := range numLines {
		if i > 0 {
			fmt.Fprintln(t.writer)
		}

		for j := range info { // header
			var label string

			if !info[j].visible() {
				continue
			}

			if i < len(info[j].Labels) {
				label = info[j].Labels[i]
			}

			fmt.Fprint(t.writer, sgr.Wrapf(t.colors.Header, "%-*s", info[j].Width, label))

			if j != lastColumn {
				fmt.Fprint(t.writer, " ")
			}
		}
	}

	fmt.Fprintln(t.writer)
}

func lastVisibleColumn(info []columnInfo) int {
	for i := len(info) - 1; i >= 0; i-- {
		if info[i].visible() {
			return i
		}
	}

	return -1
}

func (c columnInfo) visible() bool {
	return !c.IsZero && !c.Skip
}

func (t *Table) annotationsFor(start, end int, includeStart bool) []annotation {
	var annotations []annotation

	for _, a := range t.annotations {
		if a.index > end || (!includeStart && a.index == start) || a.index < start {
			continue
		}

		a.index -= start
		annotations = append(annotations, a)
	}

	return annotations
}

func (t *Table) processHeader(header reflect.Type) []columnInfo {
	numFields := header.NumField()

	columns := make([]columnInfo, numFields)

	for i := range numFields {
		field := header.Field(i)
		label := t.fieldToLabel(field.Name)

		columns[i] = columnInfo{
			Labels: []string{label},
			Width:  utf8.RuneCountInString(label),
		}

		if tag := field.Tag.Get("table"); tag != "" {
			// Parse tag: "LABEL,omitempty" -> label="LABEL", omitEmpty=true
			label, options, _ := strings.Cut(tag, ",")

			if label == "-" {
				columns[i].Skip = true

				continue
			}

			if label != "" {
				labels := strings.Split(label, "\n")
				if labels[0] == "" { // use the default label if empty
					labels[0] = columns[i].Labels[0]
				}

				columns[i].Labels = labels
				columns[i].Width = maxStringLength(labels)
			}

			columns[i].OmitEmpty = strings.Contains(options, "omitempty")
		}
	}

	return columns
}

func findRepeats(top, bottom []cell) []bool {
	r := make([]bool, len(bottom))

	if top == nil || len(top) != len(bottom) {
		return r
	}

	for i := range bottom {
		r[i] = top[i].Text == bottom[i].Text
	}

	return r
}

func isColumnZero(n int, rows [][]cell) bool {
	for i := range rows {
		if !rows[i][n].Value.IsZero() {
			return false
		}
	}

	return true
}
