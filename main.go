/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"githubActivitiesCli/cmd"
	"githubActivitiesCli/database"
	"log"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user
	database.InitTables()
	cmd.Execute()
	defer database.CloseDB()
}
