package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"me-ai/pkg/db"
	"os"
	"strings"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("DB init error: %v", err)
	}

	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Cannot read migrations dir: %v", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		path := "migrations/" + file.Name()
		fmt.Printf("Applying migration: %s\n", path)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Cannot read migration %s: %v", path, err)
		}
		_, err = db.DB.Exec(string(content))
		if err != nil {
			log.Fatalf("Migration %s failed: %v", path, err)
		}
	}
	fmt.Println("All migrations applied successfully!")
}
