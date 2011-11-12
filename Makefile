include $(GOROOT)/src/Make.inc

TARG=bot
GOFILES=\
	BotV0.go\
	BotV3.go\
	BotV4.go\
	BotV5.go\
	bot6.go\
	viz.go\
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
	targets.go\
	map.go\
	movement.go\
	sym.go\
	MyBot.go\
	main.go


include $(GOROOT)/src/Make.cmd

.PHONY: gofmt test dist
gofmt:
	gofmt -w $(GOFILES)

dist:
	zip dist/dist.$(shell date +%Y%m%d-%H%M).zip $(GOFILES)

test: $(TARG)
	./bot -V none < log/0.bot0.input
