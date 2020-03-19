package server

import (
	"bufio"
	"log"
	"net"
	"sync"
)

type Server struct {
	Addr             string
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
			break
		}
		srv.mu.Unlock()

		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("error accepting connection: %v\n", err)
			continue
		}
		log.Printf("accepted connection from: %v\n", conn.RemoteAddr())
		handle(conn)
	}
	return nil
}

func (srv *Server) ShutDown() error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.commenceShutdown = true
	log.Println("shutting down...")
	return srv.listener.Close()
}

func handle(conn net.Conn) error {
	defer func() {
		log.Printf("Closing connection from: %v\n", conn.RemoteAddr())
		conn.Close()
	}()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(r)
	for {
		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				log.Printf("%v(%v)", err, conn.RemoteAddr())
				return err
			}
			break
		}
		w.WriteString(scanner.Text() + "\n")
		w.Flush()
	}
	return nil
}
