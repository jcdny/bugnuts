include $(GOROOT)/src/Make.inc

TARG=bot
GOFILES=\
	main.go

include $(GOROOT)/src/Make.cmd

.PHONY: gofmt test dist packages $(TARG)
gofmt:
	gofmt -w $(GOFILES)

dist:
	zip dist/dist.$(shell date +%Y%m%d-%H%M).zip $(GOFILES)

test: $(TARG)
	./bot -b v7 -V none -d 1 < i.tmp 2>&1 > /dev/null

testall: $(TARG)
	./bot -b v7 -V none -d 0 < i.tmp 2>&1 > /dev/null
	./bot -b v6 -V none -d 0 < i.tmp 2>&1 > /dev/null
	./bot -b v5 -V none -d 0 < i.tmp 2>&1 > /dev/null
	./bot -b v3 -V none -d 0 < i.tmp 2>&1 > /dev/null
	./bot -b sb -V none -d 0 < i.tmp 2>&1 > /dev/null

#	./bot -V none -d=1 -w 195:@88,36 < i.tmp 2>&1 > /dev/null

$(GOFILES): packages

packages:
	+@echo Making packages
	+@$(MAKE) -C ../../pkg install || (echo FAILED; exit 1)
