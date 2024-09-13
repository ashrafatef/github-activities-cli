package store

import "github.com/sopherapps/go-scdb/scdb"

var conn *scdb.Store

func InitDB() *scdb.Store {
	if conn == nil {
		var maxKeys uint64 = 1000000
		var redundantBlocks uint16 = 1
		var poolCapacity uint64 = 10
		var compactionInterval uint32 = 1_800
		conn, err := scdb.New("githubActivitiesDB", &maxKeys, &redundantBlocks, &poolCapacity, &compactionInterval, true)
		if err != nil {
			panic(err)
		}
		return conn
	}
	return conn
}
