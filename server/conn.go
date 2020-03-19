package server

import (
	"net"
	"time"
)

type conn struct {
	net.Conn
	IdleTimeout time.Duration
}

func (c *conn) Write(p []byte) (n int, err error) {
	c.updateDeadline()
	n, err = c.Conn.Write(p)
	return
}

func (c *conn) Read(b []byte) (n int, err error) {
	c.updateDeadline()
	n, err = c.Conn.Read(b)
	return
}

func (c *conn) Close() (err error) {
	err = c.Conn.Close()
	return
}

func (c *conn) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(idleDeadline)
}
