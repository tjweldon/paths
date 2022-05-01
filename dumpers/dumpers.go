package dumpers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"tjweldon/paths/fullpath"
)

type Format func(p *fullpath.Paths) (string, error)

type Formatter struct {
	Format Format
}

// Formatter implementations

// Json formats the paths as a json array
func Json() Formatter {
	return Formatter{
		Format: func(p *fullpath.Paths) (string, error) {
			jsonBytes, err := json.MarshalIndent(p, "", strings.Repeat(" ", 4))
			if err != nil {
				return "", err
			}
			return string(jsonBytes), nil
		},
	}
}

// IndexedList prints formats the paths as a list with each entry indexed i.e.
//      0: /first/path
//      1: /next/path
//      ...
//      n: /last/path
//
func IndexedList() Formatter {
	return Formatter{
		Format: func(p *fullpath.Paths) (string, error) {
			result := ""
			for index, path := range p.Paths {
				result += fmt.Sprintf("%d: %s\n", index, path)
			}
			return result, nil
		},
	}
}

// ExportCommand formats the path as a shell export command for the deduped path set in config
func ExportCommand() Formatter {
	return Formatter{
		Format: func(p *fullpath.Paths) (string, error) {
			return fmt.Sprintf(
					"export PATH='%s'",
					strings.Join(p.Paths, ":"),
				),
				nil
		},
	}
}

type Output func(s string) error

type Outputter struct {
	Output Output
}

// Stdout dumps any string to stdout
func Stdout() Outputter {
	return Outputter{
		Output: func(s string) error {
			fmt.Println(s)
			return nil
		},
	}
}

// FileOverwrite dumps any string to the file at the path given
func FileOverwrite(path string) Outputter {
	return Outputter{
		Output: func(s string) error {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
			defer func() { _ = file.Close() }()
			if err != nil {
				return err
			}
			_, err = file.Write([]byte(s))
			return err
		},
	}
}

// FileAppend is an Outputter that will append the supplied string to the file at the path supplied.
func FileAppend(path string) Outputter {
	return Outputter{
		Output: func(s string) error {
			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			defer func() { _ = file.Close() }()
			if err != nil {
				return err
			}
			_, err = file.Write([]byte("\n" + s))
			return err
		},
	}
}

type Dumper struct {
	Formatter  Formatter
	Outputters []Outputter
}

func MakeDumper(formatter Formatter, outputters ...Outputter) Dumper {
	return Dumper{formatter, outputters}
}

func (d Dumper) Dump(p *fullpath.Paths) error {
	formatted, err := d.Formatter.Format(p)
	if err != nil {
		return err
	}

	for _, out := range d.Outputters {
		err = out.Output(formatted)
		if err != nil {
			return err
		}
	}

	return nil
}

type MultiDumper struct {
	dumpers []Dumper
}

func NewMulti() *MultiDumper {
	return &MultiDumper{dumpers: []Dumper{}}
}

func (m *MultiDumper) AddDumper(formatter Formatter, outputters ...Outputter) *MultiDumper {
	m.dumpers = append(m.dumpers, MakeDumper(formatter, outputters...))
	return m
}

func (m *MultiDumper) Dump(p *fullpath.Paths) (err error) {
	for _, dumper := range m.dumpers {
		err = dumper.Dump(p)
		if err != nil {
			return err
		}
	}

	return err
}
