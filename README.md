# maxmind_mmdb_insert_subnets

Insert subnets with country information into an existing mmdb file

( Used at Criteo for CI testing purpose )

Usage:
```./insert_subnet -source="GeoLite2-Country.mmdb" -dest="GeoLite2-Country-new.mmdb" -subnet=10.10.10.10/32 -country="SPACE"```
 
TODO:
   - Insert more than 1 subnet
   - Update not only country.


