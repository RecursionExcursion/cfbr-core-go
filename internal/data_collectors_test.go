package internal

import (
	"log"
	"testing"
)

func TestDataCollectors(t *testing.T) {
	szn, err := CompileSeason(2025)
	if err != nil {
		panic(err)
	}

	log.Println(szn)
}
