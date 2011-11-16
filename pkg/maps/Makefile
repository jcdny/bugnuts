include $(GOROOT)/src/Make.inc

TARG=bugnuts/maps
GOFILES=\
	direction.go\
	item.go\
	map.go\
	mask.go\
	point.go\
	torus.go\
	util.go


include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

