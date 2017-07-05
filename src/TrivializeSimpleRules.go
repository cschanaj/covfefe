package main

import (
	ruleset "./httpse-lib"

	"path/filepath"
	"encoding/xml"
	"regexp"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"log"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <path/https-everywhere/rules\n")
		os.Exit(1)
	}

	files, errdirio := ioutil.ReadDir(os.Args[1])
	if errdirio != nil {
		log.Fatal(errdirio)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".xml") == false {
			continue
		}

		xmlFile, err := ioutil.ReadFile(filepath.Join(os.Args[1], file.Name()))
		if err != nil {
			log.Fatal(err)
			continue
		}

		if strings.Contains(string(xmlFile), ".xml") {
			continue
		}

		var r ruleset.Ruleset
		xml.Unmarshal(xmlFile, &r)

		if len(r.Default_off) == 0 && len(r.Platform) == 0 && len(r.Rules) == 1 && len(r.Targets) == 1 {
			target := r.Targets[0].Host
			from := r.Rules[0].From
			to := r.Rules[0].To


			if strings.Contains(target, "*") {
				continue
			}


			if from != "^http:" && to == "https://" + target + "/" {
				re := regexp.MustCompile("<rule from=.*?\r?\n?.*?/>")
				pxml := re.ReplaceAllString(string(xmlFile), "<rule from=\"^http:\" to=\"https:\" />")

				err := ioutil.WriteFile(filepath.Join(os.Args[1], file.Name()), []byte(pxml), 0644)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
