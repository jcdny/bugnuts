-- games with high ranked players.
select gameid, x.N
from games g,
     (select gid, count(*) as N from gameplayer group by gid having max(challengerank) < 20) x
where x.gid = g.gid
order by gameid desc;

