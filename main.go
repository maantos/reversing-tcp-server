package main

import (
	"fmt"
	"log"
	"net"
)

const address = "localhost:3000"

type server struct {
	ln         net.Listener
	listenAddr string
	quitch     chan struct{}
}

func newServer(Addr string) *server {
	return &server{
		listenAddr: Addr,
		quitch:     make(chan struct{}),
	}
}

func (s *server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	fmt.Printf("TCP server listening on %s\n", s.listenAddr)
	defer ln.Close()
	s.ln = ln
	go s.acceptLoop()
	<-s.quitch
	return nil
}

func (s *server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Printf("accept error: %s\n", err.Error())
			continue
		}
		s.handleRequest(conn)
	}
}

func (s *server) handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	msg := buf[:n]
	msg = reverseString(msg)
	// Send a response back to person contacting us.
	conn.Write(msg)
	// Close the connection when you're done with it.
	conn.Close()
}

func reverseString(a []byte) []byte {

	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a

}

func tcpClient(messages []string) ([]string, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return []string{}, err
	}

	reversed := []string{}

	for _, msg := range messages {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return []string{}, err
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			return []string{}, err
		}

		reply := make([]byte, 1024)
		n, err := conn.Read(reply)
		if err != nil {
			return []string{}, err
		}

		fmt.Printf("reply:%s\n", string(reply))
		reversed = append(reversed, string(reply[:n]))
		conn.Close()
	}
	return reversed, nil
}

func main() {

	s := newServer(address)
	go func() {
		err := s.Start()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()
	msg := []string{"Tomasz", "ma", "jednego", "psa"}
	reversedMsg, err := tcpClient(msg)
	if err != nil {
		fmt.Println(err)
	}

	for _, rm := range reversedMsg {
		fmt.Println(rm)
	}

	// signalChnl := make(chan os.Signal, 1)
	// signal.Notify(
	//
	//	signalChnl,
	//	syscall.SIGINT,
	//	syscall.SIGTERM,
	//
	// )
	//
	// //receiving from channel is blocking operation, so its blocked until signal received
	// sig := <-signalChnl
	// log.Println("Got signal: ", sig)
}
