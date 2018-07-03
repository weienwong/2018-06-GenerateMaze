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
}

type Maze struct {
	data   [][]bool
	startx int
	starty int
	items  []item
}

type item interface {
	location() (int, int)
	symbol() byte
}

type goal struct {
	x, y     int
	required bool
}

type warp struct {
	x, y, destx, desty int
}

func (g goal) location() (int, int) {
	return g.x, g.y
}

func (g goal) symbol() byte {
	if g.required {
		return 'x'
	}
	return 'o'
}

func (w warp) location() (int, int) {
	return w.x, w.y
}

func (w warp) symbol() byte {
	return 'w'
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

	width := int(data[4])
	height := int(data[5])

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

	c := make(chan bool, 10)

	go func() {
		m.data = make([][]bool, height)

		for i := 0; i < height; i++ {
			m.data[i] = make([]bool, width)
			for j := 0; j < width; j++ {
				m.data[i][j] = <-c
			}
		}
	}()

	for _, v := range data[8 : 8+length] {
		fmt.Printf("%b\n", v)
		for i := uint(0); i < 8; i++ {
			c <- v&(128>>i) > 0
			fmt.Println(v&(128>>i) > 0)
		}
	}

	data = data[8+length:]

	m.startx = int(data[0])
	m.starty = int(data[1])

	itemCount := int(data[2])

	m.items = make([]item, itemCount)
	cursor := 3

	for index := 0; index < itemCount; index++ {
		t := data[cursor]

		switch t {
		case 0:
			// required goal
			m.items[index] = goal{x: int(data[t+1]), y: int(data[t+2]), required: true}
			fmt.Println("Required goal")
			cursor += 3
		case 1:
			// optional goal
			m.items[index] = goal{x: int(data[t+1]), y: int(data[t+2]), required: false}
			cursor += 3
		case 2:
			// warp
			m.items[index] = warp{
				x:     int(data[t+1]),
				y:     int(data[t+2]),
				destx: int(data[t+3]),
				desty: int(data[t+4]),
			}
			cursor += 5
		}

	}

	return nil

}

func main() {

	data, err := ioutil.ReadFile("maze.map")
	if err != nil {
		panic(err)
	}

	m := Maze{}
	m.UnmarshalBinary(data)

	fmt.Println(m)
}
