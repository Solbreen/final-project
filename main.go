package main

import (
	"os"

	"github.com/Solbreen/final-project/pkg/db"
	"github.com/Solbreen/final-project/pkg/server"
	"github.com/Solbreen/final-project/pkg/utils"
)

func main() {
	utils.LoadEnv()

	dbFile := "scheduler.db"

	if envPort := os.Getenv("TODO_DBFILE"); envPort != "" {
		dbFile = envPort
	}

	err := db.Init(dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server.Run()
}
