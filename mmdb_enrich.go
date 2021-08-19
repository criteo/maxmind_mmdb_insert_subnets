package main

import (
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func to_mmdb(object interface{}) mmdbtype.DataType {
	switch t := object.(type) {
	case string:
		return mmdbtype.String(t) // t has type string
	case map[string]interface{}:
		record := mmdbtype.Map{}
		for k, v := range t {
			record[mmdbtype.String(k)] = to_mmdb(v)
		}
		return record
	case []interface{}:
		slice := mmdbtype.Slice{}
		for _, v := range t {
			slice = append(slice, to_mmdb(v))
		}
		return slice
	}
	return nil
}

func main() {

	source_mmdb_file_name := flag.String("source", "GeoIP2-City.mmdb", "source file name")
	destination_mmdb_file_name := flag.String("dest", "GeoIP2-City-new.mmdb", "destination file name")
	infos_file := flag.String("infos", "XX", "filename with infos")
	flag.Parse()

	writer, err := mmdbwriter.Load(*source_mmdb_file_name, mmdbwriter.Options{IncludeReservedNetworks: true})
	if err != nil {
		log.Fatal(err)
	}
	yamlFile, err := ioutil.ReadFile(*infos_file)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var infos map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &infos)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for subnet, info := range infos {
		// But if you don't know the field types, you can use type switching to determine (safe):
		// Keep in mind that, since this is a map, the order is not guaranteed.
		_, sreNet, err := net.ParseCIDR(subnet)
		if err != nil {
			log.Fatal(err)
		}
		sreData := to_mmdb(info)
		if err := writer.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
			log.Fatal(err)
		}
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
