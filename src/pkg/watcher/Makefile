include $(GOROOT)/src/Make.inc

TARG=bugnuts/watcher
GOFILES=\
	debugging.go\
	timing.go\
	watcher.go

include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

