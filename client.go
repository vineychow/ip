package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type IpClient interface {
	QueryIp(ip string) (string, error)
}

type ipClient struct {
	mux    sync.Mutex
	uv     url.Values
	client *http.Client
	uri    string
}

func NewIpClient(uri string) IpClient {
	return &ipClient{
		uv:     make(url.Values),
		client: &http.Client{},
		uri:    uri,
	}
}

func (ic *ipClient) QueryIp(ip string) (string, error) {
	ic.mux.Lock()
	defer ic.mux.Unlock()

	resp, err := ic.queryIp(ip)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (ic *ipClient) queryIp(ip string) (string, error) {
	m := ic.uv
	m.Add("ip", ip)

	// client
	client := ic.client
	url := ic.uri + m.Encode()
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
