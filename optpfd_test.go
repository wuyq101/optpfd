package optpfd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestCompress(t *testing.T) {
	v, err := loadFromFile("example/600874.txt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("load %d integers\n", len(v))
	Compress(v)
}

func loadFromFile(file string) ([]int, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	bs := bytes.Split(b, []byte("\n"))
	v := make([]int, 0)
	for _, s := range bs {
		n, err := strconv.Atoi(string(s))
		if err != nil {
			return nil, err
		}
		v = append(v, n)
	}
	return v, nil
}

func TestMask(t *testing.T) {
	for i := 1; i <= 32; i++ {
		m := mask(i)
		fmt.Printf("%b \n", m)
	}
}
