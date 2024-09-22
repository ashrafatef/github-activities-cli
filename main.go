/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"githubActivitiesCli/cmd"
	"githubActivitiesCli/database"
)

func main() {
	database.InitTables()
	cmd.Execute()
	defer database.CloseDB()
}
