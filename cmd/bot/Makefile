include $(GOROOT)/src/Make.inc

TARG=bot
GOFILES=\
	main.go

include $(GOROOT)/src/Make.cmd

.PHONY: gofmt test dist packages
gofmt:
	gofmt -w $(GOFILES)

dist:
	zip dist/dist.$(shell date +%Y%m%d-%H%M).zip $(GOFILES)

test: $(TARG)
	./bot -V none -d=1 < testdata/test.input

$(GOFILES): packages

packages:
	+@echo Making packages
	+@$(MAKE) -C ../../pkg install || (echo FAILED; exit 1)
