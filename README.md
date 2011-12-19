# [Bugnuts Bot](http://aichallenge.org/profile.php?user=33)

My code for the [AI Challege 2011](http://aichallenge.org)

A Go learning experience...

The code lives at https://github.com/jcdny/bugnuts

The branch ``submission'' contains the code that was actually
submitted for the contest including the file submission.zip which is
the version uploaded to aichallenge.org.

The master branch is a cleaned up version of that code with more
comments and some of the copy/paste idiocy I inflicted on myself
removed.  Also includes final maps.

# Code Roadmap

The executable is built in _src/cmd/bot_, the packages are in _src/pkg_.
The _botN_ packages are particular versions of the Bot (and _statbot_ is a
logicless bot which supports taking in the true gamestate from the
replay engine), movement logic resides mostly in Bot<N>, _replay_ and
_engine_ are the packages for implementing replay parsing. Symmetry
analysis lives partly in maps and partly in torus.  State handles game
state updating, metrics, and statistics. watcher contains timers and
the watch point evaluation.

There is a fancy replay analyzer that lives in _src/cmd/analyze_ that I 
mostly wrote to learn how to use channels.








