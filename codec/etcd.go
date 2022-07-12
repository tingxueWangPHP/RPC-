package codec

import (
	"net/http"
	"strings"
	"time"
	"fmt"
	"errors"
	"sync"
)


const (
	pattern = "/test"
	headerParam = "X-Rpc-Serveraddr"
)

type serverItem struct {
	updateTime	time.Time
	server 		string
}

type serverMeta struct {
	expireTime 	time.Duration
	serverList 	[]*serverItem
	rwlock 		sync.RWMutex
}

func NewserverMeta() *serverMeta {
	return &serverMeta{
		expireTime:	time.Second * 20,
		serverList:	[]*serverItem{},
	}
}

func (s *serverMeta) GetServers() []string {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	tmpSlice := []string{}
	for _, v := range s.serverList {
		if v.updateTime.Add(s.expireTime).After(time.Now()) {
			tmpSlice = append(tmpSlice, v.server)
		}
	}

	return tmpSlice
}

func (s *serverMeta) updateServer(server string) error {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	for _, v := range s.serverList {
		if v.server == server {
			v.updateTime = time.Now()
			return nil
		}
	}

	s.serverList = append(s.serverList, &serverItem{
		updateTime:time.Now(),
		server:server,
	})

	return nil
}

func (s *serverMeta) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header()[headerParam] =  []string{strings.Join(s.GetServers(), ",")}
		fmt.Println(s.serverList[0].updateTime)
		fmt.Println(s.serverList[1].updateTime)
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
	http.ListenAndServe(":9000", NewserverMeta())
}

//心跳监测
func Heartbeat(addres string) {
	request, _ := http.NewRequest("PUT", "http://localhost:9000"+pattern, nil)
	request.Header[headerParam] = []string{addres}
	//request.Close = true

	doRequest(request)
	ch := time.Tick(time.Second * 20)
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