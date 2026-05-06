package telemetry

import (
	"fmt"
	"io"
	"strings"

	"github.com/elias/axiom/engine"
)

type Export struct {
	Keys   []string
	Values map[string]string
}

func NewExport() *Export {
	return &Export{
		Values: make(map[string]string),
	}
}
func (e *Export) Add(key string, value string) {
	e.Keys = append(e.Keys, key)
	e.Values[key] = value
}

type exportable interface {
	ExportFields() *Export
}

type TelemetryWriter struct {
	writer io.Writer
	tick   *engine.Tick
}

func NewTelemetryWriter(w io.Writer, t *engine.Tick) *TelemetryWriter {
	return &TelemetryWriter{
		writer: w,
		tick:   t,
	}
}

func (w *TelemetryWriter) Write(e *Export) {
	var vals []string
	for _, k := range e.Keys {
		vals = append(vals, fmt.Sprintf("%v", e.Values[k]))
	}

	fmt.Fprintf(w.writer, "%d,%s\n", w.tick.Tick(), strings.Join(vals, ","))
}
