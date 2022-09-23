package main

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	DisclosureURL  = "https://guoyk93.github.io/lastwill.key.txt"
	DisclosureTerm = time.Minute * 15

	FileBeacon = "lastwill.beacon.txt"
	FileKey    = "lastwill.key.txt"

	EventSchedule = "schedule"
)

var (
	envEventName = strings.ToLower(strings.TrimSpace(os.Getenv("EVENT_NAME")))
	envSecretKey = strings.TrimSpace(os.Getenv("SECRET_KEY"))
)

func checkDisclosure() bool {
	res, err := http.Get(DisclosureURL)
	if err != nil {
		log.Println("failed to check disclosure:", err.Error())
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

	var (
		optBeacon   bool
		optDisclose bool
	)

	flag.BoolVar(&optBeacon, "beacon", false, "update beacon.txt")
	flag.BoolVar(&optDisclose, "disclose", false, "disclose the secret key")
	flag.Parse()

	if optBeacon {

		err = os.WriteFile(FileBeacon, []byte(time.Now().Format(time.RFC3339)), 0640)

	} else if optDisclose {

		log.Println("triggered by:", envEventName)

		if envEventName == EventSchedule {
			if checkDisclosure() {
				err = errors.New("key is already disclosed, stopping the workflow")
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

		if time.Now().Sub(beacon) < DisclosureTerm {
			log.Println("deadline not reached")
			return
		}

		log.Println("disclosing key")

		if envSecretKey == "" {
			err = errors.New("missing environment variable SECRET_KEY")
			return
		}

		err = os.WriteFile(FileKey, []byte(envSecretKey), 0640)
	}

}
