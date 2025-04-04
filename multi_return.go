package fuego

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// DataOrTemplate is a struct that can return either data or a template
// depending on the asked type.
type DataOrTemplate[T any] struct {
	Data     T
	Template any
}

var (
	_ CtxRenderer    = DataOrTemplate[any]{} // Can render HTML (template)
	_ json.Marshaler = DataOrTemplate[any]{} // Can render JSON (data)
	_ xml.Marshaler  = DataOrTemplate[any]{} // Can render XML (data)
	_ yaml.Marshaler = DataOrTemplate[any]{} // Can render YAML (data)
	_ fmt.Stringer   = DataOrTemplate[any]{} // Can render string (data)
)

func (m DataOrTemplate[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Data)
}

func (m DataOrTemplate[T]) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	return e.Encode(m.Data)
}

func (m DataOrTemplate[T]) MarshalYAML() (any, error) {
	return m.Data, nil
}

func (m DataOrTemplate[T]) String() string {
	return fmt.Sprintf("%v", m.Data)
}

func (m DataOrTemplate[T]) Render(c context.Context, w io.Writer) error {
	switch renderer := m.Template.(type) {
	case CtxRenderer:
		return renderer.Render(c, w)
	case Renderer:
		return renderer.Render(w)
	default:
		panic("template must be either CtxRenderer or Renderer")
	}
}

// DataOrHTML is a helper function to create a [DataOrTemplate] return item without specifying the type.
func DataOrHTML[T any](data T, template any) *DataOrTemplate[T] {
	return &DataOrTemplate[T]{
		Data:     data,
		Template: template,
	}
}
