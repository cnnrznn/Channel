package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "github.com/cnnrznn/channel"
  "io/ioutil"
  "time"
)

type Config struct {
    Peers []string `json:"peers"`
}

func checkError(err error) {
    if nil != err {
        fmt.Println(err)
    }
}

func main() {
    var err error

    id := flag.Int("id", -1, "Index of this process in the config")
    confFn := flag.String("conf", "config.json", "Config file for the network")
    flag.Parse()

    // Load the JSON config here
    config := Config{}

    confData, err := ioutil.ReadFile(*confFn)
    checkError(err)

    err = json.Unmarshal([]byte(confData), &config)
    checkError(err)

    // Instantiate the channel
    ch := channel.Channel{Id: *id,
                          Peers: config.Peers}
    fmt.Println(ch)

    msgChan := make(chan channel.Msg)
    addrChan := make(chan string)
    go ch.Serve(msgChan, addrChan)

    for i:=0; i<10; i++ {
        ch.PingAll()
    }

    for i:=0; i<10*len(ch.Peers); i++ {
        fmt.Println(<-msgChan, <-addrChan)
    }

    for {
        fmt.Println("Done receiving")
        time.Sleep(5 * time.Second)
    }
}
