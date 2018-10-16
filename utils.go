package golua

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func (s *State) PrintStack() error {
	return s.PrintStackf(os.Stdout)
}

var newline = []byte{'\n'}

func (s *State) PrintStackf(w io.Writer) error {
	return s.RawPrintStack(func(f string, a ...interface{}) error {
		if _, err := fmt.Fprintf(w, f, a...); err != nil {
			return err
		}
		_, err := w.Write(newline)
		return err
	})
}

type Logger interface {
	Logf(format string, args ...interface{})
}

func (s *State) PrintStackl(l Logger) error {
	return s.RawPrintStack(func(f string, a ...interface{}) error {
		l.Logf(f, a...)
		return nil
	})
}

func (s *State) RawPrintStack(printf func(string, ...interface{}) error) error {
	if !s.RawCheckStack(2) {
		return ErrMem
	}

	for i := s.GetTop(); i > 0; i-- {
		value := s.ToString(i)
		s.Pop(1)

		if err := printf("|-- %02d: type<%v> value=`%s`", i, s.Type(i), value); err != nil {
			return err
		}
	}

	main := s.PushThread()

	ty := s.Type(-1)
	value := s.ToString(-1)
	s.Pop(2)

	if err := printf("`-- 00: type<%v> main<%v> value=`%s`", ty, main, value); err != nil {
		return err
	}

	return nil
}

type skipReader struct {
	r    *bufio.Reader
	skip bool
}

func NewSkipReader(r io.Reader) io.Reader {
	var br *bufio.Reader
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	return &skipReader{r: br}
}

var bom = []byte{0xEF, 0xBB, 0xBF}

func (r *skipReader) Read(p []byte) (n int, err error) {
	if !r.skip {
		// Skip BOM
		data, _ := r.r.Peek(3)
		if bytes.Equal(data, bom) {
			r.r.Discard(3)
		}

		// Skip Comment
		data, _ = r.r.Peek(1)
		if len(data) > 0 && data[0] == '#' {
			var b byte
			for err == nil && b != '\n' {
				b, err = r.r.ReadByte()
			}
		}

		r.skip = true
	}

	return r.r.Read(p)
}
