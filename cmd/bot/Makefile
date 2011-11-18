include $(GOROOT)/src/Make.inc

TARG=bot
GOFILES=\
	main.go


include $(GOROOT)/src/Make.cmd

.PHONY: gofmt test dist
gofmt:
	gofmt -w $(GOFILES)

dist:
	zip dist/dist.$(shell date +%Y%m%d-%H%M).zip $(GOFILES)

test: $(TARG)
	./bot -V none -d=1 < log/0.bot0.input > /dev/null
