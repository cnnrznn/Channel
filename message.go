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

type Msg struct {
    From int
    Type MsgType
    Round int
    Data string
}

func MsgFromBytes(rawBytes []byte) Msg {
    var buffer bytes.Buffer
    dec := gob.NewDecoder(&buffer)
    var msg Msg

    buffer.Write(rawBytes)
    err := dec.Decode(&msg)
    if err != nil {
        log.Fatal("decode error:", err)
    }

    return msg
}

func (m Msg) MsgToBytes() []byte {
    var buffer bytes.Buffer
    enc := gob.NewEncoder(&buffer)

    err := enc.Encode(m)
    if err != nil {
        log.Fatal("encode error:", err)
    }

    return buffer.Bytes()
}
