package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	DisclosureURL  = "https://guoyk93.github.io/lastwill.key.txt"
	DisclosureDays = 14

	FileBeacon = "beacon.txt"
	FileKey    = "lastwill.key.txt"
)

var (
	envEventName = strings.ToLower(strings.TrimSpace(os.Getenv("EVENT_NAME")))
	envSecretKey = strings.ToLower(strings.TrimSpace(os.Getenv("SECRET_KEY")))
)

func checkDisclosure() bool {
	res, err := http.Get(DisclosureURL)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return res.StatusCode == http.StatusOK
}

func main() {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()

	log.Println("triggered by:", envEventName)

	if envEventName == "schedule" {
		if checkDisclosure() {
			err = errors.New("lastwill.key.txt is already disclosured")
			return
		}
	}

	log.Println("reading beacon:", FileBeacon)

	var buf []byte
	if buf, err = os.ReadFile(FileBeacon); err != nil {
		return
	}

	buf = bytes.TrimSpace(buf)

	var beacon time.Time
	if beacon, err = time.Parse(time.RFC3339, string(buf)); err != nil {
		return
	}

	log.Println("beacon:", beacon.Format(time.RFC3339))

	if time.Now().Sub(beacon) < time.Hour*24*DisclosureDays {
		log.Println("looks good")
		return
	}

	log.Println("disclose secret key")

	err = os.WriteFile(FileKey, []byte(envSecretKey), 0640)
}
