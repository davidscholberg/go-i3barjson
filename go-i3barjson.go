// Package i3barjson provides a Go library for i3bar JSON protocol support.
package i3barjson

import (
	"encoding/json"
	"fmt"
	"io"
)

// marshalIndent returns a marshalled JSON string of the given object.
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
func (d *Header) String() string {
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
func (d *Block) String() string {
	return marshalIndent(d)
}

// StatusLine represents a full i3bar status line.
type StatusLine []*Block

// String pretty prints StatusLine objects.
func (d *StatusLine) String() string {
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
func (d *Click) String() string {
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

	err = e.e.Encode(v)
	if err != nil {
		return err
	}

	return nil
}

// newJsonArrayEncoder returns a new jsonArrayEncoder that wraps w.
func newJsonArrayEncoder(w io.Writer) *jsonArrayEncoder {
	return &jsonArrayEncoder{0, w, *json.NewEncoder(w)}
}

// Init initializes the i3bar io.
// w is the io.Writer to write to (usually os.Stdout).
// r is the io.Reader to from (usually os.Stdin) (TODO: implement).
// q is a channel that will be closed when the write loop is finished.
// The returned channel can be used to write status lines to w.
func Init(h *Header, w io.Writer, r io.Reader, q chan bool) (chan StatusLine, error) {
	if w == nil {
		return nil, fmt.Errorf("error: Writer required")
	}
	var jsonWriter = newJsonArrayEncoder(w)
	// TODO: implement read loop
	//var jsonReader *json.Decoder
	if r != nil {
		//jsonReader = json.NewDecoder(r)
	}

	msg, err := json.Marshal(h)
	if err != nil {
		return nil, fmt.Errorf("error: couldn't parse Header")
	}
	_, err = fmt.Fprintln(w, string(msg))
	if err != nil {
		return nil, fmt.Errorf("error: couldn't write to Writer")
	}

	statusChan := make(chan StatusLine)
	go writeLoop(jsonWriter, statusChan, q)
	return statusChan, nil
}

// writeLoop continuosly writes status lines sent over c to e.
// q is closed when there are no more values to read from c.
func writeLoop(e *jsonArrayEncoder, c chan StatusLine, q chan bool) {
	for block := range c {
		// TODO: proper error handling
		err := e.Encode(block)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}

	close(q)
}
