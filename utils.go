package golua

import (
	"fmt"
	"io"
	"os"
)

func (s State) PrintStack() error {
	return s.PrintStackf(os.Stdout)
}

func (s State) PrintStackl(l Logger) error {
	return s.PrintStackf(logWriter{l})
}

func (s State) PrintStackf(w io.Writer) error {
	s.CheckStack(2)

	if w == nil {
		w = os.Stdout
	}
	w = wrapWriter(w)

	for i := s.GetTop(); i > 0; i-- {
		fmt.Fprintf(w, "|-- %02d: type<%v> value=`%s`", i, s.Type(i), s.ToString(i))
		s.Pop(1)
	}

	main := s.PushThread()
	fmt.Fprintf(w, "`-- 00: type<%v> main<%v> value=`%s`", s.Type(-1), main, s.ToString(-1))
	s.Pop(2)

	return nil
}

func wrapWriter(w io.Writer) io.Writer {
	if lw, ok := w.(LineWriter); ok {
		return lineWriter{lw}
	}
	return proxyWriter{w}
}

type LineWriter interface {
	io.Writer
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

type Logger interface {
	Log(args ...interface{})
}

type logWriter struct {
	l Logger
}

func (w logWriter) Write(p []byte) (n int, err error) {
	panic("Not support")
}

func (w logWriter) WriteLine(p []byte) (n int, err error) {
	w.l.Log(string(p))
	return len(p), nil
}
