package pkg

import (
	"fmt"
	"math/rand"
	"time"
)

var DefaultClient = NewClient()

type Client struct {
	ID int64
}

func NewClient() *Client {
	rand.Seed(time.Now().Unix())
	id := rand.Int63n(100)
	return &Client{
		ID: id,
	}
}

func (c *Client) PrintID() {
	fmt.Printf("Client ID: %v\n", c.ID)
}
