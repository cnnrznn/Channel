package channel

import (
    "fmt"
    "net"
)

type Channel struct {
    Id int
    Peers []string
}

func (c Channel) String() string {
    return fmt.Sprintf("{%v, %v}", c.Id, c.Peers)
}

//func (c Channel) PingAll() {
//    for index, _ := range c.Peers {
//        go c.Send("ping", index)
//    }
//}

func (c Channel) Send(msg Message, index int) {
    conn, _ := net.Dial("udp", c.Peers[index])
    defer conn.Close()

    conn.Write(msg.ToBytes())
}

func (c Channel) Broadcast(msg Message) {
    for index, _ := range c.Peers {
        go c.Send(msg, index)
    }
}

func (c Channel) Serve(dataChan chan Message, addrChan chan string) {
    pc, _ := net.ListenPacket("udp", c.Peers[c.Id])
    defer pc.Close()

    buffer := make([]byte, 2048)

    for {
        n, addr, err := pc.ReadFrom(buffer)
        if err != nil {
            continue
        }

        // TODO try to convert the bytes into a Message
        // if I can, ack the sender and return the message
        msg := MessageFromBytes(buffer[:n])

        dataChan <- msg
        addrChan <- addr.String()
    }
}
