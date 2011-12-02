-- original csv generated in the TestMapId() func in map_test.go
create temporary table gametmp as
       select text('') as game, gameid, played, gamelength,
       text('') as challenge, workerid, text('') as server, text('') as mapid,
       matchupid, postid from games where 0 = 1;
\copy gametmp from 'games.csv' csv null E'\\N';
update gametmp set challenge = 'ants' where challenge is null;

-- create entries in maps for maps not already present
insert into maps(mapid)
       select distinct g.mapid from gametmp g
       where g.mapid is not null
         and not exists (select 1 from maps m where g.mapid = m.mapid);

-- create a server entry
insert into servers(server, challenge) 
       select distinct g.server, g.challenge from gametmp g 
       where g.server is not null
         and not exists (select 1 from servers s where s.server = g.server and s.challenge = g.challenge);

-- nuke games already inserted.
delete from gametmp g where exists (
       select 1 from servers s, games gg
       where g.server = s.server
       and g.challenge = s.challenge
       and gg.sid = s.sid
       and gg.gameid = g.gameid);
 
insert into games(sid, mid, gameid, played, gamelength, workerid, matchupid, postid)
       select s.sid, m.mid, g.gameid, g.played, g.gamelength, g.workerid, g.matchupid, g.postid
       from servers s, maps m, gametmp g
       where g.server = s.server and g.challenge = s.challenge and m.mapid = g.mapid;
