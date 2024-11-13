package stun

import (
	"encoding/json"
	"log"
	"net"
)

const (
	STUN_REQ_TYPE_PING         = iota // PING测试
	STUN_REQ_TYPE_PUBLIC              // 检测公网IP
	STUN_REQ_TYPE_FULLCONE_NAT        // 检测完全锥形NAT
)

const (
	NAT_TYPE_NONE              = iota
	NAT_TYPE_PUBLIC            // 公网IP 且 没有防火墙
	NAT_TYPE_SYMMETRY_FIREWALL // 对称型UDP防火墙
	NAT_TYPE_FULLCONE_NAT      // 完全锥形NAT
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

	newUdpServer, err := net.ListenPacket("udp", ":3479")
	if err != nil {
		log.Fatal(err)
	}
	defer newUdpServer.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			log.Printf("Error reading from STUN server: %v", err)
			continue
		}
		go response(udpServer, newUdpServer, addr, buf[:n])
	}
}

func response(udpServer, newUdpServer net.PacketConn, addr net.Addr, msg []byte) {
	reqMsg := &StunMsg{}
	err := json.Unmarshal(msg, reqMsg)
	if err != nil {
		log.Println(err)
		return
	}
	switch reqMsg.ReqType {
	case STUN_REQ_TYPE_PING:
		echoResponse(udpServer, addr)
	case STUN_REQ_TYPE_PUBLIC:
		changeIPAndPORTResponse(newUdpServer, addr)
	case STUN_REQ_TYPE_FULLCONE_NAT:
		changeIPAndPORTResponse(newUdpServer, addr)
	}
}

func echoResponse(udpServer net.PacketConn, addr net.Addr) {
	resMsg := StunMsg{
		Addr: addr.String(),
	}
	res, _ := json.Marshal(resMsg)
	udpServer.WriteTo(res, addr)
}

func changeIPAndPORTResponse(newUdpServer net.PacketConn, addr net.Addr) {
	resMsg := StunMsg{
		Addr: addr.String(),
	}
	res, _ := json.Marshal(resMsg)
	_, err := newUdpServer.WriteTo(res, addr)
	if err != nil {
		log.Printf("Error writing from NewUdpServer, Error: %v\n", err)
	}
}
