package maps

import (
	"os"
	"io/ioutil"
	"strings"
)

var AllMaps = []string{}
var MapRoot string

func init() {
	MapRoot = os.Getenv("HOME") + "/bot/src/pkg/maps/testdata/maps"
	m, err := ioutil.ReadDir(MapRoot)
	if err == nil {
		for _, f := range m {
			if strings.HasSuffix(f.Name, ".map") {
				AllMaps = append(AllMaps, f.Name[:len(f.Name)-4])
			}
		}
	}
}

func MapFile(name string) string {
	file := MapRoot + "/" + name + ".map"
	return file
}
