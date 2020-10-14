package typescript

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Code represents an item of code that can be rendered.
type Code interface {
	render(f *File, w io.Writer, s *Statement) error
	isNull(f *File) bool
}

// Save renders the file and saves to the filename provided.
func (f *File) Save(filename string) error {
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

// AppendAfter renders the file and append data to the filename provided.
// Finds the last match of string *last* and the next first match of string *first*
// appends data after *first*
func (f *File) AppendAfter(filename string, last string, first string) error {
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		return err
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if lastMatch := bytes.LastIndex(content, []byte(last)); lastMatch != -1 {
		if firstMatch := bytes.Index(content[lastMatch:], []byte(first)); firstMatch != -1 {
			tmp := []byte{}
			firstMatch += lastMatch + len(first)
			tmp = append(tmp, content[:firstMatch]...)
			tmp = append(tmp, buf.Bytes()...)
			tmp = append(tmp, content[firstMatch:]...)

			if _, err = file.Write(tmp); err != nil {
				return err
			}
			return nil
		}
	}

	if _, err = file.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (f *File) CheckRepetition(filename string, names []string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	for _, name := range names {
		if bytes.Contains(content, []byte(name)) {
			err = fmt.Errorf("Collision detected: %s already presented in %s file,"+
				" please delete this file and run go generate ./... from project root folder ",
				name, filename)
			return err
		}
	}

	return nil
}

// Render renders the file to the provided writer.
func (f *File) Render(w io.Writer) error {
	body := &bytes.Buffer{}
	if err := f.render(f, body, nil); err != nil {
		return err
	}
	source := &bytes.Buffer{}
	if len(f.headers) > 0 {
		for _, c := range f.headers {
			if err := Comment(c).render(f, source, nil); err != nil {
				return err
			}
			if _, err := fmt.Fprint(source, "\n"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(source, "\n"); err != nil {
			return err
		}
	}
	for _, c := range f.comments {
		if err := Comment(c).render(f, source, nil); err != nil {
			return err
		}
		if _, err := fmt.Fprint(source, "\n"); err != nil {
			return err
		}
	}
	if err := f.renderImports(source); err != nil {
		return err
	}
	if _, err := source.Write(body.Bytes()); err != nil {
		return err
	}
	if _, err := w.Write(source.Bytes()); err != nil {
		return err
	}
	return nil
}

func (f *File) renderImports(source io.Writer) (err error) {

	// import { JSONRPCRequest, IResponse } from '@mihanizm56/fetch-api';
	for from, def := range f.imports {
		if _, err := fmt.Fprintf(source, "import {%s} from '%s';\n", strings.Join(def.items, ","), from); err != nil {
			return err
		}
	}
	return
}
