package generator

import (
	"bytes"
	"text/template"
)

// Translation - is struct holding data for one translation to be printed to the buffer
type Translation struct {
	Comment   string
	Key       string
	Value     string
	Delimeter string
}

// FileGenerator prints a line for a translation using selected implementation
type FileGenerator interface {
	WriteFirstLine(buffer *bytes.Buffer)
	WriteLineTo(buffer *bytes.Buffer, translation *Translation)
	WriteLastLine(buffer *bytes.Buffer)
}

var tmplProperies = template.Must(template.New("").Parse(`
{{if .Comment}} # {{.Comment}}
{{end}} {{.Key}}{{.Delimeter}} {{.Value}}`))

var tmplJSON = template.Must(template.New("").Parse(
	`
	"{{.Key}}": {
		{{end}}"one": "{{.Value}}"
	},
	`))

const (
	typeProperties = "properties"
	typeJSON       = "json"
)

type (
	propertiesImpl struct{}
	jsonImpl       struct{}
)

// GetGenerator returns instance of generator
func GetGenerator(outputFileType string) FileGenerator {

	switch outputFileType {
	case typeProperties:
		return propertiesImpl{}

	case typeJSON:
		return jsonImpl{}

	default:
		return propertiesImpl{}
	}
}

// .properties implementation

func (p propertiesImpl) WriteLineTo(buffer *bytes.Buffer, translation *Translation) {
	tmplProperies.Execute(buffer, translation)
}

func (p propertiesImpl) WriteFirstLine(buffer *bytes.Buffer) {}

func (p propertiesImpl) WriteLastLine(buffer *bytes.Buffer) {}

// JSON implementation

func (j jsonImpl) WriteLineTo(buffer *bytes.Buffer, translation *Translation) {
	tmplJSON.Execute(buffer, translation)
}

func (j jsonImpl) WriteFirstLine(buffer *bytes.Buffer) { buffer.WriteString("{\n") }

func (j jsonImpl) WriteLastLine(buffer *bytes.Buffer) {

	offset := buffer.Len() - 3
	if offset > 0 {
		// remove the last comma to keep JSON valid
		buffer.Truncate(offset)
	}

	buffer.WriteString("\n}")
}
