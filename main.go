package main

import (
	"time"

	preparation "github/demo/init"
	"github/demo/repository"
	"github/demo/rest"
	"github/demo/service"
	"github/demo/utils/log"
)

func main() {
	// Init config
	cf := preparation.Init()
	cf.Watch()

	// Init repository
	e, err := repository.NewEngine(cf)
	if err != nil {
		log.Error(err)
	}

	e.Database.SetPool(10, 100, time.Hour)

	// Init service
	err = service.Init(cf, e)
	if err != nil {
		log.Error(err)
	}

	// Init rest
	router := rest.Init()
	router.Run(":8080")
}
