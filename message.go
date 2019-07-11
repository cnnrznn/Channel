package channel

import (
    "bytes"
    "encoding/gob"
    "log"
)

type MsgType int

const (
    INITIAL MsgType = iota
    ECHO
)

type Message struct {
    From int
    Type MsgType
    Round int
    Data string
}

func MessageFromBytes(rawBytes []byte) Message {
    var buffer bytes.Buffer
    dec := gob.NewDecoder(&buffer)
    var msg Message

    buffer.Write(rawBytes)
    err := dec.Decode(&msg)
    if err != nil {
        log.Fatal("decode error:", err)
    }

    return msg
}

func (m Message) ToBytes() []byte {
    var buffer bytes.Buffer
    enc := gob.NewEncoder(&buffer)

    err := enc.Encode(m)
    if err != nil {
        log.Fatal("encode error:", err)
    }

    return buffer.Bytes()
}
