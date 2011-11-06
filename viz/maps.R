maps <- c("maze_02p_01",
	"maze_02p_02",
	"maze_03p_01",
	"maze_04p_01",
	"maze_04p_02",
	"maze_05p_01",
	"maze_06p_01",
	"maze_07p_01",
	"maze_08p_01",
	"mmaze_02p_01",
	"mmaze_02p_02",
	"mmaze_03p_01",
	"mmaze_04p_01",
	"mmaze_04p_02",
	"mmaze_05p_01",
	"mmaze_07p_01",
	"mmaze_08p_01",
	"random_walk_02p_01",
	"random_walk_02p_02",
	"random_walk_03p_01",
	"random_walk_03p_02",
	"random_walk_04p_01",
	"random_walk_04p_02",
	"random_walk_05p_01",
	"random_walk_05p_02",
	"random_walk_06p_01",
	"random_walk_06p_02",
	"random_walk_07p_01",
	"random_walk_07p_02",
	"random_walk_08p_01",
	"random_walk_08p_02",
	"random_walk_09p_01",
	"random_walk_09p_02",
	"random_walk_10p_01",
	"random_walk_10p_02")



pdf(file="heat.pdf")
for (map in maps) {
  x <- as.matrix(read.csv(paste(map, ".csv", sep=""), header=F))
  filled.contour(log(x+1), color.palette=heat.colors, axes=F, plot.title=title(main=paste(map, " All points in to nearest hill")))
}
dev.off()
