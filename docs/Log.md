### 2011-10-21 FRIDAY

## The AI Challenge

I decided to compete in the AI Challenge partly since it sounded like fun and partly out of an interest in learning Go.  I think the structure of the contest favors high performance languages since the cpu budget at high levels will be pretty lean.  Go is probably ok from that standpoint though possibly at a slight disadvantage to C/C++.  I plan to use the blog mostly to track ideas and evolution of my bot and if I do well people might be interested in seeing how I went about building things.

### Here is my initial sketch of ideas and priorities:

* Gradient search on encountered enemies.
* Repulsion on ants to increase spread.
* Threat on hill possibly weighted by precomputed bfs dist and upweight coordinated attacks
* greedy step by unk + visible land + overall visibility.
* avoidance scaling - avoid early while few ants gathering and scouting.
* count of ants vs food gotten so can decide whether to plug holes on multi hole version based on threat.
* array of how much visible/reachable dirt at each point out to 5 steps.
* goal chaining - ie if I am next to an ant and I want to move into its square to achieve my goal I pass my goal to it.
* clump rounds to make my ants less subject to the line of eating.
* pathfind out of hill to a rally point.  rally point determined by safety???
* cat map type as maze or random and pathfind differently?
* death avoidance - count deads I see and avoid danger when scouting.

### Performance Ideas:

* arrays mapping whfa <-> byte
* arrays mapping byte -> threat
* arrays mapping byte -> value

### Initial Priorities:

* gather food
* max explore + visible dirt
* avoid enemies

## 2011-10-22

I decided to mostly toss the starter bot and rewrite it (with more
tests and benchmarking and hopefully more flexibility around plugging
strategies).  Motivated in part since when I tried to write a test of
the turn parser I was finding it annoying.  Also the code seemed to do
a huge amount of memory churn though that might not matter; I don't
have a good feel for Go's garbage collection.

I also had the idea that it would be fun to analyse the replay data do
see if it's possible to identify strategies.  I will start pulling
that data down to play with it.

2011-10-25

I got a version up with fewer bugs (heh) and less prone to getting
stuck but of course no pathfinding and no combat stuff means really
poor performance.

Working on pathfinding now and adding tests/code shuffling.  I am
loving Go!

2011-10-27

I spent a couple days trying to write a fill that was faster than the
lame knucklehead simple queue implementation only to find my fancy
space saving version was actually slower and probably less flexible.
I might have to revisit it though since I think saving the queue at
the first unknown and being able to regen from there could be a big
win especially in late game where there is only a small part of the
map thats unkown and queue stashing is kind of expensive in the simple
version.

