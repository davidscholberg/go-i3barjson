// Package i3bar provides an Go library for i3bar protocol support.
package i3bar

import (
    "encoding/json"
    "fmt"
    "io"
)

// Header represents the header of an i3bar message.
type Header struct {
    Version int         `json:"version"`
    Stop_signal int     `json:"stop_signal,omitempty"`
    Cont_signal int     `json:"cont_signal,omitempty"`
    Click_events bool   `json:"click_events,omitempty"`
}

// String pretty prints Header objects.
func (d *Header) String() string {
    str, err := json.MarshalIndent(d, "", "    ")
    if err != nil {
        return err.Error()
    }
    return string(str)
}

// Block represents a single block of an i3bar message.
type Block struct {
    Full_text string            `json:"full_text"`
    Short_text string           `json:"short_text,omitempty"`
    Color string                `json:"color,omitempty"`
    Min_width int               `json:"min_width,omitempty"`
    Align string                `json:"align,omitempty"`
    Name string                 `json:"name,omitempty"`
    Instance string             `json:"instance,omitempty"`
    Urgent bool                 `json:"urgent,omitempty"`
    Separator bool              `json:"separator"`
    Separator_block_width int   `json:"separator_block_width,omitempty"`
    Markup string               `json:"markup,omitempty"`
}

// String pretty prints Block objects.
func (d *Block) String() string {
    str, err := json.MarshalIndent(d, "", "    ")
    if err != nil {
        return err.Error()
    }
    return string(str)
}

// StatusLine represents a full i3bar status line.
type StatusLine []Block

// String pretty prints StatusLine objects.
func (d *StatusLine) String() string {
    str, err := json.MarshalIndent(d, "", "    ")
    if err != nil {
        return err.Error()
    }
    return string(str)
}

// Click represents an i3bar click event.
type Click struct {
    Name string     `json:"name"`
    Instance string `json:"instance"`
    X int           `json:"x"`
    Y int           `json:"y"`
    Button int      `json:"button"`
}

// String pretty prints Click objects.
func (d *Click) String() string {
    str, err := json.MarshalIndent(d, "", "    ")
    if err != nil {
        return err.Error()
    }
    return string(str)
}

func Init(h *Header, w io.Writer, r io.Reader, q chan bool) (chan StatusLine, error) {
    if (w == nil) {
        return nil, fmt.Errorf("error: Writer required")
    }
    var jsonWriter = json.NewEncoder(w)
    // TODO: implement read loop
    //var jsonReader *json.Decoder
    if (r != nil) {
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

func writeLoop(e *json.Encoder, c chan StatusLine, q chan bool) {
    for block := range c {
        // TODO: proper error handling
        err := e.Encode(block)
        if err != nil {
            fmt.Printf("%s\n", err)
        }
    }

    close(q)
}
