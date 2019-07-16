package channel

import (
    "fmt"
    "net"
    "time"
)

type Channel struct {
    Id int
    Peers []string
}

func (c Channel) String() string {
    return fmt.Sprintf("{%v, %v}", c.Id, c.Peers)
}

func (c Channel) PingAll() {
    msg := Msg{c.Id,
               INITIAL,
               45,
               "ping"}
    for index, _ := range c.Peers {
        go c.Send(msg, index)
    }
}

func (c Channel) Send(msg Msg, index int) {
    conn, _ := net.Dial("udp", c.Peers[index])
    defer conn.Close()

    data := msg.MsgToBytes()
    buff := make([]byte, 128)

    for {
        conn.Write(data)

        time.Sleep(2 * time.Second)

        n, err := conn.Read(buff)
        if err != nil {
            continue
        }
        if string(buff[:n]) == "ok" {
            return
        }
    }
}

func (c Channel) Broadcast(msg Msg) {
    for index, _ := range c.Peers {
        go c.Send(msg, index)
    }
}

func (c Channel) Serve(dataChan chan Msg, addrChan chan string) {
    pc, _ := net.ListenPacket("udp", c.Peers[c.Id])
    defer pc.Close()

    buffer := make([]byte, 2048)

    for {
        n, addr, err := pc.ReadFrom(buffer)
        if err != nil {
            continue
        }

        // TODO try to convert the bytes into a Msg
        // if I can, ack the sender and return the message
        msg := MsgFromBytes(buffer[:n])

        pc.WriteTo([]byte("ok"), addr)

        dataChan <- msg
        addrChan <- addr.String()
    }
}
