package main

import (
	"log"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("loadConfig: ", err)
	}

	if err := initSRP(config); err != nil {
		log.Fatal("initSRP: ", err)
	}
	defer config.client.Close()

	version, err := config.client.SystemVersion()
	if err != nil {
		log.Fatal("SystemVersion: ", err)
	}
	log.Println("system.version:", version)

	about, err := config.client.SystemVersion()
	if err != nil {
		log.Fatal("SystemAbout: ", err)
	}
	log.Println("system.about:", about)

	methods, err := config.client.SystemListMethods()
	if err != nil {
		log.Fatal("SystemListMethods: ", err)
	}

	for i := range methods {
		log.Println(methods[i])
	}

	root, err := config.client.GetRoot()
	if err != nil {
		log.Fatal("GetRoot: ", err)
	}

	groups, err := config.client.FindGroups("foo", root.UniqueID)
	if err != nil {
		log.Fatal("FindGroups: ", err)
	}
	log.Println(groups)
}
