package golang

import (
	"io"
	"net"
	"time"
)

type StdioConn struct {
	r io.Reader
	w io.Writer
}

func NewStdioConn(r io.Reader, w io.Writer) StdioConn {
	return StdioConn{
		r: r,
		w: w,
	}
}

// Close implements net.Conn.
func (s StdioConn) Close() error {
	return nil
}

// LocalAddr implements net.Conn.
func (s StdioConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

// RemoteAddr implements net.Conn.
func (s StdioConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

// SetDeadline implements net.Conn.
func (s StdioConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements net.Conn.
func (s StdioConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implements net.Conn.
func (s StdioConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// Write implements net.Conn.
func (s StdioConn) Write(b []byte) (n int, err error) {
	return s.w.Write(b)
}

// Read implements net.Conn.
func (s StdioConn) Read(b []byte) (n int, err error) {
	return s.r.Read(b)
}
