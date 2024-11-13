package stun

import (
	"encoding/json"
	"log"
	"net"
)

type StunMsg struct {
	Addr string
}

func StunServer() {
	udpServer, err := net.ListenPacket("udp", ":3478")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()
	for {
		buf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			log.Printf("Error reading from STUN server: %v", err)
			continue
		}
		go response(udpServer, addr, buf[:n])
	}
}

func response(udpServer net.PacketConn, addr net.Addr, msg []byte) {
	resMsg := StunMsg{
		Addr: addr.String(),
	}
	res, _ := json.Marshal(resMsg)
	udpServer.WriteTo(res, addr)
}
