package ceresdb

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jfcarter2358/ceresdb-go/connection"
)

var collectionNames = []string{"assets"}
var collectionSchemas = []map[string]string{
	{
		"id":       "STRING",
		"name":     "STRING",
		"type":     "STRING",
		"url":      "STRING",
		"created":  "STRING",
		"updated":  "STRING",
		"tags":     "LIST",
		"metadata": "DICT",
		"models":   "LIST",
	},
}

func VerifyDatabase(databaseName string) error {
	databases, err := connection.Query("get database")
	if err != nil {
		return err
	}
	for _, db := range databases {
		if db["name"].(string) == databaseName {
			log.Printf("Database %v exists!", databaseName)
			return nil
		}
	}
	log.Printf("Database %v does not exist, creating now", databaseName)
	_, err = connection.Query(fmt.Sprintf("post database %v", databaseName))
	return err
}

func VerifyCollections(databaseName string) error {
	for idx, collectionName := range collectionNames {
		collections, err := connection.Query(fmt.Sprintf("get collection %v", databaseName))
		if err != nil {
			return err
		}
		exists := false
		for _, col := range collections {
			if col["name"].(string) == collectionName {
				log.Printf("Collection %v exists!", collectionName)
				exists = true
				continue
			}
		}
		if exists {
			continue
		}
		log.Printf("Collection %v does not exist, creating now", collectionName)
		schemaData, _ := json.Marshal(&collectionSchemas[idx])
		_, err = connection.Query(fmt.Sprintf("post collection %v.%v %v", databaseName, collectionName, string(schemaData)))
		if err != nil {
			return err
		}
	}
	return nil
}
