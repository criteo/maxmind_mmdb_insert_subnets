# maxmind_mmdb_insert_subnets

Enrich mmdb file base on a json source.

Used at Criteo for CI testing purpose. 

Usage:
```./mmdb_enrich -source="GeoLite2-Country.mmdb" -dest="GeoLite2-Country-new.mmdb" -infos=test/enrich.json```
 

Base on https://github.com/maxmind/mmdb-from-go-blogpost

TODO
- Create tests