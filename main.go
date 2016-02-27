package main

import (
	"flag"
	"log"
	"runtime"
	"runtime/debug"

	"github.com/maxzerbini/dingo/explorer"
	"github.com/maxzerbini/dingo/generators"
	"github.com/maxzerbini/dingo/model"
	"github.com/maxzerbini/dingo/producers"
)

var configPath string = "./config.json"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(300)
	flag.StringVar(&configPath, "conf", "./config.json", "path of the file config.json")
}

// Start the code generator
func main() {
	flag.Parse()
	log.Printf("DinGo Code Generator\r\n")
	log.Printf("Examining configuration file %s\r\n", configPath)
	config := model.LoadConfiguration(configPath)
	schema := explorer.ExploreSchema(&config)
	modelpkg := producers.ProduceModelPackage(&config, schema)
	daopkg := producers.ProduceDaoPackage(&config, schema, modelpkg)
	generators.GenerateModel(&config, modelpkg)
	generators.GenerateDao(&config, daopkg)
	log.Printf("Code generation done.\r\n")
}