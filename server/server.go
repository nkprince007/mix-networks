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

	listener         net.Listener
	commenceShutdown bool
	mu               sync.Mutex
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

		conn := &conn{
			Conn:        newConn,
			IdleTimeout: srv.IdleConnTimeout,
		}
		conn.SetDeadline(time.Now().Add(srv.IdleConnTimeout))

		srv.Add(1)
		go func() {
			srv.Done()
			handle(conn)
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

func handle(c *conn) error {
	defer func() {
		log.Printf("Closing connection from: %v\n", c.RemoteAddr())
		c.Close()
	}()

	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)
	scanner := bufio.NewScanner(r)
	deadline := time.After(c.IdleTimeout)
	for {
		select {
		case <-deadline:
			return nil
		default:
			scanned := scanner.Scan()
			if !scanned {
				if err := scanner.Err(); err != nil {
					log.Printf("%v(%v)", err, c.RemoteAddr())
					return err
				}
				break
			}
			w.WriteString(scanner.Text() + "\n")
			w.Flush()
			deadline = time.After(c.IdleTimeout)
		}
	}
}
