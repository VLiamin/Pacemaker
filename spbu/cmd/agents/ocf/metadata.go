package ocf

import (
	"encoding/xml"
	"time"
)

type (
	MetaData struct {
		Name            string     `xml:"name,attr"`
		Version         string     `xml:"version"`
		LongDesription  string     `xml:"longdesc"`
		ShortDesription string     `xml:"shortdesc"`
		Parameters      Parameters `xml:"parameters>parameter"`
		Actions         Actions    `xml:"actions>action"`
	}
)

func (md MetaData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Tmp MetaData
	return e.EncodeElement(struct {
		Tmp
		VersionAttr string `xml:"version,attr"`
	}{Tmp(md), md.Version}, xml.StartElement{Name: xml.Name{Local: `resource-agent`}})
}

type (
	Parameters []Parameter
	Parameter  struct {
		Name            string
		Uniq            bool
		Required        bool
		Kind            string
		Default         string
		LongDesription  string
		ShortDesription string
	}
)

func oneOrZero(v bool) int {
	if v {
		return 1
	}
	return 0
}

func (p Parameter) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type descr struct {
		Lang  string `xml:"lang,attr"`
		Value string `xml:",chardata"`
	}
	type content struct {
		Kind    string `xml:"type,attr"`
		Default string `xml:"default,attr,omitempty"`
	}
	descFor := func(name, descr string) string {
		if len(descr) > 0 {
			return descr
		}
		return `parameter ` + name
	}
	return e.EncodeElement(struct {
		Name            string  `xml:"name,attr"`
		Uniq            int     `xml:"unique,attr"`
		Required        int     `xml:"required,attr"`
		LongDesription  descr   `xml:"longdesc"`
		ShortDesription descr   `xml:"shortdesc"`
		Content         content `xml:"content"`
	}{
		p.Name,
		oneOrZero(p.Uniq),
		oneOrZero(p.Required),
		descr{`en`, descFor(p.Name, p.LongDesription)},
		descr{`en`, descFor(p.Name, p.ShortDesription)},
		content{p.Kind, p.Default},
	}, xml.StartElement{Name: xml.Name{Local: `parameter`}})
}

type (
	Actions []Action
	Action  struct {
		Name     string
		Interval time.Duration
		Timeout  time.Duration
	}
)

func (a Action) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		Name     string `xml:"name,attr"`
		Interval int    `xml:"interval,attr"`
		Timeout  int    `xml:"timeout,attr"`
	}{
		a.Name,
		int(a.Interval / time.Millisecond),
		int(a.Timeout / time.Millisecond),
	}, xml.StartElement{Name: xml.Name{Local: `action`}})
}
