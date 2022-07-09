package codec

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type selectMode int

const (
	RandomSelect selectMode = iota
	RoundRobinSelect
)

type Discovery interface {
	Update([]string) error
	Get(selectMode) (string, error)
	GetAll() []string
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
	return m.servers
}
