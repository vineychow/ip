package main

type IpService interface {
	Query(ip string) (string, error)
	Close()
}

// ------------------------------------------
// ipService

type serviceTask struct {
	ip     string
	result chan *serviceResult
}

type serviceResult struct {
	data string
	err  error
}

type ipService struct {
	route      QueryRoute
	closeQueue chan bool
	queue      chan *serviceTask
}

func NewIpService(routePath string) (IpService, error) {
	route, err := NewQueryRoute(routePath)
	if err != nil {
		return nil, err
	}
	is := &ipService{
		route:      route,
		closeQueue: make(chan bool),
		queue:      make(chan *serviceTask, 1024*10),
	}
	go is.runQuery()

	return is, nil
}

func (is *ipService) Close() {
	is.closeQueue <- true
	CloseDrivers()
}

func (is *ipService) Query(ip string) (string, error) {
	task := &serviceTask{ip: ip, result: make(chan *serviceResult)}
	is.queue <- task
	r := <-task.result
	return r.data, r.err
}

func (is *ipService) runQuery() {
	for {
		select {
		case <-is.closeQueue:
			return
		case t := <-is.queue:
			go func(ip string) {
				data, err := is.queryIp(ip)
				t.result <- &serviceResult{data, err}
			}(t.ip)
		}
	}
}

func (is *ipService) queryIp(ip string) (string, error) {
	drv := is.route.RouteBy(ip)
	return drv.Query(ip)
}
