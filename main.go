package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
)

type Header struct {
	byteOrder byte
	magic     [2]byte
	version   byte
	size      [2]byte
	length    [2]byte
}

type Maze struct {
	data []byte
}

func sliceEquals(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}

	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}

func (m *Maze) UnmarshalBinary(data []byte) error {
	byteOrder := data[0]
	magic := data[1:3]

	if !sliceEquals(magic, []byte{'\x5D', '\x90'}) {
		return errors.New("Magic number incorrect")
	}

	v := data[3]
	if v != 1 {
		return errors.New("Unsupported version")
	}

	// width := data[4]
	// heigth := data[5]

	var length uint16

	fmt.Println(data[6:8])

	if byteOrder == 1 {
		// use little endian
		length = binary.LittleEndian.Uint16(data[6:8])
		fmt.Println("Using little endian, data is ", length, " bytes long")
	} else {
		length = binary.BigEndian.Uint16(data[6:8])
		fmt.Println("Using big endian, data is ", length, " bytes long")
	}

	return nil

	// version :=
}

func main() {

	data, err := ioutil.ReadFile("maze.map")
	if err != nil {
		panic(err)
	}

	m := Maze{}
	m.UnmarshalBinary(data)

	fmt.Println()
}
