package client

import (
	"fmt"
	"log"
	"sync"

	"github.com/reiver/go-telnet"
)

type Client struct {
	t        *telnet.Conn
	mutex    sync.Mutex
	login    string
	password string
}

func (c *Client) readUntil(read bool, delims ...string) ([]byte, int, error) {
	if len(delims) == 0 {
		return nil, 0, nil
	}
	p := make([]string, len(delims))
	for i, s := range delims {
		if len(s) == 0 {
			return nil, 0, nil
		}
		p[i] = s
	}
	var line []byte
	b := make([]byte, 1)
	for {
		_, err := c.t.Read(b)
		if err != nil {
			return nil, 0, err
		}
		if read {
			line = append(line, b...)
		}
		for i, s := range p {
			if s[0] == b[0] {
				if len(s) == 1 {
					return line, i, nil
				}
				p[i] = s[1:]
			} else {
				p[i] = delims[i]
			}
		}
	}
}

func (c *Client) ExecuteCommand(cmd string) []byte {
	c.mutex.Lock()
	c.t.Write([]byte(cmd + "\n"))
	result, _ := c.ReadUntil(">")
	c.mutex.Unlock()
	return result
}

func (c *Client) ReadUntil(delims ...string) ([]byte, error) {
	d, _, err := c.readUntil(true, delims...)
	return d, err
}

func (c *Client) SkipUntil(delims ...string) error {
	_, _, err := c.readUntil(false, delims...)
	return err
}

func expect(c *Client, d ...string) {
	c.SkipUntil(d...)
}

func (c *Client) Authorize() {
	c.mutex.Lock()

	expect(c, "Login:")
	fmt.Println("l")
	c.t.Write([]byte(c.login + "\n"))
	expect(c, "Password:")
	fmt.Println("p")
	c.t.Write([]byte(c.password + "\n"))
	expect(c, ">")

	c.mutex.Unlock()
}

func New(dst string, login string, password string) (*Client, error) {
	t, err := telnet.DialTo(dst)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		t:        t,
		mutex:    sync.Mutex{},
		login:    login,
		password: password,
	}, nil
}
