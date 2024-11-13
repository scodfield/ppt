package stun

import (
	"encoding/json"
	"log"
	"net"
)

const (
	STUN_REQ_TYPE_PING     = iota // PING测试
	STUN_REQ_TYPE_FIREWALL        // 检测客户端防火墙
)

type StunMsg struct {
	ReqType int
	Addr    string `json:"addr,omitempty"`
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
