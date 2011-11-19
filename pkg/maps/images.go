package maps

import (
	"fmt"
	"flag"
	"log"
	"os"
	"image"
	"image/png"
)

type ImageHelper struct {
	m     *Map
	pixel func(row, col int) image.NRGBAColor
}

func (ih ImageHelper) ColorModel() image.ColorModel {
	return image.NRGBAColorModel
}

func (ih ImageHelper) Bounds() image.Rectangle {
	// fmt.Printf("Bounds are 0,0,%d,%d\n", ih.m.Cols*4, ih.m.Rows*4)
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
	return ItemColor(m.Grid[loc])
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
func ItemColor(o Item) image.NRGBAColor {
	switch o {
	case UNKNOWN:
		return image.NRGBAColor{0xb0, 0xb0, 0xb0, 0xff}
	case WATER:
		return image.NRGBAColor{0x10, 0x10, 0x50, 0xff}
	case FOOD:
		return image.NRGBAColor{0xe0, 0xe0, 0xc0, 0xff}
	case LAND:
		return image.NRGBAColor{0x8b, 0x45, 0x13, 0xff}
	case MY_ANT:
		return image.NRGBAColor{0xf0, 0x00, 0x00, 0xff}
	case PLAYER1:
		return image.NRGBAColor{0xf0, 0xf0, 0x00, 0xff}
	case PLAYER2:
		return image.NRGBAColor{0x00, 0xf0, 0x00, 0xff}
	case PLAYER3:
		return image.NRGBAColor{0x00, 0x00, 0xf0, 0xff}
	case PLAYER4:
		return image.NRGBAColor{0xf0, 0x00, 0xf0, 0xff}
	case PLAYER5:
		return image.NRGBAColor{0xf0, 0xf0, 0xf0, 0xff}
	case PLAYER6:
		return image.NRGBAColor{0x80, 0x80, 0x00, 0xff}
	case PLAYER7:
		return image.NRGBAColor{0x00, 0x80, 0x80, 0xff}
	case PLAYER8:
		return image.NRGBAColor{0x80, 0x00, 0x80, 0xff}
	case PLAYER9:
		return image.NRGBAColor{0x80, 0x00, 0xf0, 0xff}
	}
	return image.NRGBAColor{0xff, 0xff, 0xff, 0xff}
}
