package stun

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func StunClient() int {
	udpServer, err := net.ResolveUDPAddr("udp", ":3478")
	if err != nil {
		log.Fatal(err)
	}
	lAddr, err := net.ResolveUDPAddr("udp", ":4008")
	conn, err := net.DialUDP("udp", lAddr, udpServer)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, isPublic, isNat := stepPing(conn)
	if isPublic {
		if stepPublic(conn) {
			fmt.Printf("step public is true...\n")
			return NAT_TYPE_PUBLIC
		}
		fmt.Printf("step public is false...\n")
		return NAT_TYPE_SYMMETRY_FIREWALL
	}

	if isNat {
		if stepNat(conn) {
			return NAT_TYPE_FULLCONE_NAT
		}
	}

	if !stepSymmetryNat(conn, lAddr) {
		return NAT_TYPE_SYMMETRY_NAT
	}

	if stepPortRestricted(conn) {
		return NAT_TYPE_PORT_RESTRICTED_NAT
	}

	return NAT_TYPE_RESTRICTED_NAT
}

// stepPortRestricted: 是否端口受限型NAT
// @return true 端口受限型NAT, 外服的ip/port不能变
// @return false 受限型NAT, 只限制外服的ip
func stepPortRestricted(conn *net.UDPConn) bool {
	reqBytes, err := json.Marshal(&StunMsg{
		ReqType: STUN_REQ_TYPE_PORT_RESTRICTED,
	})
	if err != nil {
		log.Fatal(err)
	}
	if _, timeoutCount, _ := commonSend(conn, reqBytes, 3); timeoutCount >= 3 {
		return true
	}
	return false
}

// stepSymmetryNat: 是否对称型NAT
// @return true 受限锥形NAT, 需进一步判断是否为端口受限型NAT
// @return false 对称型NAT,每次连接分配不同的映射端口
func stepSymmetryNat(conn *net.UDPConn, lAddr *net.UDPAddr) bool {
	reqBytes, err := json.Marshal(&StunMsg{
		ReqType: STUN_REQ_TYPE_PING,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(reqBytes)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Println("Error setting read deadline:", err)
	}
	recv1 := make([]byte, 1024)
	n, err := conn.Read(recv1)
	if err != nil {
		log.Fatal(err)
	}
	recvMsg := &StunMsg{}
	json.Unmarshal(recv1[:n], recvMsg)
	firstEchoAddr := recvMsg.Addr
	newUdpServer, err := net.ResolveUDPAddr("udp", ":3479")
	if err != nil {
		log.Fatal(err)
	}
	newConn, err := net.DialUDP("udp", lAddr, newUdpServer)
	if err != nil {
		log.Fatal(err)
	}
	defer newConn.Close()
	_, err = newConn.Write(reqBytes)
	if err != nil {
		log.Fatal(err)
	}
	err = newConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Println("Error setting read deadline:", err)
	}
	recv2 := make([]byte, 1024)
	n, err = newConn.Read(recv2)
	if err != nil {
		log.Fatal(err)
	}
	recvMsg = &StunMsg{}
	json.Unmarshal(recv2[:n], recvMsg)
	secondEchoAddr := recvMsg.Addr
	return firstEchoAddr == secondEchoAddr
}

// stepNat: NAT类型检测
// @return true 完全锥形NAT
// @return false 需进一步检查具体NAT类型
func stepNat(conn *net.UDPConn) bool {
	reqBytes, err := json.Marshal(&StunMsg{
		ReqType: STUN_REQ_TYPE_FULLCONE_NAT,
	})
	if err != nil {
		log.Fatal(err)
	}
	if _, timeoutCount, _ := commonSend(conn, reqBytes, 3); timeoutCount >= 3 {
		return false
	}
	return true
}

// stepPublic: public 检测
// @return true 公网IP 且 无防火墙
// @return false 对称型udp防火墙
func stepPublic(conn *net.UDPConn) bool {
	reqBytes, err := json.Marshal(&StunMsg{
		ReqType: STUN_REQ_TYPE_PUBLIC,
	})
	if err != nil {
		log.Fatal(err)
	}

	if _, timeoutCount, _ := commonSend(conn, reqBytes, 3); timeoutCount >= 3 {
		return false
	}
	return true
}

// stepPing: ping 检测
// @return ping=true 能够进行udp通信
// @return public=true 具有公网IP 或 对称型udp
// @return nat=true nat网关, 需进一步判断nat类型
func stepPing(conn *net.UDPConn) (ping, public, nat bool) {
	reqBytes, err := json.Marshal(&StunMsg{
		ReqType: STUN_REQ_TYPE_PING,
	})
	if err != nil {
		log.Fatal(err)
	}

	sendFailCount, timeoutCount, sameCount := commonSend(conn, reqBytes, 3)
	fmt.Printf("sendFailCount:%d, timeoutCount:%d, sameCount:%d\n", sendFailCount, timeoutCount, sameCount)
	if sendFailCount <= 0 {
		ping = true
	}
	if sendFailCount <= 0 && timeoutCount <= 0 {
		if sameCount >= 3 {
			public = true
		} else {
			nat = true
		}
	}

	return
}

func commonSend(conn *net.UDPConn, msg []byte, count int) (sendFailCount, timeoutCount, sameCount int) {
	for i := 0; i < count; i++ {
		_, err := conn.Write(msg)
		if err != nil {
			sendFailCount++
			continue
		}

		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Println("Error setting read deadline:", err)
		}
		recv := make([]byte, 1024)
		n, err := conn.Read(recv)
		if err != nil {
			var netError net.Error
			if errors.As(err, &netError) && netError.Timeout() {
				timeoutCount++
				continue
			}
		}
		recvMsg := &StunMsg{}
		json.Unmarshal(recv[:n], recvMsg)
		if conn.LocalAddr().String() == recvMsg.Addr {
			sameCount++
		}
	}
	return
}
