-- original csv generated in the TestMapId() func in map_test.go
create temporary table gptmp as
       select text('') as dummy, text('') as player, 0 as gameid, turns, score, 
       rank, bonus, status, userid, 0 as submissionid, challengerank, challengeskill
       from gameplayer gp, players p where 0 = 1;
\copy gptmp from 'players.csv' csv null E'\\N';
update gametmp set challenge = 'ants' where challenge is null;

-- create players that do not exist
-- TODO need to store a uuid for games
-- TODO support playername -> versionname
insert into players(sid, player, userid)
       select distinct s.sid, player, userid
       from servers s, gptmp gp
       where s.server = 'aichallenge.org' and s.challenge = 'ants'
       and not exists (select 1 from players p where p.sid = s.sid and p.player = gp.player);

-- create versions for the given player
insert into versions(pid, submissionid)
       select distinct p.pid, submissionid
       from players p, servers s, gptmp gp
       where s.server = 'aichallenge.org' and s.challenge = 'ants'
       and p.sid = s.sid and p.player = gp.player
       and not exists (select 1 from versions v where v.pid = p.pid and v.submissionid = gp.submissionid);

-- nuke records that exist already
delete from gptmp gpt where exists (
       select 1
       from servers s, players p, games g, gameplayer gp
       where s.server = 'aichallenge.org' and s.challenge = 'ants'
       and p.sid = s.sid
       and p.player = gpt.player
       and g.sid = s.sid
       and g.gameid = gpt.gameid
       and gp.pid = p.pid
       and gp.gid = g.gid);

insert into gameplayer(gid, pid, vid, turns, score, bonus, status, challengerank, challengeskill)
       select g.gid, p.pid, v.vid, turns, score, bonus, status, challengerank, challengeskill
       from gptmp gpt, servers s, players p, games g, versions v
       where s.server = 'aichallenge.org' and s.challenge = 'ants'
       and p.sid = s.sid
       and p.player = gpt.player
       and g.sid = s.sid
       and g.gameid = gpt.gameid
       and v.pid = p.pid
       and v.submissionid = gpt.submissionid;

