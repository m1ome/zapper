package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"sort"
)

type zapper struct {
	reader *bufio.Reader
	writer io.Writer
}

type entry struct {
	level   string
	message string
	caller  string
	ts      string
	trace   []string
	fields  map[string]interface{}
}

type options struct {
	short bool
}

func (z *zapper) pipe() error {
	var stop bool

	for !stop {
		input, err := z.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				stop = true
			} else {
				fmt.Fprintf(z.writer, "Error reading from pipe: %s\n", err)
				return err
			}
		}

		entry := &entry{
			fields: make(map[string]interface{}),
		}

		// Parsing entry object
		err = jsonparser.ObjectEach(input, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			switch string(key) {
			case "level":
				// Binding level
				entry.level = string(value)
			case "msg":
				// Binding message
				entry.message = string(value)
			case "ts":
				// Binding timestamp
				ts, err := strconv.ParseFloat(string(value), 64)
				if err != nil {
					return err
				}

				entry.ts = time.Unix(int64(ts), 0).Format("2006-01-02T15:04:05.999Z")
			case "caller":
				// Binding caller
				entry.caller = string(value)
			case "stacktrace":
				// Binding trace
				trace := string(value)
				if len(trace) > 0 {
					entry.trace = strings.Split(trace, `\n`)
				}
			default:
				entry.fields[string(key)] = string(value)
			}

			return nil
		})

		if err != nil {
			fmt.Fprintf(z.writer, "Error: %s\n", err)
			continue
		}

		z.write(entry)
	}

	return nil
}

func (z *zapper) write(e *entry) {
	tsColor := color.New(color.Faint)
	ts := tsColor.Sprint(e.ts)

	level := "[%s]"
	switch e.level {
	case "debug":
		level = color.New(color.Faint, color.Bold).Sprintf(level, "DEBG")
	case "info":
		level = color.New(color.Bold, color.FgGreen).Sprintf(level, "INFO")
	case "warn":
		level = color.New(color.Bold, color.FgYellow).Sprintf(level, "WARN")
	case "error":
		level = color.New(color.Bold, color.FgRed).Sprintf(level, "ERRO")
	case "dpanic":
		level = color.New(color.Bold, color.FgHiRed).Sprintf(level, "DPAN")
	case "panic":
		level = color.New(color.Bold, color.FgHiRed, color.BlinkSlow).Sprintf(level, "PANC")
	case "fatal":
		level = color.New(color.Bold, color.FgHiRed, color.BlinkRapid).Sprintf(level, "FATL")
	}

	caller := ""
	if len(e.caller) > 0 {
		caller = color.New(color.Faint, color.Italic).Sprintf("@ %s", e.caller)
	}

	fmt.Fprintf(z.writer, "%s %s: %s %s\n", ts, level, e.message, caller)

	if len(e.fields) > 0 {
		fmt.Fprintf(z.writer, "\t%s", color.New(color.Bold, color.FgGreen).Sprint("Fields:\n"))
		keys := make([]string, 0, len(e.fields))
		for key := range e.fields {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(z.writer, "\t\t%s:%v\n", color.New(color.Bold).Sprint(k), color.New(color.Italic).Sprint(e.fields[k]))
		}
	}

	if len(e.trace) > 0 {
		fmt.Fprintf(z.writer, "\t%s", color.New(color.Bold, color.FgGreen).Sprint("Stacktrace:\n"))
		for _, tr := range e.trace {
			fmt.Fprintf(z.writer, "\t\t%s\n", color.New(color.Faint).Sprint(tr))
		}
	}
}
