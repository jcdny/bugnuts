include $(GOROOT)/src/Make.inc

TARG=MyBot
GOFILES=\
	ants.go\
	map.go\
	main.go\
	debugging.go\
	MyBot.go\

include $(GOROOT)/src/Make.cmd

dist:
	zip tmp/dist.zip $(GOFILES)


.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

