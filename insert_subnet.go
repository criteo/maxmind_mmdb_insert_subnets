package main

import (
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"log"
	"net"
	"os"
)

func main() {

	source_mmdb_file_name := flag.String("source", "GeoIP2-City.mmdb", "source file name")
	destination_mmdb_file_name := flag.String("dest", "GeoIP2-City-new.mmdb", "destination file name")
	subnet := flag.String("subnet", "10.0.0.1/32", "subnet to be inserted")
	country := flag.String("country", "XX", "country")
	flag.Parse()

	writer, err := mmdbwriter.Load(*source_mmdb_file_name, mmdbwriter.Options{IncludeReservedNetworks: true})
	if err != nil {
		log.Fatal(err)
	}
	// Define and insert the new data.
	_, sreNet, err := net.ParseCIDR(*subnet)
	if err != nil {
		log.Fatal(err)
	}
	sreData := mmdbtype.Map{
		"country": mmdbtype.Map{
			"iso_code": mmdbtype.String(*country),
		},
	}
	if err := writer.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
		log.Fatal(err)
	}

	// Write the newly enriched DB to the filesystem.
	fh, err := os.Create(*destination_mmdb_file_name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
