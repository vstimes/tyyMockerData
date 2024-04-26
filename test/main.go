package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

var request = map[string]interface{}{
	"token": "fakeToken2",
}

type TcpData struct {
	MessageLength    string
	RegionCode       string
	ServiceType      string
	AgencyCode       string
	BusinessUnitCode string
	InitiatorFlag    string
	TransactionCode  string
	Message          map[string]interface{}
}

func padRight(str string, length int) string {
	l := length - len(str)
	ss := make([]string, l)
	for i := 0; i < l; i++ {
		ss = append(ss, " ")
	}

	return str + strings.Join(ss, "")
}
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:29401")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}

	msg, _ := json.Marshal(request)

	header := // MessageLength: 8 bytes
		"RCode" + // RegionCode: 4 bytes
			"TYPE01" + // ServiceType: 6 bytes
			padRight("ACode01", 15) + // AgencyCode: 15 bytes
			padRight("BCode01", 15) + // BusinessUnitCode: 15 bytes
			padRight("IFlag01", 8) + // InitiatorFlag: 8 bytes
			padRight("TYY4501", 15) // TransactionCode: 15 bytes

	all := header + string(msg)
	all = padRight(fmt.Sprintf("%d", 8+len([]byte(all))), 8) + all

	_, err = conn.Write([]byte(all))
	if err != nil {
		fmt.Println("write error:", err)
		return
	}

	buf := make([]byte, 8)
	_, err = io.ReadFull(conn, buf)
	if err == io.EOF {
		fmt.Println("Client disconnected")
		return
	}
	if err != nil {
		fmt.Println("Unexpected read error:", err)
		return
	}
	messageLen, _ := strconv.Atoi(string(buf)) // Convert string to int
	messageBuf := make([]byte, messageLen)
	_, err = io.ReadFull(conn, messageBuf)
	if err != nil {
		fmt.Println("Unexpected read error:", err)
		return
	}
	println(string(messageBuf))

	conn.Close()
}
