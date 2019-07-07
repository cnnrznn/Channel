package channel

import (
    "fmt"
)

type Channel struct {
    Id int
    Peers []string
}

func (c Channel) String() string {
    return fmt.Sprintf("{%v, %v}", c.Id, c.Peers)
}

func (c Channel) PingAll() int {
    return -1
}
