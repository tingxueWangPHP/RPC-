package codec

import (
	"context"
	"errors"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
	//"fmt"
)

type DataFormat struct {
	Header string
	Body   []interface{}
}

type (
	Server struct{
		relation map[string]interface{}
		lock 	sync.Mutex
	}

	Client struct {
		conn 	net.Conn
		lock 	sync.Mutex
		address string
		isClosed bool
	}
)

type Xclient struct {
	clientCache map[string]*Client
	lock 	sync.Mutex
}

var once sync.Once
var x *Xclient


func NewServer() *Server {
	return &Server{
		relation:make(map[string]interface{}),
	}
}

func (s *Server) Register(structItem interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.relation[reflect.TypeOf(structItem).Name()] = structItem
}

func (s *Server) parseHeader(header string) (string, string) {
	ss := strings.Split(header, ".")
	return ss[0], ss[1]
}

func ListenAndServe(s *Server, protocol, address string) error {
	listen, err := net.Listen(protocol, address)
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			defer conn.Close()
			
			data, err := Decode(conn)
			for err == nil {
				structName, methodName := s.parseHeader(data.Header)
				fn := reflect.ValueOf(s.relation[structName]).MethodByName(methodName)

				//构建参数
				tempParam := []reflect.Value{}
				for i := 0; i < len(data.Body); i++ {
					tempParam = append(tempParam, reflect.ValueOf(data.Body[i]))
				}
				values := fn.Call(tempParam)

				d := DataFormat{
					Header: "reply",
				}
				for _, item := range values {
					d.Body = append(d.Body, item.Interface())
				}

				Encode(conn, d)
				data, err = Decode(conn)
			}
		}(conn)
	}

	return nil
}

func dialServer(protocol, address string) (*Client, error) {
	x.lock.Lock()
	defer x.lock.Unlock()
	//判断cache里是否有这个链接
	if v, ok := x.clientCache[address]; ok {
		if !v.isClosed {
			return v, nil
		}
	}
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return nil, err
	}

	x.clientCache[address] = &Client{
		conn: conn,
		address: address,
	}

	return x.clientCache[address], nil

}

func DialServer(discovery Discovery, mode selectMode) (*Client, error) {
	//优化
	once.Do(func(){
		x = &Xclient{
			clientCache:make(map[string]*Client),
		}
	})
	
	if addr, err := discovery.Get(mode); err != nil {
		return nil, err
	} else {
		return dialServer("tcp", addr)
	}
}

func ClientsClose() error {
	for k, v := range x.clientCache {
		v.Close()
		delete(x.clientCache, k)
	}

	return nil
}

func (c *Client) Close() error {
	c.isClosed = true
	return c.conn.Close()
}

func (c *Client) Call(cli *Client, method string, timeout time.Duration, ret interface{}, args ...interface{}) error {
	if reflect.TypeOf(ret).Kind() != reflect.Ptr {
		return errors.New("type error")
	}

	//防止协程因为没有消费者而阻塞 产生泄露
	resCh := make(chan *DataFormat, 1)
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	go func() {
		cli.lock.Lock()
		defer cli.lock.Unlock()
		Encode(cli.conn, DataFormat{
			Header: method,
			Body:   args,
		})
		res, _ := Decode(cli.conn)
		resCh <- res
	}()

	select {
	case res := <-resCh:
		if res == nil {
			return errors.New("network error")
		}
		for _, item := range res.Body {
			if err := recursionGet2(item, ret); err != nil {
				return err
			}
		}
		close(resCh)
	case <-ctx.Done():
		return errors.New("time out")
	}

	return nil
}

func recursionGet2(src, des interface{}) error {
	var (
		desElem  = reflect.ValueOf(des).Elem()
		srcValue = reflect.ValueOf(src)
		srcType  = reflect.TypeOf(src)

		result reflect.Value
	)

	//判断是否是指针
	if !desElem.CanSet() {
		return errors.New("cannot set")
	}

	switch srcValue.Kind() {
	case reflect.Slice:
		newSlice := reflect.MakeSlice(srcType, srcValue.Len(), srcValue.Cap())
		for i := 0; i < newSlice.Len(); i++ {
			p := reflect.New(srcValue.Index(i).Type())
			recursionGet2(srcValue.Index(i).Interface(), p.Interface())
			newSlice.Index(i).Set(p.Elem())
		}
		result = newSlice
	case reflect.Map:
		newMap := reflect.MakeMap(srcType)
		for _, v := range srcValue.MapKeys() {
			p := reflect.New(srcValue.MapIndex(v).Type())
			recursionGet2(srcValue.MapIndex(v).Interface(), p.Interface())
			newMap.SetMapIndex(v, p.Elem())
		}
		result = newMap
	default:
		result = srcValue
	}

	desElem.Set(result)

	return nil
}
