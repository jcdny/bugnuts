CREATE SEQUENCE servers_sid_seq;

CREATE TABLE servers (
       sid BIGINT DEFAULT nextval('servers_sid_seq')
         PRIMARY KEY
       , server text UNIQUE NOT NULL
       , challenge text DEFAULT 'ants'
       , url text
       , user_url text
       , game_url text
       , rank_url text
);

CREATE SEQUENCE players_pid_seq;

CREATE TABLE players (
       pid BIGINT DEFAULT nextval('players_pid_seq')
         PRIMARY KEY
       , sid integer
         REFERENCES servers ON DELETE CASCADE ON UPDATE CASCADE
       , player text NOT NULL
       , userid integer
       , lastseen timestamp
       , mu double precision
       , sigma double precision
       , skill double precision
       , CONSTRAINT players_uq UNIQUE (sid, player)
);

CREATE SEQUENCE versions_vid_seq;

CREATE TABLE versions (
       vid BIGINT DEFAULT nextval('versions_vid_seq')
         PRIMARY KEY
       , pid BIGINT
         REFERENCES players ON DELETE CASCADE ON UPDATE CASCADE
       , submissionid integer
       , versionname text
       , CONSTRAINT versioncheck CHECK (submissionid is not null or versionname is not null)
       , CONSTRAINT version_uq UNIQUE (pid, submissionid, versionname)
);

CREATE SEQUENCE maps_mid_seq;

CREATE TABLE maps (
       mid BIGINT DEFAULT nextval('maps_mid_seq')
         PRIMARY KEY
       , mapid text NOT NULL
       , mapname text
       , nrows integer
       , ncols integer
       , players integer
       , hills integer
       , mapdata text
);

CREATE SEQUENCE games_gid_seq;

CREATE TABLE games (
       gid BIGINT DEFAULT nextval('games_gid_seq')
         PRIMARY KEY
       , sid BIGINT
         REFERENCES servers ON DELETE CASCADE ON UPDATE CASCADE
       , mid BIGINT
         REFERENCES maps ON DELETE CASCADE ON UPDATE CASCADE
       , gameid integer
       , played timestamp
       , gamelength integer
       , workerid integer
       , matchupid integer
       , postid integer
       , CONSTRAINT games_uq UNIQUE (sid, gameid)
);

CREATE TABLE gameplayer (
       gid BIGINT
         REFERENCES games ON DELETE CASCADE ON UPDATE CASCADE
       , pid BIGINT
         REFERENCES players ON DELETE CASCADE ON UPDATE CASCADE
       , vid BIGINT
         REFERENCES versions ON DELETE CASCADE ON UPDATE CASCADE
       , turns integer
       , score integer
       , rank integer
       , bonus integer
       , status text
       , challengerank integer
       , challengeskill double precision
       , CONSTRAINT gameplayer_pkey PRIMARY KEY (pid, gid)
);


