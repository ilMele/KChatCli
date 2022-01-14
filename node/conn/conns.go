package conn

import (
	"fmt"
	"net"
)

type Client struct {
	Addr     string
	Name     string
	Password string
	Target   string
}

var client_list map[string]*Client

func AddRuntime(ser *net.UDPConn, ch chan *Client, del chan string) {
	client_list = make(map[string]*Client)

	for {
		select {
		case c := <-ch:
			addToList(c)
			fmt.Println("add client:")
			checkMatch(ser)
			showList()

		case a := <-del:
			removeClient(a)
		}
	}
}

func Get(name string) string {
	return ""
}

func SetTarget(name string, target string) {
	client_list[name].Target = target
}

func NameExists(name string) bool {
	if client_list[name] != nil {
		return true
	} else {
		return false
	}
}

func checkMatch(ser *net.UDPConn) {
	for i, v := range client_list {
		if client_list[v.Target] == nil {
			continue
		}
		if i == client_list[v.Target].Target {
			swapAddress(ser, v, client_list[v.Target])
		}
	}
}

func swapAddress(ser *net.UDPConn, a, b *Client) {
	Bdata := []byte(b.Addr + "\n")
	Aaddr, aerr := net.ResolveUDPAddr("udp", a.Addr)
	if aerr != nil {
		fmt.Println("unresolved: " + a.Addr)
	}

	Adata := []byte(a.Addr + "\n")
	Baddr, berr := net.ResolveUDPAddr("udp", b.Addr)
	if berr != nil {
		fmt.Println("unresolved: " + b.Addr)
	}
	if aerr == nil && berr == nil {
		ser.WriteToUDP(Bdata, Aaddr)
		ser.WriteToUDP(Adata, Baddr)
	}

	removeClient(a.Name)
	removeClient(b.Name)
}

func addToList(c *Client) {
	client_list[c.Name] = c
}

func removeClient(a string) {
	delete(client_list, a)
	fmt.Println("after delete:")
	showList()
}

func showList() {
	for i, v := range client_list {
		fmt.Printf("  name: %s, addr: %s\n", i, v.Addr)
	}
}
