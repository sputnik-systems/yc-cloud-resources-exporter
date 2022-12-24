package main

import (
	"log"

	"github.com/sputnik-systems/yc-cloud-resources-exporter/internal/resources"
)

func main() {
	if err := resources.Run(); err != nil {
		log.Fatal(err)
	}
}
