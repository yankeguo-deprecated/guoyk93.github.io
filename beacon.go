package main

import (
	"os"
	"time"
)

func main() {
	os.WriteFile("beacon.txt", []byte(time.Now().Format(time.RFC3339)), 0640)
}
