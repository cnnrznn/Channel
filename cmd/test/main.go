package main

import (
  "channel"
  "flag"
)

func main() {
    id := *flag.Int("id", -1, "Index of this process in the config")
    conf := *flag.String("conf", "config.json", "Config file for the network")
    flag.Parse()

    // Load the JSON config here

    // Instantiate the channel
    ch := channel.Channel{Id: id}
}
