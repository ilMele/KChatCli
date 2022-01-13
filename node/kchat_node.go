package main

import(
    "fmt"
    "net"
    "strings"
    "kiwichat/node/conn"
)

const port = 50002

func main(){

    addr := net.UDPAddr{
        Port: port,
        IP: net.ParseIP("0.0.0.0"),
    }

    ser, err := net.ListenUDP("udp", &addr)

    if err != nil {
        panic(err)
    }
   
    defer ser.Close()
    
    data := make([]byte, 51) 

    ch := make(chan *conn.Client)
    del := make(chan string)

    go conn.AddRuntime(ser, ch, del)
    
    for {
        _, remote, err := ser.ReadFromUDP(data)

        if err != nil {
            panic(err)
        }
        go requestHandler(ser, remote, data, ch)
    }
        
}

func requestHandler(ser *net.UDPConn, addr *net.UDPAddr, res []byte, ch chan *conn.Client){
    datastr := shrinkNames(res)
    fmt.Println(datastr)
    names := strings.Split(datastr, " ")
    fmt.Println(len(names[0]), len(names[1]))

    if conn.NameExists(names[0]){
        errormsg := []byte{0}
        errormsg = append(errormsg, []byte("name already exists")...)
        ser.WriteToUDP(errormsg, addr) 
        return
    }

    ch <- &conn.Client{Addr: addr.String(), Name: names[0], Target: names[1]}
    ser.WriteToUDP([]byte{1}, addr) 
}

func shrinkNames(data []byte) string{
    var str string
    for _, v := range data{
        if v == '\n'{
            break
        }
        str += string(v)
    }
    return str
}
