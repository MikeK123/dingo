package main

import (
	"flag"
	"log"
	"runtime"
	"runtime/debug"

	"github.com/MikeK123/dingo/explorer"
	"github.com/MikeK123/dingo/generators"
	"github.com/MikeK123/dingo/model"
	"github.com/MikeK123/dingo/producers"
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
	log.Printf("Processing configuration file %s\r\n", configPath)
	config := model.LoadConfiguration(configPath)
	exp := createExplorer(&config)
	schema := exp.ExploreSchema(&config)
	modelpkg := producers.ProduceModelPackage(&config, schema)
	modelmixedpkg := producers.ProduceMixedModelPackage(&config)
	daopkg := producers.ProduceDaoPackage(&config, schema, modelpkg)
	daomixedpkg := producers.ProduceDaoMixedPackage(&config, schema, modelpkg)
	viewmodelpkg := producers.ProduceViewModelPackage(&config, schema)
	viewmodelmixedpkg := producers.ProduceMixedViewModelPackage(&config)
	bizpkg := producers.ProduceBizPackage(&config, modelpkg, daopkg, viewmodelpkg)
	bizmixedpkg := producers.ProduceBizMixedPackage(&config, modelmixedpkg, daomixedpkg, viewmodelmixedpkg)
	srvpkg := producers.ProduceServicePackage(&config, viewmodelpkg, bizpkg)
	generators.GenerateModel(&config, modelpkg)
	generators.GenerateMixedModel(&config, modelmixedpkg)
	if !config.SkipDaoGeneration {
		generators.GenerateDao(&config, daopkg)
		generators.GenerateMixedDao(&config, daomixedpkg)
		if !config.SkipBizGeneration {
			generators.GenerateViewModel(&config, viewmodelpkg)
			generators.GenerateMixedViewModel(&config, viewmodelmixedpkg)
			generators.GenerateBiz(&config, bizpkg)
			generators.GenerateMixedBiz(&config, bizmixedpkg)
			if !config.SkipServiceGeneration {
				generators.GenerateService(&config, srvpkg)
				generators.GenerateConfig(&config)
				generators.GenerateMain(&config, srvpkg)
				generators.GenerateCustomResources(&config)
			}
		}
	}
	log.Printf("Code generation done.\r\n")
}

func createExplorer(conf *model.Configuration) explorer.DatabaseExplorer {
	switch conf.DatabaseType {
	case "mysql":
		return explorer.NewMySqlExplorer()
	case "postgres":
		return explorer.NewPostgreSqlExplorer()
	default:
		return explorer.NewMySqlExplorer()
	}
}
