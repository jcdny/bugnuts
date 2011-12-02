-- original csv generated in the TestMapId() func in map_test.go
create temporary table maptmp (LIKE maps);
alter table maptmp drop column mid;
\copy maptmp from 'maps.csv' csv null E'\\N';
update maptmp set mapid = lower(mapid);
delete from maptmp where exists (select 1 from maps where maptmp.mapid = maps.mapid);
insert into maps select nextval('maps_mid_seq'), * from maptmp;
