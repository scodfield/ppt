package stun

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const (
	STUN_REQ_TYPE_PING            = iota // PING测试
	STUN_REQ_TYPE_PUBLIC                 // 检测公网IP
	STUN_REQ_TYPE_FULLCONE_NAT           // 检测完全锥形NAT
	STUN_REQ_TYPE_PORT_RESTRICTED        // 端口受限NAT
)

const (
	NAT_TYPE_NONE                = iota
	NAT_TYPE_PUBLIC              // 公网IP 且 没有防火墙
	NAT_TYPE_SYMMETRY_FIREWALL   // 对称型UDP防火墙
	NAT_TYPE_FULLCONE_NAT        // 完全锥形NAT
	NAT_TYPE_SYMMETRY_NAT        // 对称型NAT
	NAT_TYPE_PORT_RESTRICTED_NAT // 端口受限型NAT
	NAT_TYPE_RESTRICTED_NAT      // 受限型NAT
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
	fmt.Println("STUN server listening on :3478")

	newUdpServer, err := net.ListenPacket("udp", ":3479")
	if err != nil {
		log.Fatal(err)
	}
	defer newUdpServer.Close()
	fmt.Printf("NewudpServer Local Addr: %s.\n", newUdpServer.LocalAddr().String())

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
	case STUN_REQ_TYPE_PORT_RESTRICTED:
		changePortResponse(udpServer, addr)
	}
}

func echoResponse(udpServer net.PacketConn, addr net.Addr) {
	resMsg := StunMsg{
		Addr: addr.String(),
	}
	res, _ := json.Marshal(resMsg)
	n, err := udpServer.WriteTo(res, addr)
	fmt.Printf("EchoPing, %d bytes written to:%s, and error: %v.\n", n, addr.String(), err)
}

func changeIPAndPORTResponse(newUdpServer net.PacketConn, addr net.Addr) {
	resMsg := StunMsg{
		Addr: addr.String(),
	}
	res, _ := json.Marshal(resMsg)
	n, err := newUdpServer.WriteTo(res, addr)
	if err != nil {
		log.Printf("NewUdpServer Error writing back, send bytes:%d,Error: %v\n", n, err)
	} else {
		log.Printf("changeIPAndPortResponse, NewUdpServer Write bytes: %d\n", n)
	}
}

func changePortResponse(udpServer net.PacketConn, addr net.Addr) {
	resMsg := StunMsg{
		Addr: addr.String(),
	}
	res, _ := json.Marshal(resMsg)
	udpServer.WriteTo(res, addr)

	// another port
	lAddr, err := net.ResolveUDPAddr("udp", ":3480")
	if err != nil {
		log.Println(err)
	}
	raddr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		log.Println(err)
	}
	newConn, err := net.DialUDP("udp", lAddr, raddr)
	if err != nil {
		log.Println(err)
	}
	defer newConn.Close()
	_, err = newConn.Write(res)
	if err != nil {
		log.Println(err)
	}
}
