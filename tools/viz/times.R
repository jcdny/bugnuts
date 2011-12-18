
ggplot(x, aes(turn, accumulated)) +facet_wrap(~ name, scale="free_y")+geom_line()
