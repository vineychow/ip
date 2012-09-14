package main

const (
	taobaoName = "taobao"
	ipHost     = "http://ip.taobao.com"
	ipAddr     = ipHost + "/service/getIpInfo.php?"
)

func init() {
	RegisterDriver(newTaobaoService())
}

type taobaoTask struct {
	ip     string
	result chan *taobaoResult
}

type taobaoResult struct {
	data string
	err  error
}

type taobaoService struct {
	closeQueue chan bool
	queue      chan *taobaoTask
	client     IpClient
}

func newTaobaoService() IpDriver {
	ts := &taobaoService{
		closeQueue: make(chan bool),
		queue:      make(chan *taobaoTask, 1024),
		client:     NewIpClient(ipAddr),
	}
	go ts.runQuery()
	return ts
}

func (ts *taobaoService) Query(ip string) (string, error) {
	task := &taobaoTask{ip: ip, result: make(chan *taobaoResult)}
	ts.queue <- task
	r := <-task.result
	return r.data, r.err
}

func (ts *taobaoService) runQuery() {
	for {
		select {
		case <-ts.closeQueue:
			return
		case t := <-ts.queue:
			data, err := ts.client.QueryIp(t.ip)
			t.result <- &taobaoResult{data, err}
		}
	}
}

func (ts *taobaoService) Name() string {
	return taobaoName
}

func (ts *taobaoService) Close() error {
	ts.closeQueue <- true
	return nil
}
