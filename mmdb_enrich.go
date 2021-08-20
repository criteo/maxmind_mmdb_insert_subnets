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
	"reflect"
)

func toMmdb(object interface{}) mmdbtype.DataType {
	switch t := object.(type) {
	case string:
		return mmdbtype.String(t) // t has type string
	case int:
		return mmdbtype.Uint64(t)
	case float32:
		return mmdbtype.Float32(t)
	case float64:
		return mmdbtype.Float64(t)
	case map[string]interface{}:
		record := mmdbtype.Map{}
		for k, v := range t {
			record[mmdbtype.String(k)] = toMmdb(v)
		}
		return record
	case []interface{}:
		slice := mmdbtype.Slice{}
		for _, v := range t {
			slice = append(slice, toMmdb(v))
		}
		return slice
	default:
		log.Fatal("Type unknown :", reflect.TypeOf(object))
	}
	return nil
}

func main() {

	mmdbFileName := flag.String("source", "GeoIP2-City.mmdb", "source file name")
	destinationMmdbFileName := flag.String("dest", "GeoIP2-City-new.mmdb", "destination file name")
	infosFile := flag.String("infos", "XX", "filename with infos")
	flag.Parse()

	writer, err := mmdbwriter.Load(*mmdbFileName, mmdbwriter.Options{IncludeReservedNetworks: true})
	if err != nil {
		log.Fatal(err)
	}
	yamlFile, err := ioutil.ReadFile(*infosFile)
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
		sreData := toMmdb(info)
		if err := writer.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
			log.Fatal(err)
		}
	}

	// Write the newly enriched DB to the filesystem.
	fh, err := os.Create(*destinationMmdbFileName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
