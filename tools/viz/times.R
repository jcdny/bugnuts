library(ggplot2)
x <- read.csv("/tmp/v8.csv")
ggplot(x, aes(turn, accumulated)) +facet_wrap(~ name, scale="free_y")+geom_line()
xs <- x[which(x$name == "scoring" & x$count > 0),]
ggplot(xs, aes(count, accumulated)) + geom_point()
ggplot(xs, aes(msper)) + geom_histogram()
