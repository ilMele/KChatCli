package conn

import (
	"net"
	"strings"
)

var node *net.UDPAddr
var addr net.UDPAddr
var dest *net.UDPAddr
var ser *net.UDPConn
var err error

func Init() *net.UDPConn {
	node, err = net.ResolveUDPAddr("udp", "151.81.61.158:50002")
	if err != nil {
		panic(err)
	}

	addr = net.UDPAddr{
		IP: net.ParseIP("0.0.0.0"),
	}

	ser, err = net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}

	return ser
}

func LoginRequest(you, user string) string {
	ser.WriteToUDP([]byte(you+" "+user+"\n"), node)
	data := make([]byte, 100)
	_, _, err = ser.ReadFromUDP(data)
	if err != nil {
		panic(err)
	}
	if data[0] == byte(1) {
		return ""
	}
	return string(data[1:])
}

func WaitAddr() {
	data := make([]byte, 20)

	_, _, err = ser.ReadFromUDP(data)
	if err != nil {
		panic(err)
	}
	addrstr := strings.SplitN(string(data), "\n", 2)[0]

	userAddr, uaerr := net.ResolveUDPAddr("udp", addrstr)
	if uaerr != nil {
		panic(uaerr)
	}
	dest = userAddr
}

func SendMsg(msg string) {
	_, err = ser.WriteTo([]byte(msg+"\n"), dest)
	if err != nil {
		panic(err)
	}
}

func ReadResponse(user string, msg chan string) {
	data := make([]byte, 512)
	for {
		ser.Read(data)
		msg <- user + strings.SplitN(string(data), "\n", 2)[0]
		data = make([]byte, 512)
	}
}
