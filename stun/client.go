package stun

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"
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

	_, _, _ = stepPing(conn)
}

// stepPing: ping 检测
// Return:
// ping -- 能够进行udp通信
// public -- 具有公网IP 或 对称型udp
// nat -- nat网关, 需进一步判断nat类型
func stepPing(conn *net.UDPConn) (ping, public, nat bool) {
	reqMsg := &StunMsg{
		ReqType: STUN_REQ_TYPE_PING,
	}
	reqBytes, err := json.Marshal(reqMsg)
	if err != nil {
		log.Fatal(err)
	}

	sendFailCount, timeoutCount, sameCount := 0, 0, 0
	received := make([]byte, 1024)
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Println("Error setting read deadline:", err)
	}
	for i := 0; i < 3; i++ {
		_, err = conn.Write(reqBytes)
		if err != nil {
			sendFailCount++
			continue
		}

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
