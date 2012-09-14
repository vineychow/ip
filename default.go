package main

import (
	"errors"
)

func init() {
	RegisterDriver(newDefaultService())
}

type defaultService struct {
}

func newDefaultService() IpDriver {
	return &defaultService{}
}

func (ds *defaultService) Name() string {
	return "default"
}

func (ds *defaultService) Close() error {
	return nil
}

func (ds *defaultService) Query(ip string) (string, error) {
	return "", errors.New("NotImpl: defaultService")
}
