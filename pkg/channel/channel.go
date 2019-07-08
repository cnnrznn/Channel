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

func (c Channel) PingAll() int {
    ch := make(chan int)

    for index, _ := range c.Peers {
        go c.Send("ping", index, ch)
    }

    for range c.Peers {
        <-ch
    }

    return 0
}

func (c Channel) Send(msg string, index int, ch chan int) {
    fmt.Println("Sending", msg, "to", c.Peers[index])

    conn, _ := net.Dial("udp", c.Peers[index])
    defer conn.Close()

    conn.Write([]byte(msg))

    ch <- 0
}

func (c Channel) Serve(ch chan string) {
    pc, _ := net.ListenPacket("udp", c.Peers[c.Id])
    defer pc.Close()

    buffer := make([]byte, 1024)

    for {
        n, _, err := pc.ReadFrom(buffer)
        if err != nil {
            continue
        }
        ch <- string(buffer[:n])
    }
}
