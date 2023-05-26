package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Kseleven/psql2csv/pkg"
)

var (
	DBPwd      string
	ConfigFile string
)

func main() {
	flag.StringVar(&ConfigFile, "f", "", "config file")
	flag.Parse()

	fmt.Println("enter password for db:")
	_, err := fmt.Scanln(&DBPwd)
	catchError(err)
	if DBPwd == "" {
		catchError(fmt.Errorf("password is empty"))
	}
	if ConfigFile == "" {
		catchError(fmt.Errorf("config path is empty"))
	}

	conf, err := pkg.LoadConfig(ConfigFile)
	catchError(err)
	catchError(conf.Valid())

	conf.DBPassword = DBPwd
	db, err := pkg.NewDB(conf)
	catchError(err)

	if conf.ImportAction() {
		catchError(pkg.ImportCsv2DB(db, conf))
	} else {
		catchError(pkg.Export2Csv(db, conf))
	}

	fmt.Println("done")
}

func catchError(err error) {
	if err == nil {
		return
	}
	fmt.Println(err.Error())
	os.Exit(1)
}
