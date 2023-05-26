package main

import (
	"flag"
	"fmt"
	csv "github.com/Kseleven/psql2csv"
	"os"
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

	conf, err := csv.LoadConfig(ConfigFile)
	catchError(err)
	catchError(conf.Valid())

	conf.DBPassword = DBPwd
	db, err := csv.NewDB(conf)
	catchError(err)

	if conf.ImportAction() {
		catchError(csv.ImportCsv2DB(db, conf))
	} else {
		catchError(csv.Export2Csv(db, conf))
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
