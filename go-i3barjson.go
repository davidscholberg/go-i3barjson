// Package i3barjson provides a Go library for i3bar JSON protocol support.
package i3barjson

import (
	"encoding/json"
	"fmt"
	"io"
)

var jsonWriter jsonArrayEncoder

// marshalIndent returns a marshalled JSON string of the given object.
// NOTE: if the marshal call fails, then a non-JSON error string is returned.
func marshalIndent(d interface{}) string {
	str, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return err.Error()
	}
	return string(str)
}

// Header represents the header of an i3bar message.
type Header struct {
	Version     int  `json:"version"`
	StopSignal  int  `json:"stop_signal,omitempty"`
	ContSignal  int  `json:"cont_signal,omitempty"`
	ClickEvents bool `json:"click_events,omitempty"`
}

// String pretty prints Header objects.
func (d Header) String() string {
	return marshalIndent(d)
}

// Block represents a single block of an i3bar message.
type Block struct {
	FullText            string `json:"full_text"`
	ShortText           string `json:"short_text,omitempty"`
	Color               string `json:"color,omitempty"`
	MinWidth            string `json:"min_width,omitempty"`
	Align               string `json:"align,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Separator           bool   `json:"separator"`
	SeparatorBlockWidth int    `json:"separator_block_width,omitempty"`
	Markup              string `json:"markup,omitempty"`
}

// String pretty prints Block objects.
func (d Block) String() string {
	return marshalIndent(d)
}

// StatusLine represents a full i3bar status line.
type StatusLine []*Block

// String pretty prints StatusLine objects.
func (d StatusLine) String() string {
	return marshalIndent(d)
}

// Click represents an i3bar click event.
type Click struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Button   int    `json:"button"`
}

// String pretty prints Click objects.
func (d Click) String() string {
	return marshalIndent(d)
}

// jsonArrayEncoder is an object that streams an infinite JSON array.
type jsonArrayEncoder struct {
	count int
	w     io.Writer
	e     json.Encoder
}

// Encode streams an infinite JSON array.
// Each call adds another element to the array.
func (e *jsonArrayEncoder) Encode(v interface{}) error {
	linePrefix := ","
	if e.count == 0 {
		linePrefix = "["
		e.count++
	}
	_, err := e.w.Write([]byte(linePrefix))
	if err != nil {
		return err
	}

	return e.e.Encode(v)
}

// newJsonArrayEncoder returns a new jsonArrayEncoder that wraps w.
func newJsonArrayEncoder(w io.Writer) jsonArrayEncoder {
	return jsonArrayEncoder{0, w, *json.NewEncoder(w)}
}

// Init initializes the i3bar io.
// w is the io.Writer to write to (usually os.Stdout).
// r is the io.Reader to from (usually os.Stdin) (TODO: implement).
// h is the Header object to send as the first line to i3bar.
func Init(w io.Writer, r io.Reader, h Header) error {
	if w == nil {
		return fmt.Errorf("error: Writer required")
	}
	jsonWriter = newJsonArrayEncoder(w)
	// TODO: implement read loop
	//var jsonReader *json.Decoder
	if r != nil {
		//jsonReader = json.NewDecoder(r)
	}

	msg, err := json.Marshal(h)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(jsonWriter.w, string(msg))
	if err != nil {
		return err
	}

	return nil
}

// Update sends a new StatusLine to i3bar.
// NOTE: this function is not thread safe. Only call it from a single thread.
func Update(s StatusLine) error {
	return jsonWriter.Encode(s)
}
