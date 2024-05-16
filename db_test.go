package main

import (
	"log"
	"testing"
)

// var testDB *sql.DB

func TestXxx(t *testing.T) {
	testDB, err := InitDB("../data.db")
	if err != nil {
		log.Fatal(err)
	}

	err = InsertBlob(testDB, "0x124", 123)
	if err != nil {
		t.Error(err)
	}

	height, err := GetBlobHeight(testDB, "0x123")
	if err != nil {
		t.Error(err)
	}

	t.Logf("found %d in db", height)
}
