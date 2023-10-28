package main

import (
	"flag"
)

//#	   (  (          (  (          (  (
//#    )\))(  (      )\))(  (      )\))(  (
//#   ((_))\  )\    ((_))\  )\    ((_))\  )\
//#   (()(_)((_)    (()(_)((_)    (()(_)((_)
//#	  / _` |/ _ \   / _` |/ _ \   / _` |/ _ \
//#	  \__, |\___/   \__, |\___/   \__, |\___/
//#	  |___/         |___/         |___/  route

var DIR string

func init() {
	DIR = env("DIR", "#_PRIVATE_#")
	flag.Bool("load", false, "Импорт данных из api с сохранением в MongoDB")
	flag.Bool("write", false, "Выгрузить статику из базы MongoDB в metalapi")
	flag.Parse()
}

func do(command string) {
	// go routing arguments like this:
	// ./locations --load city          > load-city
	// ./locations --load --write city  > load-city && write-city
	// ./locations --load city tariff   > load-city && load-tariff
	routing := map[string]func(){
		"load-city":  loadCities,
		"load-zones": loadZones,
		"load-pvz":   loadPvz,
		"write-city": writeCities,
		"write-pvz":  writePvz,
	}
	if routing[command] != nil {
		routing[command]()
	}
}

func main() {
	// go get databases
	mysql.Init()
	mongodb.Init()
	defer mysql.Close()
	defer mongodb.Close()

	for _, command := range flag.Args() {
		if flag.Lookup("load").Value.String() == "true" {
			do("load-" + command)
		}
		if flag.Lookup("write").Value.String() == "true" {
			do("write-" + command)
		}
		do(command)
	}
}
