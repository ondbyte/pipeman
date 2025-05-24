package golang

import (
	"net"
	"os"
)

type StdioListener struct {
	r, w *os.File
}

// Accept implements net.Listener.
func (s *StdioListener) Accept() (net.Conn, error) {
	return StdioConn{s.r, s.w}, nil
}

// Addr implements net.Listener.
func (s *StdioListener) Addr() net.Addr {
	return &net.TCPAddr{}
}

// Close implements net.Listener.
func (s *StdioListener) Close() error {
	return nil
}
