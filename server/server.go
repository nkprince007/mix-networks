package server

import (
	"bufio"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	sync.WaitGroup

	Addr            string
	IdleConnTimeout time.Duration
	BufferSize      int64

	listener         net.Listener
	commenceShutdown bool
	mu               sync.Mutex
	conns            []*conn
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = "8000"
	}

	log.Printf("Starting TCP server on %v\n", addr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	srv.listener = listener
	for {
		srv.mu.Lock()
		if srv.commenceShutdown {
			srv.Wait()
			return nil
		}
		srv.mu.Unlock()

		listener.SetDeadline(time.Now().Add(1e9))
		newConn, err := listener.AcceptTCP()

		if err != nil {
			netOpError, ok := err.(*net.OpError)
			if ok && (netOpError.Err.Error() == "use of closed network connection" || netOpError.Timeout()) {
				continue
			}

			log.Printf("error accepting connection: %v\n", err)
			continue
		}

		log.Printf("accepted connection from: %v\n", newConn.RemoteAddr())

		c := &conn{
			Conn:        newConn,
			IdleTimeout: srv.IdleConnTimeout,
			BufferSize:  srv.BufferSize,
		}
		c.SetDeadline(time.Now().Add(srv.IdleConnTimeout))

		srv.Add(1)
		srv.trackConnection(c)
		go func() {
			srv.handle(c)
			srv.Done()
		}()
	}
}

func (srv *Server) ShutDown() error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.commenceShutdown = true
	log.Println("shutting down...")
	return srv.listener.Close()
}

func (srv *Server) trackConnection(c *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	if srv.conns == nil {
		srv.conns = make([]*conn, 0)
	}
	srv.conns = append(srv.conns, c)
}

func (srv *Server) untrackConnection(c *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	for i, connection := range srv.conns {
		if connection == c {
			srv.conns[i] = srv.conns[len(srv.conns)-1]
			srv.conns[len(srv.conns)-1] = nil
			srv.conns = srv.conns[:len(srv.conns)-1]
		}
	}
}

func (srv *Server) handle(c *conn) error {
	defer func() {
		log.Printf("Closing connection from: %v\n", c.RemoteAddr())
		c.Close()
		srv.untrackConnection(c)
	}()

	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)
	scanner := bufio.NewScanner(r)

	sc := make(chan bool, 1)
	deadline := time.After(c.IdleTimeout)
	for {
		go func(s chan bool) {
			s <- scanner.Scan()
		}(sc)

		select {
		case <-deadline:
			return nil
		case scanned := <-sc:
			if !scanned {
				if err := scanner.Err(); err != nil {
					log.Printf("%v(%v)", err, c.RemoteAddr())
					return err
				}
				return nil
			}
			w.WriteString(scanner.Text() + "\n")
			w.Flush()
			deadline = time.After(c.IdleTimeout)
		}
	}
}
