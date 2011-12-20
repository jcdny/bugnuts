# [Bugnuts Bot](http://aichallenge.org/profile.php?user=33)

My code for the [AI Challege 2011](http://aichallenge.org)

A Go learning experience...

This code lives at https://github.com/jcdny/bugnuts

The branch
[__submission__](https://github.com/jcdny/bugnuts/tree/submission)
contains the code that was actually submitted for the contest
including the file submission.zip which is the version uploaded to
aichallenge.org.

The master branch is a cleaned up version of that code with more
comments and some of the copy/paste idiocy I inflicted on myself
removed.  Also includes final maps.

## Code Roadmap

The executable is built in `src/cmd/bot` (see [*main.go*][1] for
that), the packages are in `src/pkg`.  The _botN_ packages are
particular versions of the Bot, [*bot8.go*][2] is the last version
which was submitted for the contest, and [*statbot.go*][3] is a
logicless bot which supports taking in the true gamestate from the
replay engine and stats generated internally to validate my code),
movement logic resides mostly in Bot<N>, `replay` and `engine` are the
packages for implementing replay parsing. Symmetry analysis lives
partly in `maps` and partly in `torus`.  The `state`package handles
game state updating, metrics, and statistics; `combat` does all the 
combat evaluation and risk calculations; and `watcher` contains
timers and watch point evaluation; and `debug` manages debug flags.

[1]: https://github.com/jcdny/bugnuts/blob/master/src/cmd/bot/main.go
[2]: https://github.com/jcdny/bugnuts/blob/master/src/pkg/bot8/bot8.go
[3]: https://github.com/jcdny/bugnuts/blob/master/src/pkg/statbot/statbot.go

There is a fancy replay analyzer that lives in `src/cmd/analyze` that I
mostly wrote to learn how to use channels.

## Documentation and Writeups

See the [wiki](https://github.com/jcdny/bugnuts/wiki) for more gory details and pretty graphs.

### Copyright and License

Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>, All rights reserved.

The source code is licensed under GPL version 2 or later, see the
file LICENSE for details.

Documentation and images here are licensed under the Creative Commons
Attribution-NonCommercial-ShareAlike 3.0 Unported License. To view a
copy of this license, visit
http://creativecommons.org/licenses/by-nc-sa/3.0/ or send a letter to
Creative Commons, 444 Castro Street, Suite 900, Mountain View,
California, 94041, USA.






