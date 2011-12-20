// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

/*

The actual bot.

Flags are:
	-T=65535
		Maximum turn to run, will override the parsed turn limit.
	-V="[flag,...]" Visualization flags
		all,none,useful,targets,vcount,monte,sym,goals,
		horizon,threat,path,error,symgen, combat,tborder,risk.
	-b="v8"	Which bot to run.
		sb: Statbot - noop bot that collects statistics
		v3: V3 - diffusion bot
		v5: V5 - goal seeker
		v6: V6 - Final Noncombat bot
		v7: V7 - combat bot
		v8: V8 - combat bot
	-d=0
		Debug level
	-m=""
		Map file -- Used to validate generated map, hill guessing etc.
	-w=""
		Watch points "T1:T2@R,C,N[;T1:T2...]", 2:4:0@10,10,2 
		watches turns 2-4 for player 0 in at location 10,10 +/2
		cells; ":" will watch everything.

*/
package documentation
