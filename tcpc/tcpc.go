package tcpc

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

type TCPC[T any] struct {
	listenAddr string
	remoteAddr string

	SendChan     chan T
	RecvChan     chan T
	outboundConn net.Conn
	ln           net.Listener
}

func New[T any](listenAddr, remoteAddr string) (*TCPC[T], error) {
	var err error
	tcpc := &TCPC[T]{
		listenAddr: listenAddr,
		SendChan:   make(chan T),
		remoteAddr: remoteAddr,
		RecvChan:   make(chan T),
	}

	tcpc.ln, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	go tcpc.loop()
	go tcpc.acceptLoop()
	go tcpc.dialRemoteAndRead()

	return tcpc, nil
}

func (t *TCPC[T]) loop() {
	for {
		msg := <-t.SendChan
		log.Println("sending message over the wire: ", msg)
		if err := gob.NewEncoder(t.outboundConn).Encode(&msg); err != nil {
			log.Print("error: ", err)
			continue
		}
	}
}

func (t *TCPC[T]) acceptLoop() {

	for {
		conn, err := t.ln.Accept()
		defer func() {
			t.ln.Close()
		}()

		if err != nil {
			log.Printf("accept error: %v", err)
			return

		}
		log.Printf("sender connected: %v\n", conn.RemoteAddr())

		go t.handleConn(conn)
	}
}

func (t *TCPC[T]) handleConn(conn net.Conn) {
	for {
		var msg T
		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
			log.Println(err)
			continue
		}
		t.RecvChan <- msg
	}
}

func (t *TCPC[T]) dialRemoteAndRead() {
	conn, err := net.Dial("tcp", t.remoteAddr)
	if err != nil {
		log.Printf("dial error: %v ", err)
		time.Sleep(time.Second * 3)
		t.dialRemoteAndRead()
	}

	t.outboundConn = conn
}
