package generators

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/MikeK123/dingo/model"
)

const (
	bizDirectory = "/biz"
)

var (
	BizTemplate      = "./templates/biz.tpl"
	BizFile          = "biz.go"
	BizMixedTemplate = "./templates/biz_mixed.tpl"
	BizMixedFile     = "biz_mixed.go"
)

func GenerateBiz(config *model.Configuration, pkg *model.BizPackage) {
	// load template
	file, err := ioutil.ReadFile(BizTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", BizTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + bizDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + BizFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating Biz file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(pkg, tpl, w)
	w.Flush()
}

func GenerateMixedBiz(config *model.Configuration, pkg *model.BizPackage) {
	// load template
	file, err := ioutil.ReadFile(BizMixedTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", BizMixedTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + bizDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + BizMixedFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating Biz file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(pkg, tpl, w)
	w.Flush()
}
