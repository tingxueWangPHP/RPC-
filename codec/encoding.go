package codec

import (
	"encoding/gob"
	"io"
	"reflect"
)

func Encode(conn io.Writer, data DataFormat) error {
	//判断是否要注册
	for _, item := range data.Body {
		//if reflect.TypeOf(item).Kind() == reflect.Struct || reflect.TypeOf(item).Kind() == reflect.Map {
		gob.Register(reflect.ValueOf(item).Interface())
		//}
	}
	encoder := gob.NewEncoder(conn)

	return encoder.Encode(data)
}

func Decode(conn io.Reader) (*DataFormat, error) {
	d := new(DataFormat)
	decoder := gob.NewDecoder(conn)
	if err := decoder.Decode(d); err != nil {
		return nil, err
	}

	return d, nil
}
