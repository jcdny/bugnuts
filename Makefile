include $(GOROOT)/src/Make.inc

TARG=MyBot
GOFILES=\
	BotV0.go\
	Botv3.go\
	Botv4.go\
	pallettes.go\
	debugging.go\
	queue.go\
	fill.go\
	util.go\
	tables.go\
	state.go\
	map.go\
	MyBot.go\
	main.go


include $(GOROOT)/src/Make.cmd

dist:
	zip tmp/dist.zip $(GOFILES)


.PHONY: gofmt test
gofmt:
	gofmt -w $(GOFILES)

test:
	./MyBot < testdata/stream1.dat
