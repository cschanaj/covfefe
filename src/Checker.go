package main

import (
	ruleset "./httpse-lib"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup

	for _, file := range files {
		go func(filename string) {
			wg.Add(1)
			defer wg.Done()

			shddel, err := parse_and_check(os.Args[1] + "/" + filename)
			if err != nil {
				log.Fatal(err)
			}

			if shddel {
				mutex.Lock()
				fmt.Println(filename)
				mutex.Unlock()
			}
		}(file.Name())
	}

	wg.Wait()
}

func parse_and_check(filename string) (bool, error) {
	xmlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	var r ruleset.Ruleset
	xml.Unmarshal(xmlFile, &r)

	if len(r.Default_off) > 0 && len(r.Targets) < 3 {
		// ignore ruleset with a wildcard target
		for _, target := range r.Targets {
			if strings.Contains(target.Host, "*") {
				return false, nil
			}
		}

		for _, target := range r.Targets {
			client := &http.Client{
				Timeout: 60 * time.Second,
			}

			_, err := client.Get("https://" + target.Host)

			if err == nil {
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}
