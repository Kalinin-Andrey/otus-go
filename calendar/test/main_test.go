package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const delay = 5 * time.Second

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "progress", // Замените на "pretty" для лучшего вывода
		Paths:     []string{"features"},
		Randomize: 0, // Последовательный порядок исполнения
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

