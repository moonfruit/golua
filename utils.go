package golua

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func (s State) PrintStack() error {
	return s.PrintStackf(os.Stdout)
}

func (s State) PrintStackf(w io.Writer) error {
	if !s.CheckStack(1) {
		return errors.New("stack has not enough space")
	}

	if w == nil {
		w = os.Stdout
	}
	w = wrapWriter(w)

	for i := s.GetTop(); i > 0; i-- {
		fmt.Fprintf(w, "|-- %02d: type[%v] value[%s]", i, s.Type(i), s.ToString(i))
	}

	main := s.PushThread()
	fmt.Fprintf(w, "`-- 00: type[%v] value[%s] main[%v]", s.Type(-1), s.ToString(-1), main)
	s.Pop(1)

	return nil
}

func wrapWriter(w io.Writer) io.Writer {
	if lw, ok := w.(LineWriter); ok {
		return lineWriter{lw}
	}
	return proxyWriter{w}
}

type LineWriter interface {
	WriteLine(p []byte) (n int, err error)
}

type lineWriter struct {
	LineWriter
}

func (w lineWriter) Write(p []byte) (int, error) {
	return w.WriteLine(p)
}

type proxyWriter struct {
	w io.Writer
}

func (w proxyWriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	if err != nil {
		return n, err
	}

	n2, err := w.w.Write([]byte{'\n'})
	return n + n2, err
}
