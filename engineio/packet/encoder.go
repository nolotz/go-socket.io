package packet

import (
	"io"

	"github.com/nolotz/go-socket.io/engineio/frame"
)

// FrameWriter is the writer which supports framing.
type FrameWriter interface {
	NextWriter(typ frame.Type) (io.WriteCloser, error)
}

type Encoder struct {
	w FrameWriter
}

func NewEncoder(w FrameWriter) *Encoder {
	return &Encoder{
		w: w,
	}
}

func (e *Encoder) NextWriter(ft frame.Type, pt Type) (io.WriteCloser, error) {
	w, err := e.w.NextWriter(ft)
	if err != nil {
		return nil, err
	}

	var b [1]byte
	if ft == frame.String {
		b[0] = pt.StringByte()
	} else {
		b[0] = pt.BinaryByte()
	}

	if _, err := w.Write(b[:]); err != nil {
		w.Close()
		return nil, err
	}

	return w, nil
}
