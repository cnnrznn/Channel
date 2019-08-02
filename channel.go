package channel

import (
    "fmt"
    "log"
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
    c.Broadcast(msg)
}

func (c Channel) Send(msg Msg, index int) {
    conn, err := net.Dial("udp", c.Peers[index])
    if err != nil {
        log.Fatal("Could not dial:", err)
    }
    defer conn.Close()

    data := msg.MsgToBytes()
    buff := make([]byte, 128)

    for {
        fmt.Println("Sending", msg, "to", c.Peers[index])
        n, err := conn.Write(data)
        if err != nil {
            log.Fatal(err)
        }

        conn.SetReadDeadline(time.Now().Add(2 * time.Second))
        n, err = conn.Read(buff)

        if err != nil {
            continue
        }
        if string(buff[:n]) == "ok" {
            fmt.Println("Received ack from", c.Peers[index])
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
        fmt.Println("Received", n, "bytes from", addr)

        // try to convert the bytes into a Msg
        // if I can, ack the sender and return the message
        msg := MsgFromBytes(buffer[:n])

        n, err = pc.WriteTo([]byte("ok"), addr)
        if err != nil {
            log.Fatal("Failed to send 'ok' ack")
        }
        fmt.Println("Sent ack to", addr, ":", msg)

        dataChan <- msg
        addrChan <- addr.String()
    }
}
