package codec

import (
	"errors"
	"math/rand"
	"sync"
	"time"
	"strings"
	"net/http"
)

type selectMode int

const (
	RandomSelect selectMode = iota
	RoundRobinSelect
)

type Discovery interface {
	Refresh() error
	Update([]string) error
	Get(selectMode) (string, error)
	GetAll() []string
}

type EtcdServersDiscovery struct {
	updateTime 	time.Time
	expire 		time.Duration
	*MultiServersDiscovery
}

func (e *EtcdServersDiscovery) Refresh() error {
	e.lock.Lock()
	defer e.lock.Unlock()
	//判断是否缓存过期
	if e.updateTime.Add(e.expire).After(time.Now()) {
		return nil
	}
	//测试 要通过http get请求来拿服务列表
	//e.servers = Meta.getServers()
	if response, err := http.Get("http://localhost:9000"+pattern); err != nil {
		return err
	} else {
		if response != nil {
			defer response.Body.Close()
		}

		e.servers = strings.Split(response.Header[headerParam][0], ",")
	}
	
	e.updateTime = time.Now()
	return nil
}

func (e *EtcdServersDiscovery) Update(servers []string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.servers = servers
	e.updateTime = time.Now()
	return nil
}

func (e *EtcdServersDiscovery) GetAll() []string {
	e.Refresh()
	return e.GetAll()
}

func (e *EtcdServersDiscovery) Get(mode selectMode) (string, error) {
	e.Refresh()
	return e.MultiServersDiscovery.Get(mode)
}

func NewEtcdServersDiscovery() *EtcdServersDiscovery {
	e := &EtcdServersDiscovery{
		expire:time.Minute,
		MultiServersDiscovery:NewMultiServersDiscovery(),
	}

	return e
}


type MultiServersDiscovery struct {
	servers []string
	index   int //轮询索引
	lock    sync.Mutex
	rand    *rand.Rand
}

func NewMultiServersDiscovery() *MultiServersDiscovery {
	m := &MultiServersDiscovery{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	//初始化随机
	m.index = m.rand.Int()

	return m
}

func (m *MultiServersDiscovery) Update(servers []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.servers = servers
	return nil
}

func (m *MultiServersDiscovery) Get(mode selectMode) (string, error) {
	n := len(m.servers)
	switch mode {
	case RoundRobinSelect:
		m.lock.Lock()
		defer m.lock.Unlock()
		m.index = (m.index + 1) % n
		return m.servers[m.index%n], nil
	case RandomSelect:
		return m.servers[m.rand.Intn(n)], nil
	default:
		return "", errors.New("mode error")
	}
}

func (m *MultiServersDiscovery) GetAll() []string {
	m.lock.Lock()
	defer m.lock.Unlock()
	tempSlice := make([]string, len(m.servers))
	copy(tempSlice, m.servers)
	return tempSlice
}

func (m *MultiServersDiscovery) Refresh() []string {
	return nil
}