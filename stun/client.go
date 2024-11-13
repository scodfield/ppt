package stun

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func StunClient() {
	udpServer, err := net.ResolveUDPAddr("udp", ":3478")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("STUN Client"))
	if err != nil {
		log.Fatal(err)
	}

	received := make([]byte, 1024)
	n, err := conn.Read(received)
	if err != nil {
		log.Fatal(err)
	}
	recvMsg := &StunMsg{}
	json.Unmarshal(received[:n], recvMsg)
	fmt.Printf("Local Addr:%v, Recv Addr:%s, Is same Addr:%t\n", conn.LocalAddr().String(), recvMsg.Addr, conn.LocalAddr().String() == recvMsg.Addr)
}
