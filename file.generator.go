package main

import (
	"bytes"
)

// FileGenerator interface
type FileGenerator interface {
	WriteFirstLine(buffer *bytes.Buffer)
	WriteLineTo(buffer *bytes.Buffer, code string, translation string, comment string, delimeter string)
	WriteLastLine(buffer *bytes.Buffer)
}

type propertiesImpl struct{}
type jsonImpl struct{}

// GetGenerator returns instance of generator
func GetGenerator(outputFileType string) FileGenerator {

	switch outputFileType {
	case "properties":
		return propertiesImpl{}

	case "json":
		return jsonImpl{}

	default:
		return propertiesImpl{}
	}
}

func (p propertiesImpl) WriteLineTo(buffer *bytes.Buffer, code string, translation string, comment string, delimeter string) {

	if len(comment) > 0 {
		buffer.WriteString("# " + comment + "\n")
	}

	buffer.WriteString(code + delimeter + translation + "\n")
}

func (p propertiesImpl) WriteFirstLine(buffer *bytes.Buffer) {}
func (p propertiesImpl) WriteLastLine(buffer *bytes.Buffer)  {}

func (j jsonImpl) WriteLineTo(buffer *bytes.Buffer, code string, translation string, comment string, delimeter string) {

	buffer.WriteString("\"" + code + "\": {\n")

	if len(comment) > 0 {
		buffer.WriteString("  \"_comment\" : \"" + comment + "\",\n")
	}

	buffer.WriteString("  \"one\" : \"" + translation + "\"\n } \n")
}
func (j jsonImpl) WriteFirstLine(buffer *bytes.Buffer) {}
func (j jsonImpl) WriteLastLine(buffer *bytes.Buffer)  {}
