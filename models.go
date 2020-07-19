package main

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type SCPTag struct {
	Name    string
	Entries []SCPEntry
}

type SCPEntry struct {
	Points      int
	Url         string
	ObjectClass string
	Tags        []SCPTag
}

func (entry SCPEntry) Create() error {
	session := get_db_session()
	_, err := session.Run(`
		CREATE (n:SCPEntry { points: $points, url: $url })
	`, map[string]interface{}{
		"points": entry.Points,
		"url":    entry.Url,
	})
	if err != nil {
		return err
	}
	session.Close()
	return err
}

func (entry SCPTag) Create() error {
	session := get_db_session()
	_, err := session.Run(`
		CREATE (n:SCPTag {name: $name})
	`, map[string]interface{}{
		"name": entry.Name,
	})
	if err != nil {
		return err
	}
	session.Close()
	return err
}

func (entry SCPEntry) CreateOrUpdate() error {
	session := get_db_session()
	_, err := session.Run(`
		MERGE (n:SCPEntry {url: $url})
			SET n.points = $points
			SET n.object_class = $object_class
	`, map[string]interface{}{
		"points":       entry.Points,
		"url":          entry.Url,
		"object_class": entry.ObjectClass,
	})
	if err != nil {
		return err
	}
	for _, tag := range entry.Tags {
		_, err := session.Run(`
		MATCH (scp:SCPEntry {url: $url})
		MATCH (scptag:SCPTag {name: $name})
		MERGE (scp)-[:HAS]->(scptag)
		`, map[string]interface{}{
			"url":  entry.Url,
			"name": tag.Name,
		})
		if err != nil {
			return err
		}
	}
	session.Close()
	return err
}

func (tag SCPTag) CreateOrUpdate() error {
	session := get_db_session()
	_, err := session.Run(`
		MERGE (n:SCPTag {name: $name})
	`, map[string]interface{}{
		"name": tag.Name,
	})
	if err != nil {
		return err
	}
	session.Close()
	return err
}

func get_db_session() neo4j.Session {
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "password", ""), configForNeo4j40)
	if err != nil {
		log.Fatal(err)
	}
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		log.Fatal(err)
	}
	return session
}
