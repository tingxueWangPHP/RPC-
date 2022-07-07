package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

func main() {
	/*var buf bytes.Buffer

	coder := gob.NewEncoder(&buf)
	_ = coder.Encode(5)



	/*decoder := gob.NewDecoder(&buf)

	var a int
	decoder.Decode(&a)

	fmt.Println(a)*/

	f, err := os.Create("./2.txt")
	if err != nil {
		panic(err)
	}

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(5); err != nil {
		panic(err)
	}
	defer f.Close()

	f, err = os.Open("./2.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	decoder := gob.NewDecoder(f)

	var a int
	decoder.Decode(&a)

	fmt.Println(a)
}
