package stun

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"
)

func StunClient() int {
	udpServer, err := net.ResolveUDPAddr("udp", ":3478")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Println("Error setting read deadline:", err)
	}

	received := make([]byte, 0, 1024)

	_, isPublic, isNat := stepPing(conn, received)
	if isPublic {
		if stepPublic(conn, received) {
			return NAT_TYPE_PUBLIC
		}
		return NAT_TYPE_SYMMETRY_FIREWALL
	}

	if isNat {
		stepNat(conn, received)
	}
	return NAT_TYPE_NONE
}

// stepNat: NAT类型检测
// @return true 完全锥形NAT
// @return false 需进一步检查具体NAT类型
func stepNat(conn *net.UDPConn, received []byte) bool {
	reqBytes, err := json.Marshal(&StunMsg{
		ReqType: STUN_REQ_TYPE_FULLCONE_NAT,
	})
	if err != nil {
		log.Fatal(err)
	}
	if _, timeoutCount, _ := commonSend(conn, reqBytes, received); timeoutCount >= 3 {
		return false
	}
	return true
}

// stepPublic: public 检测
// @return true 公网IP 且 无防火墙
// @return false 对称型udp防火墙
func stepPublic(conn *net.UDPConn, received []byte) bool {
	reqBytes, err := json.Marshal(&StunMsg{
		ReqType: STUN_REQ_TYPE_PUBLIC,
	})
	if err != nil {
		log.Fatal(err)
	}

	if _, timeoutCount, _ := commonSend(conn, reqBytes, received); timeoutCount >= 3 {
		return false
	}
	return true
}

// stepPing: ping 检测
// @return ping=true 能够进行udp通信
// @return public=true 具有公网IP 或 对称型udp
// @return nat=true nat网关, 需进一步判断nat类型
func stepPing(conn *net.UDPConn, received []byte) (ping, public, nat bool) {
	reqBytes, err := json.Marshal(&StunMsg{
		ReqType: STUN_REQ_TYPE_PING,
	})
	if err != nil {
		log.Fatal(err)
	}

	sendFailCount, timeoutCount, sameCount := commonSend(conn, reqBytes, received)
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

func commonSend(conn *net.UDPConn, msg []byte, received []byte) (sendFailCount, timeoutCount, sameCount int) {
	for i := 0; i < 3; i++ {
		_, err := conn.Write(msg)
		if err != nil {
			sendFailCount++
			continue
		}

		received = received[0:0]
		n, err := conn.Read(received)
		if err != nil {
			var netError net.Error
			if errors.As(err, &netError) && netError.Timeout() {
				timeoutCount++
				continue
			}
		}
		recvMsg := &StunMsg{}
		json.Unmarshal(received[:n], recvMsg)
		if conn.LocalAddr().String() == recvMsg.Addr {
			sameCount++
		}
	}
	return
}
