package main

type QueryRoute interface {
	RouteBy(ip string) IpDriver
}

type queryRoute struct {
	path string
}

func NewQueryRoute(routePath string) (QueryRoute, error) {
	qr := &queryRoute{
		path: routePath,
	}
	return qr, nil
}

func (qr *queryRoute) RouteBy(ip string) IpDriver {
	// taobao
	drv, err := DriverByName("taobao")
	if err == nil {
		return drv
	}

	// default
	drv, err = DriverByName("default")
	if drv == nil {
		panic("drv!=nil")
	}

	return drv

}
