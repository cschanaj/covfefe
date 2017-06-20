package main

import (
	ruleset "./httpse-lib"

	"encoding/xml"
	"fmt"
	"strings"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// total number of rulesets
	total := len(files)

	// default_off count
	cdo := 0

	cdofrt := 0

	// mixed content count
	cmc := 0

	dmap := make(map[string]int)

	for _, file := range files {
		xmlFile, err := ioutil.ReadFile(os.Args[1] + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
			break
		}

		var r ruleset.Ruleset
		xml.Unmarshal(xmlFile, &r)

		if len(r.Default_off) > 0 {
			cdo = cdo + 1

			dmap[r.Default_off] = 1

			if r.Default_off == "failed ruleset test" {
				cdofrt = cdofrt + 1
			}
		}

		if len(r.Platform) > 0 && strings.Contains(r.Platform, "mixedcontent") {
			cmc = cmc + 1
		}
	}

	fmt.Printf("| %d | %d | %d | %d | %d |\n", total, cmc, cdo, cdofrt, len(dmap))
}
