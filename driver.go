package main

import (
	"errors"
	"io"
)

type IpDriver interface {
	io.Closer
	Name() string
	Query(ip string) (string, error)
}

var (
	drivers = make(map[string]IpDriver)
)

func RegisterDriver(drv IpDriver) {
	name := drv.Name()
	if _, dup := drivers[name]; dup {
		panic("RegisterDriver: repeated driver " + name)
	}
	drivers[name] = drv
}

func DriverByName(name string) (IpDriver, error) {
	if drv, ok := drivers[name]; ok {
		return drv, nil
	}
	return nil, errors.New("DriverByName: no such driver " + name)
}

func CloseDrivers() error {
	errStr := ""
	var hasErr bool
	for k, v := range drivers {
		if err := v.Close(); err != nil {
			hasErr = true
			errStr += (k + ", " + err.Error() + "; ")
		}
		delete(drivers, k)
	}
	if !hasErr {
		return nil
	}
	return errors.New(errStr)
}
