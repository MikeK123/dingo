package generators

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/MikeK123/dingo/model"
)

const (
	daoDirectory = "/dao"
)

var (
	DaoTemplate      = "./templates/dao.tpl"
	DaoMixedTemplate = "./templates/dao_mixed.tpl"
	DaoFile          = "dao.go"
	DaoMixedFile     = "dao_mixed.go"
)

func GenerateDao(config *model.Configuration, pkg *model.DaoPackage) {
	if config.DatabaseType == "postgres" {
		DaoTemplate = "./templates/postgres_dao.tpl"
	}
	// load template
	file, err := ioutil.ReadFile(DaoTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", DaoTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + daoDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + DaoFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating Dao file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(pkg, tpl, w)
	w.Flush()
}

func GenerateMixedDao(config *model.Configuration, pkg *model.DaoMixedPackage) {
	if config.DatabaseType == "postgres" {
		DaoMixedTemplate = "./templates/postgres_dao_mixed.tpl"
	}
	// load template
	file, err := ioutil.ReadFile(DaoMixedTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", DaoMixedTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + daoDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + DaoMixedFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating Mixed Dao file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(pkg, tpl, w)
	w.Flush()
}
