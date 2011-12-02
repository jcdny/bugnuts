include $(GOROOT)/src/Make.inc

TARG=mapper
GOFILES=\
	job.go\
	mapper.go\
	main.go

include $(GOROOT)/src/Make.cmd

test: $(TARG)
	./mapper testdata
