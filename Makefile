include $(GOROOT)/src/Make.inc

TARG=MyBot
GOFILES=\
	BotV0.go\
	BotV3.go\
	BotV4.go\
	BotV5.go\
	direction.go\
	pallettes.go\
	parameters.go\
	debugging.go\
	queue.go\
	fill.go\
	mask.go\
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
