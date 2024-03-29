package codec

import (
	"net/http"
	"strings"
	"time"
	"fmt"
	"errors"
	"sync"
	"sort"
)


const (
	pattern = "/test"
	headerParam = "X-Rpc-Serveraddr"

	expireTime = time.Second * 20
)

type serverItem struct {
	updateTime	time.Time
	server 		string
}

type serverMeta struct {
	expireTime 	time.Duration
	serverList 	map[string]*serverItem
	rwlock 		sync.RWMutex
}

var Meta = NewserverMeta()

func NewserverMeta() *serverMeta {
	return &serverMeta{
		expireTime:	expireTime,
		serverList:	map[string]*serverItem{},
	}
}

func (s *serverMeta) getServers() []string {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	tmpSlice := []string{}
	for k, v := range s.serverList {
		if v.updateTime.Add(s.expireTime).After(time.Now()) {
			tmpSlice = append(tmpSlice, v.server)
		} else {
			delete(s.serverList, k)
		}
	}

	sort.Strings(tmpSlice)
	return tmpSlice
}

func (s *serverMeta) updateServer(server string) error {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	if v, ok := s.serverList[server]; ok {
		v.updateTime = time.Now()
		return nil
	}

	s.serverList[server] = &serverItem{
		updateTime:time.Now(),
		server:server,
	}

	return nil
}

func (s *serverMeta) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header()[headerParam] =  []string{strings.Join(s.getServers(), ",")}
		for _, v := range s.serverList {
			fmt.Println(v.updateTime)
		}
	case "PUT":
		s.updateServer(r.Header[headerParam][0])
	default:
		http.NotFound(w, r)
	}
}

func test() {
	fmt.Println(strings.Join([]string{"a"}, ","))
}

func EtcdServer() {
	http.Handle(pattern, Meta)
	http.ListenAndServe(":9000", nil)
}

//心跳监测
func Heartbeat(addres string) {
	request, _ := http.NewRequest("PUT", "http://localhost:9000"+pattern, nil)
	request.Header[headerParam] = []string{addres}
	//request.Close = true

	doRequest(request)
	ch := time.Tick(expireTime - 1)
	for range ch {
		if err := doRequest(request); err != nil {
			return 
		}
	}
}

func doRequest(r *http.Request) error {
	if response, err := http.DefaultClient.Do(r); err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New("request error")
	} else {
		response.Body.Close()
	}

	return nil
}