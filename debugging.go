package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"image"
	"image/png"
	"strconv"
)

type ImageHelper struct {
	m     *Map
	pixel func(row, col int) image.NRGBAColor
}

func (ih ImageHelper) ColorModel() image.ColorModel {
	return image.NRGBAColorModel
}

func (ih ImageHelper) Bounds() image.Rectangle {
	fmt.Printf("Bounds are 0,0,%d,%d\n",ih.m.Cols*4, ih.m.Rows*4)
	return image.Rect(0, 0, ih.m.Cols*4, ih.m.Rows*4)
}

func (ih ImageHelper) At(x, y int) image.Color {
	return ih.pixel(y/4, x/4)
}


//implement Image for fancy image debugging
func (m *Map) ColorModel() image.ColorModel {
	return image.NRGBAColorModel
}
func (m *Map) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.Cols*4, m.Rows*4)
}
func (m *Map) At(x, y int) image.NRGBAColor {
	loc := m.ToLocation(Point{y, x})
	return m.Grid[loc].Color()
}


func (m *Map) WriteDebugImage(Desc string, seq int, At func(row, col int) image.NRGBAColor) {

	//use -imgprefix="bot0" to make a series of images (bot0.0.png ... bot0.N.png) which
	//illustrate the bot's knowledge of the map at each turn. If you want the images in a
	//subdirectory, make sure you create the directory first. (e.g., -imgprefix="images/bot0")
	imageOutPrefix := flag.String("imgprefix", "bot", "prefix for helpful debugging images")
	if imageOutPrefix == nil {
		return
	}

	fname := fmt.Sprintf("%s.%s.%3.3d.png", *imageOutPrefix, Desc, seq)
	fmt.Printf("making image: %s\n", fname)

	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Panicf("Couldn't open %s (%s)", fname, err)
	}
	defer f.Close()

	err = png.Encode(f, ImageHelper{m, At})

	if err != nil {
		log.Panicf("Couldn't encode png (%s)", err)
	}
}





func (s *State) DumpSeen() {
	mseen := Max(s.Map.Seen)
	str := ""

	for r := 0; r < s.Rows; r++ {
		for c := 0; c < s.Cols; c++ {
			str += strconv.Itoa(s.Map.Seen[r*s.Cols+c] * 10 / (mseen + 1))
		}
		str += "\n"
	}

	log.Printf("Turn %d\n%v\n", s.Turn, str)
}

func (s *State) DumpMap() {
	m := make([]byte, len(s.Map.Grid))
	str := ""

	for i, o := range s.Map.Grid {
		m[i] = o.ToSymbol()
	}

	for r := 0; r < s.Rows; r++ {
		for c := 0; c < s.Cols; c++ {
			str += string(m[r*s.Cols+c])
		}
		str += "\n"
	}

	log.Printf("Turn %d\n%v\n", s.Turn, str)
}
