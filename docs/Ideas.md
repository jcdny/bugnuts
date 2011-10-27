## An initial stab at ideas to try

* Gradient search on encountered enemies.
* Repulsion on ants to increase spread
* Threat on hill possibly weighted by precomputed bfs dist and upweight coordinated attacks
* greedy step by unk + visible land + overall visibility
* avoidance scaling - avoid early while few ants gathering and scouting.
* count of ants vs food gotten so can decide whether to plug holes on multi hole version based on threat
* array of how much visible/reachable dirt at each point out to 5 steps
* goal chaining - ie if I am next to an ant and I want to move into its square to achieve my goal I pass my goal to it.
* clump rounds to make my ants less subject to the line of eating.
* pathfind out of hill to a rally point.  rally point determined by safety???
* cat map type as maze or random and pathfind differently?
* death avoidance - count deads I see and avoid danger when scouting.

### Performance

arrays mapping wall/hill/food/ant <-> byte
arrays mapping byte -> threat
arrays mapping byte -> value




