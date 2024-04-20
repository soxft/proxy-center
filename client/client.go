package client

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	ServerAddr string
	TimeOut    int
	EndTime    int64
}

func NewClient(serverAddr string, timeOut int, endTime int64) *Client {
	return &Client{
		ServerAddr: serverAddr,
		TimeOut:    timeOut,
		EndTime:    endTime,
	}
}

func (c *Client) Ping() error {
	client := resty.New().R()

	var r resp
	_, err := client.SetResult(&r).Get(fmt.Sprintf("http://%s/ping", c.ServerAddr))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetProxy() (ProxyData, error) {
	client := resty.New().R()

	var r resp
	_, err := client.SetResult(&r).Get(fmt.Sprintf("http://%s/getProxy/%d/%d", c.ServerAddr, c.TimeOut, c.EndTime))
	if err != nil {
		return ProxyData{}, err
	}

	if r.Success == false {
		return ProxyData{}, errors.New(r.Message)
	}

	return r.Data, nil
}
