package optpfd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	//"math"
)

func Compress(v []int) ([]byte, error) {
	if len(v) == 0 {
		return nil, errors.New("no data to compress.")
	}
	var b bytes.Buffer
	segmentSize := 128
	for len(v) > segmentSize {
		pre := v[0:segmentSize]
		tmp, err := segmentCompress(pre)
		if err != nil {
			return nil, err
		}
		b.Write(tmp)
		v = v[segmentSize:]
	}
	return nil, nil
}

// segmentCompress compress v []int
// 0		the lenght of v, (0,128]
// 1 - 4	the first integer
// 5 - 8	the min delta
// 9		x bit for majority
// 10		y bit for all
// 11   	z the number of integer that can be expressed in x bit.
// 			x bit for every interger
//			the rest part for exceptional intergers
// end
func segmentCompress(v []int) ([]byte, error) {
	var b bytes.Buffer
	//byte 0
	b.WriteByte(byte(len(v)))
	//byte 1-4 first
	err := binary.Write(&b, binary.LittleEndian, int32(v[0]))
	if err != nil {
		return nil, err
	}
	if len(v) == 1 {
		return b.Bytes(), nil
	}
	// delta
	pre := v[0]
	min := v[1] - v[0]
	max := min
	for i := 1; i < len(v); i++ {
		delta := v[i] - pre
		pre = v[i]
		v[i] = delta
		if delta < min {
			min = delta
		}
		if delta > max {
			max = delta
		}
	}
	for i := 1; i < len(v); i++ {
		v[i] -= min
	}
	//byte 5-8 min
	err = binary.Write(&b, binary.LittleEndian, int32(min))
	if err != nil {
		return nil, err
	}
	//count
	c := [32]int{}
	for i := 1; i < len(v); i++ {
		for j := 0; j < len(pow); j++ {
			if uint32(v[i]) <= pow[j] {
				c[j] += 1
			}
		}
	}
	//choose y, max-min, 可以确定所有的数，使用y位都可以保存
	y := 32
	for j := len(pow) - 1; j >= 0; j-- {
		if pow[j] >= uint32(max-min) {
			y = j + 1
			continue
		}
		break
	}
	//choose x
	numberOfBits := 32 * (len(v) - 1)
	x := 32
	for j := 0; j < len(c); j++ {
		tmp := (j+1)*(len(v)-1) + (len(v)-1-c[j])*(8+y-j-1)
		if tmp < numberOfBits {
			numberOfBits = tmp
			x = j + 1
		}
	}
	fmt.Printf("choose x: %d,y: %d, num %d,min: %d, max %d,  compress ratio %f\n", x, y, c[x-1], min, max, float64(numberOfBits)/float64((32*(len(v)-1))))
	// byte 9 x
	b.WriteByte(byte(x))
	// byte 10 y
	b.WriteByte(byte(y))
	// byte 11 z
	b.WriteByte(byte(c[x-1]))
	// write x bit for every interger
	//xbits := getXpart(v, x)
	return b.Bytes(), nil
}

func getXpart(v []int, x int) []byte {
	cnt := len(v) - 1
	size := int(math.Ceil(float64(x*cnt) / float64(8)))
	buf := make([]byte, size)
	mask := mask(x)
	bitIndex, byteIndex := 0, 0
	for i := 1; i < len(v); i++ {
		t := v[i] & mask
		put(i, t, buf, byteIndex, bitIndex)
	}
	return buf
}

func put(i, t int, b []byte, byteIndex, bitIndex int) {

}

func mask(x int) int {
	t := 0
	for i := 0; i < x; i++ {
		t |= (1 << uint32(i))
	}
	return t
}

var pow []uint32 = []uint32{1, 3, 7, 15, 31, 63, 127, 255, 511, 1023, 2047, 4095, 8191, 16383, 32767, 65535, 131071, 262143, 524287, 1048575, 2097151, 4194303, 8388607, 16777215, 33554431, 67108863, 134217727, 268435455, 536870911, 1073741823, 2147483647, 4294967295}

func Decompress(b []byte) ([]int, error) {
	return nil, nil
}
