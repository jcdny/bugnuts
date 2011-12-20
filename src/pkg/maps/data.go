// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package maps

import (
	"os"
	"io/ioutil"
	"strings"
	. "bugnuts/torus"
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

type MapSymData struct {
	Cycle            int
	Equivlen         int
	Gen              int
	Label            string
	Tiles            int
	Rows             int
	Cols             int
	Translate        Point
	RM1, RM2, MR, MC int
	Origin           Point
}

var MapSym = map[string]MapSymData{
	"maze_p02_26":        MapSymData{2, 2, 6, "NDiag", 1, 100, 100, Point{0, 0}, 23, -1, -1, -1, Point{0, 0}},
	"maze_08p_01":        MapSymData{8, 8, 12, "Rota8", 1, 120, 120, Point{0, 0}, -1, -1, -1, -1, Point{37, 84}},
	"maze_p08_02":        MapSymData{8, 8, 12, "Rota8", 1, 108, 108, Point{0, 0}, -1, -1, -1, -1, Point{28, 94}},
	"maze_p08_03":        MapSymData{8, 8, 12, "Rota8", 1, 110, 110, Point{0, 0}, -1, -1, -1, -1, Point{41, 62}},
	"maze_p08_04":        MapSymData{8, 8, 12, "Rota8", 1, 110, 110, Point{0, 0}, -1, -1, -1, -1, Point{4, 90}},
	"maze_p08_05":        MapSymData{8, 8, 12, "Rota8", 1, 150, 150, Point{0, 0}, -1, -1, -1, -1, Point{24, 36}},
	"maze_p08_06":        MapSymData{8, 8, 12, "Rota8", 1, 120, 120, Point{0, 0}, -1, -1, -1, -1, Point{36, 43}},
	"maze_p08_07":        MapSymData{8, 8, 12, "Rota8", 1, 120, 120, Point{0, 0}, -1, -1, -1, -1, Point{39, 110}},
	"mmaze_08p_01":       MapSymData{8, 8, 12, "Rota8", 1, 128, 128, Point{0, 0}, -1, -1, -1, -1, Point{17, 17}},
	"maze_p04_05":        MapSymData{4, 4, 10, "DBoth", 1, 114, 114, Point{0, 0}, 32, 37, -1, -1, Point{0, 0}},
	"maze_p04_06":        MapSymData{4, 4, 10, "DBoth", 1, 112, 112, Point{0, 0}, 78, 23, -1, -1, Point{0, 0}},
	"maze_p04_16":        MapSymData{4, 4, 10, "DBoth", 1, 100, 100, Point{0, 0}, 47, 10, -1, -1, Point{0, 0}},
	"maze_p04_18":        MapSymData{4, 4, 10, "DBoth", 1, 120, 120, Point{0, 0}, 87, 110, -1, -1, Point{0, 0}},
	"maze_p04_19":        MapSymData{4, 4, 10, "DBoth", 1, 114, 114, Point{0, 0}, 17, 100, -1, -1, Point{0, 0}},
	"maze_p04_22":        MapSymData{4, 4, 10, "DBoth", 1, 132, 132, Point{0, 0}, 37, 114, -1, -1, Point{0, 0}},
	"maze_p04_28":        MapSymData{4, 4, 10, "DBoth", 1, 96, 96, Point{0, 0}, 79, 92, -1, -1, Point{0, 0}},
	"maze_p04_31":        MapSymData{4, 4, 10, "DBoth", 1, 96, 96, Point{0, 0}, 83, 26, -1, -1, Point{0, 0}},
	"maze_p04_33":        MapSymData{4, 4, 10, "DBoth", 1, 64, 64, Point{0, 0}, 19, 48, -1, -1, Point{0, 0}},
	"maze_02p_02":        MapSymData{2, 2, 7, "PDiag", 1, 72, 72, Point{0, 0}, -1, 16, -1, -1, Point{0, 0}},
	"maze_p02_02":        MapSymData{2, 2, 7, "PDiag", 1, 80, 80, Point{0, 0}, -1, 3, -1, -1, Point{0, 0}},
	"maze_p02_05":        MapSymData{2, 2, 7, "PDiag", 1, 66, 66, Point{0, 0}, -1, 54, -1, -1, Point{0, 0}},
	"maze_p02_12":        MapSymData{2, 2, 7, "PDiag", 1, 150, 150, Point{0, 0}, -1, 24, -1, -1, Point{0, 0}},
	"maze_p02_16":        MapSymData{2, 2, 7, "PDiag", 1, 88, 88, Point{0, 0}, -1, 23, -1, -1, Point{0, 0}},
	"maze_p02_19":        MapSymData{2, 2, 7, "PDiag", 1, 120, 120, Point{0, 0}, -1, 94, -1, -1, Point{0, 0}},
	"maze_p02_20":        MapSymData{2, 2, 7, "PDiag", 1, 120, 120, Point{0, 0}, -1, 117, -1, -1, Point{0, 0}},
	"maze_p02_22":        MapSymData{2, 2, 7, "PDiag", 1, 150, 150, Point{0, 0}, -1, 17, -1, -1, Point{0, 0}},
	"maze_p02_23":        MapSymData{2, 2, 7, "PDiag", 1, 120, 120, Point{0, 0}, -1, 104, -1, -1, Point{0, 0}},
	"maze_p02_24":        MapSymData{2, 2, 7, "PDiag", 1, 110, 110, Point{0, 0}, -1, 100, -1, -1, Point{0, 0}},
	"maze_p02_01":        MapSymData{2, 2, 6, "NDiag", 1, 126, 126, Point{0, 0}, 64, -1, -1, -1, Point{0, 0}},
	"maze_p02_03":        MapSymData{2, 2, 6, "NDiag", 1, 84, 84, Point{0, 0}, 3, -1, -1, -1, Point{0, 0}},
	"maze_p02_04":        MapSymData{2, 2, 6, "NDiag", 1, 84, 84, Point{0, 0}, 37, -1, -1, -1, Point{0, 0}},
	"maze_p02_10":        MapSymData{2, 2, 6, "NDiag", 1, 150, 150, Point{0, 0}, 69, -1, -1, -1, Point{0, 0}},
	"maze_p02_11":        MapSymData{2, 2, 6, "NDiag", 1, 80, 80, Point{0, 0}, 69, -1, -1, -1, Point{0, 0}},
	"maze_p02_13":        MapSymData{2, 2, 6, "NDiag", 1, 108, 108, Point{0, 0}, 26, -1, -1, -1, Point{0, 0}},
	"maze_p02_15":        MapSymData{2, 2, 6, "NDiag", 1, 96, 96, Point{0, 0}, 82, -1, -1, -1, Point{0, 0}},
	"maze_p02_17":        MapSymData{2, 2, 6, "NDiag", 1, 96, 96, Point{0, 0}, 86, -1, -1, -1, Point{0, 0}},
	"maze_p02_25":        MapSymData{2, 2, 6, "NDiag", 1, 78, 78, Point{0, 0}, 26, -1, -1, -1, Point{0, 0}},
	"maze_p02_28":        MapSymData{2, 2, 6, "NDiag", 1, 108, 108, Point{0, 0}, 30, -1, -1, -1, Point{0, 0}},
	"maze_p02_29":        MapSymData{2, 2, 6, "NDiag", 1, 90, 90, Point{0, 0}, 80, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_02":   MapSymData{4, 4, 4, "R_180", 2, 180, 96, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_07":   MapSymData{4, 4, 4, "R_180", 2, 52, 90, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_08":   MapSymData{4, 4, 4, "R_180", 2, 56, 88, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_17":   MapSymData{4, 4, 4, "R_180", 2, 24, 128, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_18":   MapSymData{4, 4, 4, "R_180", 2, 76, 100, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_02":   MapSymData{6, 6, 4, "R_180", 3, 114, 58, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_07":   MapSymData{6, 6, 4, "R_180", 3, 90, 66, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_09":   MapSymData{6, 6, 4, "R_180", 3, 60, 36, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_10":   MapSymData{6, 6, 4, "R_180", 3, 114, 66, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_11":   MapSymData{6, 6, 4, "R_180", 3, 84, 66, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_12":   MapSymData{6, 6, 4, "R_180", 3, 42, 54, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_14":   MapSymData{6, 6, 4, "R_180", 3, 72, 54, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_16":   MapSymData{6, 6, 4, "R_180", 3, 114, 66, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p08_03":   MapSymData{8, 8, 4, "R_180", 4, 48, 38, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p10_01":   MapSymData{10, 10, 4, "R_180", 5, 50, 34, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p10_03":   MapSymData{10, 10, 4, "R_180", 5, 60, 40, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p02_07":        MapSymData{2, 2, 4, "R_180", 1, 120, 128, Point{0, 0}, -1, -1, -1, -1, Point{4, 58}},
	"maze_p02_08":        MapSymData{2, 2, 4, "R_180", 1, 96, 108, Point{0, 0}, -1, -1, -1, -1, Point{2, 19}},
	"maze_p02_14":        MapSymData{2, 2, 4, "R_180", 1, 72, 84, Point{0, 0}, -1, -1, -1, -1, Point{24, 24}},
	"maze_p02_30":        MapSymData{2, 2, 4, "R_180", 1, 80, 150, Point{0, 0}, -1, -1, -1, -1, Point{33, 142}},
	"maze_p02_33":        MapSymData{2, 2, 4, "R_180", 1, 110, 140, Point{0, 0}, -1, -1, -1, -1, Point{24, 137}},
	"maze_p04_01":        MapSymData{4, 4, 3, "Rotat", 1, 98, 98, Point{0, 0}, -1, -1, -1, -1, Point{31, 94}},
	"maze_p04_02":        MapSymData{4, 4, 3, "Rotat", 1, 108, 108, Point{0, 0}, -1, -1, -1, -1, Point{2, 58}},
	"maze_p04_03":        MapSymData{4, 4, 3, "Rotat", 1, 80, 80, Point{0, 0}, -1, -1, -1, -1, Point{36, 55}},
	"maze_p04_04":        MapSymData{4, 4, 3, "Rotat", 1, 108, 108, Point{0, 0}, -1, -1, -1, -1, Point{51, 102}},
	"maze_p04_07":        MapSymData{4, 4, 3, "Rotat", 1, 90, 90, Point{0, 0}, -1, -1, -1, -1, Point{44, 22}},
	"maze_p04_08":        MapSymData{4, 4, 3, "Rotat", 1, 104, 104, Point{0, 0}, -1, -1, -1, -1, Point{43, 66}},
	"maze_p04_09":        MapSymData{4, 4, 3, "Rotat", 1, 84, 84, Point{0, 0}, -1, -1, -1, -1, Point{7, 75}},
	"maze_p04_10":        MapSymData{4, 4, 3, "Rotat", 1, 100, 100, Point{0, 0}, -1, -1, -1, -1, Point{17, 36}},
	"maze_p04_12":        MapSymData{4, 4, 3, "Rotat", 1, 78, 78, Point{0, 0}, -1, -1, -1, -1, Point{36, 9}},
	"maze_p04_13":        MapSymData{4, 4, 3, "Rotat", 1, 90, 90, Point{0, 0}, -1, -1, -1, -1, Point{16, 71}},
	"maze_p04_15":        MapSymData{4, 4, 3, "Rotat", 1, 96, 96, Point{0, 0}, -1, -1, -1, -1, Point{1, 34}},
	"maze_p04_17":        MapSymData{4, 4, 3, "Rotat", 1, 78, 78, Point{0, 0}, -1, -1, -1, -1, Point{22, 6}},
	"maze_p04_20":        MapSymData{4, 4, 3, "Rotat", 1, 104, 104, Point{0, 0}, -1, -1, -1, -1, Point{38, 89}},
	"maze_p04_23":        MapSymData{4, 4, 3, "Rotat", 1, 130, 130, Point{0, 0}, -1, -1, -1, -1, Point{14, 58}},
	"maze_p04_24":        MapSymData{4, 4, 3, "Rotat", 1, 140, 140, Point{0, 0}, -1, -1, -1, -1, Point{33, 46}},
	"maze_p04_25":        MapSymData{4, 4, 3, "Rotat", 1, 144, 144, Point{0, 0}, -1, -1, -1, -1, Point{60, 82}},
	"maze_p04_26":        MapSymData{4, 4, 3, "Rotat", 1, 84, 84, Point{0, 0}, -1, -1, -1, -1, Point{25, 67}},
	"maze_p04_27":        MapSymData{4, 4, 3, "Rotat", 1, 72, 72, Point{0, 0}, -1, -1, -1, -1, Point{0, 22}},
	"maze_p04_30":        MapSymData{4, 4, 3, "Rotat", 1, 96, 96, Point{0, 0}, -1, -1, -1, -1, Point{18, 76}},
	"maze_p04_32":        MapSymData{4, 4, 3, "Rotat", 1, 112, 112, Point{0, 0}, -1, -1, -1, -1, Point{48, 29}},
	"maze_p04_35":        MapSymData{4, 4, 3, "Rotat", 1, 104, 104, Point{0, 0}, -1, -1, -1, -1, Point{37, 64}},
	"maze_p04_37":        MapSymData{4, 4, 3, "Rotat", 1, 140, 140, Point{0, 0}, -1, -1, -1, -1, Point{44, 133}},
	"mmaze_04p_02":       MapSymData{4, 4, 3, "Rotat", 1, 120, 120, Point{0, 0}, -1, -1, -1, -1, Point{26, 85}},
	"cell_maze_p06_01":   MapSymData{6, 6, 2, "MirrR", 3, 180, 66, Point{0, 0}, -1, -1, 0, -1, Point{0, 0}},
	"cell_maze_p06_05":   MapSymData{6, 6, 2, "MirrR", 3, 72, 66, Point{0, 0}, -1, -1, 0, -1, Point{0, 0}},
	"cell_maze_p06_08":   MapSymData{6, 6, 2, "MirrR", 3, 114, 66, Point{0, 0}, -1, -1, 0, -1, Point{0, 0}},
	"cell_maze_p06_13":   MapSymData{6, 6, 2, "MirrR", 3, 120, 66, Point{0, 0}, -1, -1, 0, -1, Point{0, 0}},
	"cell_maze_p08_06":   MapSymData{8, 8, 2, "MirrR", 4, 56, 50, Point{0, 0}, -1, -1, 0, -1, Point{0, 0}},
	"maze_p02_31":        MapSymData{2, 2, 2, "MirrR", 1, 80, 120, Point{0, 0}, -1, -1, 10, -1, Point{0, 0}},
	"cell_maze_p02_01":   MapSymData{2, 2, 0, "Trans", 2, 134, 98, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_02":   MapSymData{2, 2, 0, "Trans", 2, 74, 180, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_03":   MapSymData{2, 2, 0, "Trans", 1, 130, 176, Point{65, 88}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_04":   MapSymData{2, 2, 0, "Trans", 1, 110, 114, Point{55, 57}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_05":   MapSymData{2, 2, 0, "Trans", 1, 136, 198, Point{68, 99}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_06":   MapSymData{2, 2, 0, "Trans", 2, 20, 48, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_07":   MapSymData{2, 2, 0, "Trans", 1, 68, 86, Point{34, 43}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_09":   MapSymData{2, 2, 0, "Trans", 1, 70, 134, Point{35, 67}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_12":   MapSymData{2, 2, 0, "Trans", 1, 28, 116, Point{14, 58}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_13":   MapSymData{2, 2, 0, "Trans", 1, 80, 110, Point{40, 55}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_15":   MapSymData{2, 2, 0, "Trans", 2, 43, 106, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_16":   MapSymData{2, 2, 0, "Trans", 2, 50, 92, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p02_17":   MapSymData{2, 2, 0, "Trans", 2, 48, 54, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_01":   MapSymData{3, 3, 0, "Trans", 1, 102, 129, Point{34, 43}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_02":   MapSymData{3, 3, 0, "Trans", 1, 174, 183, Point{58, -61}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_03":   MapSymData{3, 3, 0, "Trans", 1, 189, 192, Point{63, -64}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_04":   MapSymData{3, 3, 0, "Trans", 1, 66, 111, Point{22, 37}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_05":   MapSymData{3, 3, 0, "Trans", 3, 39, 37, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_06":   MapSymData{3, 3, 0, "Trans", 1, 78, 129, Point{26, 43}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_07":   MapSymData{3, 3, 0, "Trans", 1, 126, 129, Point{42, 43}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_08":   MapSymData{3, 3, 0, "Trans", 3, 63, 66, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_09":   MapSymData{3, 3, 0, "Trans", 1, 84, 171, Point{28, 57}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_10":   MapSymData{3, 3, 0, "Trans", 1, 30, 138, Point{10, 46}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_11":   MapSymData{3, 3, 0, "Trans", 1, 102, 108, Point{34, -36}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_12":   MapSymData{3, 3, 0, "Trans", 1, 87, 120, Point{29, -40}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_13":   MapSymData{3, 3, 0, "Trans", 1, 75, 162, Point{25, -54}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_14":   MapSymData{3, 3, 0, "Trans", 1, 66, 165, Point{22, 55}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_15":   MapSymData{3, 3, 0, "Trans", 1, 60, 198, Point{20, 66}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_16":   MapSymData{3, 3, 0, "Trans", 1, 45, 168, Point{15, -56}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_17":   MapSymData{3, 3, 0, "Trans", 3, 42, 45, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_18":   MapSymData{3, 3, 0, "Trans", 1, 45, 195, Point{15, -65}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_19":   MapSymData{3, 3, 0, "Trans", 3, 48, 32, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p03_20":   MapSymData{3, 3, 0, "Trans", 1, 36, 150, Point{12, -50}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_01":   MapSymData{4, 4, 0, "Trans", 2, 196, 98, Point{98, 49}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_03":   MapSymData{4, 4, 0, "Trans", 2, 44, 42, Point{22, 21}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_04":   MapSymData{4, 4, 0, "Trans", 2, 52, 82, Point{26, 41}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_06":   MapSymData{4, 4, 0, "Trans", 4, 22, 88, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_10":   MapSymData{4, 4, 0, "Trans", 1, 72, 128, Point{18, -32}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_12":   MapSymData{4, 4, 0, "Trans", 1, 92, 200, Point{23, -50}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p04_13":   MapSymData{4, 4, 0, "Trans", 4, 68, 39, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_01":   MapSymData{5, 5, 0, "Trans", 1, 90, 165, Point{36, 33}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_02":   MapSymData{5, 5, 0, "Trans", 1, 145, 155, Point{29, -31}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_03":   MapSymData{5, 5, 0, "Trans", 1, 145, 190, Point{29, 38}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_04":   MapSymData{5, 5, 0, "Trans", 1, 50, 90, Point{20, -18}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_05":   MapSymData{5, 5, 0, "Trans", 1, 80, 120, Point{32, -24}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_06":   MapSymData{5, 5, 0, "Trans", 5, 115, 25, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_07":   MapSymData{5, 5, 0, "Trans", 1, 85, 155, Point{17, 31}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_08":   MapSymData{5, 5, 0, "Trans", 1, 125, 185, Point{25, 37}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_09":   MapSymData{5, 5, 0, "Trans", 1, 60, 120, Point{24, 24}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_10":   MapSymData{5, 5, 0, "Trans", 1, 55, 155, Point{22, 31}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_11":   MapSymData{5, 5, 0, "Trans", 1, 55, 200, Point{11, 40}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_12":   MapSymData{5, 5, 0, "Trans", 1, 80, 200, Point{32, -40}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_13":   MapSymData{5, 5, 0, "Trans", 1, 105, 200, Point{21, 40}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_14":   MapSymData{5, 5, 0, "Trans", 1, 105, 200, Point{42, -40}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_15":   MapSymData{5, 5, 0, "Trans", 1, 70, 130, Point{14, 26}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_16":   MapSymData{5, 5, 0, "Trans", 1, 40, 175, Point{16, -35}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_17":   MapSymData{5, 5, 0, "Trans", 1, 45, 135, Point{9, 27}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_18":   MapSymData{5, 5, 0, "Trans", 1, 90, 110, Point{36, 22}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_19":   MapSymData{5, 5, 0, "Trans", 1, 110, 200, Point{44, -40}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p05_20":   MapSymData{5, 5, 0, "Trans", 1, 120, 170, Point{24, 34}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_03":   MapSymData{6, 6, 0, "Trans", 3, 48, 40, Point{24, 20}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_04":   MapSymData{6, 6, 0, "Trans", 2, 39, 144, Point{13, -48}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_06":   MapSymData{6, 6, 0, "Trans", 3, 96, 66, Point{48, 33}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_15":   MapSymData{6, 6, 0, "Trans", 2, 96, 99, Point{32, 33}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p06_18":   MapSymData{6, 6, 0, "Trans", 3, 48, 66, Point{24, 33}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_01":   MapSymData{7, 7, 0, "Trans", 1, 126, 168, Point{18, -48}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_02":   MapSymData{7, 7, 0, "Trans", 1, 91, 168, Point{26, -24}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_03":   MapSymData{7, 7, 0, "Trans", 1, 49, 112, Point{21, 16}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_04":   MapSymData{7, 7, 0, "Trans", 7, 56, 25, Point{0, 0}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_05":   MapSymData{7, 7, 0, "Trans", 1, 119, 196, Point{34, -28}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_06":   MapSymData{7, 7, 0, "Trans", 1, 63, 98, Point{18, 14}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_07":   MapSymData{7, 7, 0, "Trans", 1, 84, 189, Point{24, -27}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_08":   MapSymData{7, 7, 0, "Trans", 1, 119, 189, Point{17, -54}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_09":   MapSymData{7, 7, 0, "Trans", 1, 105, 196, Point{30, -28}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_10":   MapSymData{7, 7, 0, "Trans", 1, 77, 105, Point{11, -30}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_11":   MapSymData{7, 7, 0, "Trans", 1, 119, 189, Point{34, 27}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p07_12":   MapSymData{7, 7, 0, "Trans", 1, 119, 196, Point{17, -56}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p08_01":   MapSymData{8, 8, 0, "Trans", 4, 32, 40, Point{16, 20}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p08_02":   MapSymData{8, 8, 0, "Trans", 2, 96, 80, Point{24, -20}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p08_07":   MapSymData{8, 8, 0, "Trans", 1, 56, 184, Point{21, -23}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p08_08":   MapSymData{8, 8, 0, "Trans", 4, 28, 128, Point{14, 64}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p09_01":   MapSymData{9, 9, 0, "Trans", 1, 99, 189, Point{11, -42}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p09_02":   MapSymData{9, 9, 0, "Trans", 3, 72, 39, Point{24, 13}, -1, -1, -1, -1, Point{0, 0}},
	"cell_maze_p09_03":   MapSymData{9, 9, 0, "Trans", 3, 117, 66, Point{39, 22}, -1, -1, -1, -1, Point{0, 0}},
	"maze_02p_01":        MapSymData{2, 2, 0, "Trans", 1, 60, 96, Point{30, 48}, -1, -1, -1, -1, Point{0, 0}},
	"maze_03p_01":        MapSymData{3, 3, 0, "Trans", 1, 96, 96, Point{32, -32}, -1, -1, -1, -1, Point{0, 0}},
	"maze_05p_01":        MapSymData{5, 5, 0, "Trans", 1, 120, 120, Point{24, -48}, -1, -1, -1, -1, Point{0, 0}},
	"maze_06p_01":        MapSymData{6, 6, 0, "Trans", 1, 144, 144, Point{24, 24}, -1, -1, -1, -1, Point{0, 0}},
	"maze_07p_01":        MapSymData{7, 7, 0, "Trans", 1, 126, 126, Point{36, 18}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p02_18":        MapSymData{2, 2, 0, "Trans", 1, 80, 120, Point{40, 60}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p02_27":        MapSymData{2, 2, 0, "Trans", 1, 72, 144, Point{36, 72}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p03_01":        MapSymData{3, 3, 0, "Trans", 1, 72, 126, Point{24, 42}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p03_02":        MapSymData{3, 3, 0, "Trans", 1, 96, 96, Point{32, 32}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p03_03":        MapSymData{3, 3, 0, "Trans", 1, 144, 144, Point{48, 48}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p03_04":        MapSymData{3, 3, 0, "Trans", 1, 90, 108, Point{30, 36}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p04_11":        MapSymData{4, 4, 0, "Trans", 1, 96, 128, Point{24, 32}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p04_29":        MapSymData{4, 4, 0, "Trans", 1, 72, 96, Point{18, -24}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p05_01":        MapSymData{5, 5, 0, "Trans", 1, 100, 100, Point{20, -40}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p05_02":        MapSymData{5, 5, 0, "Trans", 1, 120, 120, Point{24, 24}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p05_03":        MapSymData{5, 5, 0, "Trans", 1, 120, 120, Point{24, -24}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p05_04":        MapSymData{5, 5, 0, "Trans", 1, 80, 120, Point{16, 24}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p06_01":        MapSymData{6, 6, 0, "Trans", 1, 108, 144, Point{18, 24}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p06_02":        MapSymData{6, 6, 0, "Trans", 1, 72, 144, Point{12, -24}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p06_03":        MapSymData{6, 6, 0, "Trans", 1, 108, 108, Point{18, 18}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p06_04":        MapSymData{6, 6, 0, "Trans", 1, 144, 144, Point{24, 24}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p07_01":        MapSymData{7, 7, 0, "Trans", 1, 70, 140, Point{10, -40}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p07_02":        MapSymData{7, 7, 0, "Trans", 1, 84, 84, Point{12, 12}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p07_03":        MapSymData{7, 7, 0, "Trans", 1, 140, 140, Point{40, -20}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p07_04":        MapSymData{7, 7, 0, "Trans", 1, 70, 140, Point{20, -20}, -1, -1, -1, -1, Point{0, 0}},
	"maze_p08_01":        MapSymData{8, 8, 0, "Trans", 1, 128, 128, Point{16, 16}, -1, -1, -1, -1, Point{0, 0}},
	"mmaze_03p_01":       MapSymData{3, 3, 0, "Trans", 1, 108, 144, Point{36, -48}, -1, -1, -1, -1, Point{0, 0}},
	"mmaze_04p_01":       MapSymData{4, 4, 0, "Trans", 1, 64, 96, Point{16, 24}, -1, -1, -1, -1, Point{0, 0}},
	"mmaze_05p_01":       MapSymData{5, 5, 0, "Trans", 1, 150, 150, Point{30, -30}, -1, -1, -1, -1, Point{0, 0}},
	"mmaze_07p_01":       MapSymData{7, 7, 0, "Trans", 1, 84, 84, Point{24, 12}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_02p_01": MapSymData{2, 2, 0, "Trans", 1, 64, 64, Point{32, 32}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_02p_02": MapSymData{2, 2, 0, "Trans", 1, 92, 76, Point{46, 38}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_03p_01": MapSymData{3, 3, 0, "Trans", 1, 66, 129, Point{22, 43}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_03p_02": MapSymData{3, 3, 0, "Trans", 1, 126, 114, Point{42, 38}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_04p_01": MapSymData{4, 4, 0, "Trans", 1, 60, 116, Point{15, 29}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_04p_02": MapSymData{4, 4, 0, "Trans", 1, 128, 112, Point{32, 28}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_05p_01": MapSymData{5, 5, 0, "Trans", 1, 110, 140, Point{44, -28}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_05p_02": MapSymData{5, 5, 0, "Trans", 1, 115, 135, Point{23, 27}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_06p_01": MapSymData{6, 6, 0, "Trans", 1, 138, 132, Point{23, -22}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_06p_02": MapSymData{6, 6, 0, "Trans", 1, 144, 126, Point{24, 21}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_07p_01": MapSymData{7, 7, 0, "Trans", 1, 119, 126, Point{34, 18}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_07p_02": MapSymData{7, 7, 0, "Trans", 1, 119, 147, Point{34, -21}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_08p_01": MapSymData{8, 8, 0, "Trans", 1, 88, 88, Point{11, 33}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_08p_02": MapSymData{8, 8, 0, "Trans", 1, 128, 136, Point{48, -17}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_09p_01": MapSymData{9, 9, 0, "Trans", 1, 126, 126, Point{14, 28}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_09p_02": MapSymData{9, 9, 0, "Trans", 1, 135, 117, Point{15, -26}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_10p_01": MapSymData{10, 10, 0, "Trans", 1, 100, 110, Point{30, 11}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_10p_02": MapSymData{10, 10, 0, "Trans", 1, 150, 140, Point{15, -42}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_01": MapSymData{2, 2, 0, "Trans", 1, 100, 80, Point{50, 40}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_02": MapSymData{2, 2, 0, "Trans", 1, 80, 72, Point{40, 36}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_03": MapSymData{2, 2, 0, "Trans", 1, 60, 54, Point{30, 27}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_04": MapSymData{2, 2, 0, "Trans", 1, 60, 96, Point{30, 48}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_05": MapSymData{2, 2, 0, "Trans", 1, 72, 86, Point{36, 43}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_06": MapSymData{2, 2, 0, "Trans", 1, 88, 74, Point{44, 37}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_07": MapSymData{2, 2, 0, "Trans", 1, 64, 56, Point{32, 28}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_08": MapSymData{2, 2, 0, "Trans", 1, 52, 82, Point{26, 41}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_09": MapSymData{2, 2, 0, "Trans", 1, 86, 96, Point{43, 48}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_10": MapSymData{2, 2, 0, "Trans", 1, 76, 72, Point{38, 36}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_11": MapSymData{2, 2, 0, "Trans", 1, 94, 78, Point{47, 39}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_12": MapSymData{2, 2, 0, "Trans", 1, 62, 64, Point{31, 32}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_13": MapSymData{2, 2, 0, "Trans", 1, 54, 62, Point{27, 31}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_14": MapSymData{2, 2, 0, "Trans", 1, 88, 80, Point{44, 40}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_15": MapSymData{2, 2, 0, "Trans", 1, 52, 70, Point{26, 35}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_16": MapSymData{2, 2, 0, "Trans", 1, 82, 68, Point{41, 34}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_17": MapSymData{2, 2, 0, "Trans", 1, 100, 94, Point{50, 47}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_18": MapSymData{2, 2, 0, "Trans", 1, 58, 92, Point{29, 46}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_19": MapSymData{2, 2, 0, "Trans", 1, 82, 90, Point{41, 45}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_20": MapSymData{2, 2, 0, "Trans", 1, 70, 64, Point{35, 32}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_21": MapSymData{2, 2, 0, "Trans", 1, 90, 90, Point{45, 45}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_22": MapSymData{2, 2, 0, "Trans", 1, 78, 64, Point{39, 32}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_23": MapSymData{2, 2, 0, "Trans", 1, 76, 74, Point{38, 37}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_24": MapSymData{2, 2, 0, "Trans", 1, 88, 94, Point{44, 47}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_25": MapSymData{2, 2, 0, "Trans", 1, 64, 74, Point{32, 37}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_26": MapSymData{2, 2, 0, "Trans", 1, 84, 70, Point{42, 35}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_27": MapSymData{2, 2, 0, "Trans", 1, 58, 52, Point{29, 26}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_28": MapSymData{2, 2, 0, "Trans", 1, 68, 72, Point{34, 36}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_29": MapSymData{2, 2, 0, "Trans", 1, 54, 62, Point{27, 31}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_30": MapSymData{2, 2, 0, "Trans", 1, 100, 86, Point{50, 43}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_31": MapSymData{2, 2, 0, "Trans", 1, 62, 52, Point{31, 26}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_32": MapSymData{2, 2, 0, "Trans", 1, 50, 50, Point{25, 25}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_33": MapSymData{2, 2, 0, "Trans", 1, 82, 82, Point{41, 41}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_34": MapSymData{2, 2, 0, "Trans", 1, 80, 76, Point{40, 38}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_35": MapSymData{2, 2, 0, "Trans", 1, 94, 96, Point{47, 48}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_36": MapSymData{2, 2, 0, "Trans", 1, 88, 76, Point{44, 38}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_37": MapSymData{2, 2, 0, "Trans", 1, 52, 68, Point{26, 34}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_38": MapSymData{2, 2, 0, "Trans", 1, 98, 98, Point{49, 49}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_39": MapSymData{2, 2, 0, "Trans", 1, 88, 76, Point{44, 38}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_40": MapSymData{2, 2, 0, "Trans", 1, 86, 100, Point{43, 50}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_41": MapSymData{2, 2, 0, "Trans", 1, 90, 82, Point{45, 41}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_42": MapSymData{2, 2, 0, "Trans", 1, 100, 86, Point{50, 43}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_43": MapSymData{2, 2, 0, "Trans", 1, 78, 76, Point{39, 38}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_44": MapSymData{2, 2, 0, "Trans", 1, 94, 96, Point{47, 48}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_45": MapSymData{2, 2, 0, "Trans", 1, 60, 68, Point{30, 34}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p02_46": MapSymData{2, 2, 0, "Trans", 1, 74, 92, Point{37, 46}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_01": MapSymData{3, 3, 0, "Trans", 1, 78, 96, Point{26, 32}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_02": MapSymData{3, 3, 0, "Trans", 1, 129, 117, Point{43, -39}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_03": MapSymData{3, 3, 0, "Trans", 1, 111, 123, Point{37, 41}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_04": MapSymData{3, 3, 0, "Trans", 1, 120, 114, Point{40, 38}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_05": MapSymData{3, 3, 0, "Trans", 1, 120, 144, Point{40, 48}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_06": MapSymData{3, 3, 0, "Trans", 1, 144, 120, Point{48, -40}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_07": MapSymData{3, 3, 0, "Trans", 1, 141, 135, Point{47, -45}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_08": MapSymData{3, 3, 0, "Trans", 1, 135, 135, Point{45, 45}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_09": MapSymData{3, 3, 0, "Trans", 1, 90, 114, Point{30, 38}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_10": MapSymData{3, 3, 0, "Trans", 1, 78, 72, Point{26, 24}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_11": MapSymData{3, 3, 0, "Trans", 1, 96, 81, Point{32, 27}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_12": MapSymData{3, 3, 0, "Trans", 1, 132, 132, Point{44, -44}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_13": MapSymData{3, 3, 0, "Trans", 1, 117, 132, Point{39, -44}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_14": MapSymData{3, 3, 0, "Trans", 1, 123, 150, Point{41, 50}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_15": MapSymData{3, 3, 0, "Trans", 1, 54, 105, Point{18, 35}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_16": MapSymData{3, 3, 0, "Trans", 1, 129, 132, Point{43, -44}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_17": MapSymData{3, 3, 0, "Trans", 1, 132, 111, Point{44, 37}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p03_18": MapSymData{3, 3, 0, "Trans", 1, 78, 150, Point{26, -50}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p04_01": MapSymData{4, 4, 0, "Trans", 1, 148, 132, Point{37, 33}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p04_02": MapSymData{4, 4, 0, "Trans", 1, 124, 120, Point{31, 30}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p04_03": MapSymData{4, 4, 0, "Trans", 1, 120, 144, Point{30, 36}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p04_04": MapSymData{4, 4, 0, "Trans", 1, 88, 148, Point{22, 37}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p04_05": MapSymData{4, 4, 0, "Trans", 1, 100, 112, Point{25, -28}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p04_06": MapSymData{4, 4, 0, "Trans", 1, 80, 144, Point{20, -36}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p04_07": MapSymData{4, 4, 0, "Trans", 1, 116, 148, Point{29, -37}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p04_08": MapSymData{4, 4, 0, "Trans", 1, 100, 128, Point{25, 32}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p05_01": MapSymData{5, 5, 0, "Trans", 1, 85, 150, Point{17, -30}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p05_02": MapSymData{5, 5, 0, "Trans", 1, 75, 125, Point{30, 25}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p05_03": MapSymData{5, 5, 0, "Trans", 1, 150, 120, Point{30, -48}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p05_04": MapSymData{5, 5, 0, "Trans", 1, 80, 90, Point{32, -18}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p05_05": MapSymData{5, 5, 0, "Trans", 1, 125, 125, Point{25, 25}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p06_01": MapSymData{6, 6, 0, "Trans", 1, 132, 126, Point{22, 21}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p06_02": MapSymData{6, 6, 0, "Trans", 1, 150, 126, Point{25, 21}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p06_03": MapSymData{6, 6, 0, "Trans", 1, 138, 132, Point{23, -22}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p06_04": MapSymData{6, 6, 0, "Trans", 1, 150, 126, Point{25, 21}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p06_05": MapSymData{6, 6, 0, "Trans", 1, 144, 120, Point{24, -20}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p07_01": MapSymData{7, 7, 0, "Trans", 1, 105, 119, Point{30, -17}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p07_02": MapSymData{7, 7, 0, "Trans", 1, 119, 105, Point{34, -15}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p08_01": MapSymData{8, 8, 0, "Trans", 1, 80, 104, Point{30, 13}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p08_02": MapSymData{8, 8, 0, "Trans", 1, 144, 120, Point{18, -45}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p10_01": MapSymData{10, 10, 0, "Trans", 1, 140, 130, Point{14, -39}, -1, -1, -1, -1, Point{0, 0}},
	"random_walk_p10_02": MapSymData{10, 10, 0, "Trans", 1, 140, 140, Point{14, 42}, -1, -1, -1, -1, Point{0, 0}},
	"test":               MapSymData{5, 5, 0, "Trans", 1, 100, 100, Point{20, -40}, -1, -1, -1, -1, Point{0, 0}},
}
