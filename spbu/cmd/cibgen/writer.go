package main

import (
	encoding "encoding/xml"
	"fmt"
	"io"
	"log"
)

type (
	ioWriterStackElement struct {
		local       string
		hasChildren bool
	}
	XmlWriter struct {
		io.Writer
		stack []ioWriterStackElement
	}
)

func (w *XmlWriter) MustStartElement(local string) {
	err := w.StartElement(local)
	if err != nil {
		log.Fatalf(`failed to write start element %s. %v`, local, err)
	}
}

func (w *XmlWriter) MustEndElement(local string) {
	err := w.EndElement(local)
	if err != nil {
		log.Fatalf(`failed to write end element %s. %v`, local, err)
	}
}

func (w *XmlWriter) MustAttributes(attributes [][2]string) {
	for _, attribute := range attributes {
		w.MustAttribute(attribute[0], attribute[1])
	}
}

func (w *XmlWriter) MustAttribute(local, value string) {
	err := w.Attribute(local, value)
	if err != nil {
		log.Fatalf(`failed to write attribute %s=%s. %v`, local, value, err)
	}
}

func (w *XmlWriter) StartElement(local string) (err error) {
	if completePreviousTag := len(w.stack) > 0 && !w.stack[len(w.stack)-1].hasChildren; completePreviousTag {
		if _, err = io.WriteString(w, `>`); err != nil {
			return
		}
		if _, err = io.WriteString(w, "\n"); err != nil {
			return
		}
	}
	for range w.stack {
		if _, err = io.WriteString(w, `  `); err != nil {
			return
		}
	}
	if _, err = io.WriteString(w, `<`); err != nil {
		return
	}
	if _, err = io.WriteString(w, local); err != nil {
		return
	}
	if len(w.stack) > 0 {
		w.stack[len(w.stack)-1].hasChildren = true
	}
	w.stack = append(w.stack, ioWriterStackElement{local, false})
	return
}
func (w *XmlWriter) Attribute(local, value string) (err error) {
	if _, err = io.WriteString(w, ` `); err != nil {
		return
	}
	if _, err = io.WriteString(w, local); err != nil {
		return
	}
	if _, err = io.WriteString(w, `="`); err != nil {
		return
	}
	if err = encoding.EscapeText(w, []byte(value)); err != nil {
		return
	}
	if _, err = io.WriteString(w, `"`); err != nil {
		return
	}
	return
}
func (w *XmlWriter) EndElement(local string) (err error) {
	current := w.stack[len(w.stack)-1]
	if local != current.local {
		return fmt.Errorf("cannot close %s before %s", local, current.local)
	}
	newStack := w.stack[:len(w.stack)-1]
	if current.hasChildren {
		for range newStack {
			if _, err = io.WriteString(w, `  `); err != nil {
				return
			}
		}
		if _, err = io.WriteString(w, `</`); err != nil {
			return
		}
		if _, err = io.WriteString(w, current.local); err != nil {
			return
		}
		if _, err = io.WriteString(w, `>`); err != nil {
			return
		}
	} else {
		if _, err = io.WriteString(w, `/>`); err != nil {
			return
		}
	}
	if _, err = io.WriteString(w, "\n"); err != nil {
		return
	}
	w.stack = newStack
	return
}
