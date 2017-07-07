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
		// full path of file
		filename := filepath.Join(os.Args[1], file.Name())

		// extension must be ".xml"
		if strings.HasSuffix(filename, ".xml") == false {
			continue
		}

		xmlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Print(err)
			continue
		}

		// make sure no linked ruleset
		if strings.Contains(string(xmlFile), ".xml") {
			continue
		}

		// parse xml
		var r ruleset.Ruleset
		xml.Unmarshal(xmlFile, &r)

		// no default_off and platform
		if len(r.Default_off) > 0 || len(r.Platform) > 0 {
			continue
		}

		// no exclusion, must rewrite everything
		if len(r.Exclusions) != 0 {
			continue
		}

		// #rule == 1
		if len(r.Rules) != 1 {
			continue
		}

		// #non-wildcard target == 1
		if len(r.Targets) != 1 || strings.Contains(r.Targets[0].Host, "*") {
			continue
		}



		target := r.Targets[0].Host
		from   := r.Rules[0].From
		to     := r.Rules[0].To


		trivial_from := "^http://" + regexp.QuoteMeta(target) + "/"
		trivial_to   := "https://" + target + "/"

		// apply rewrite to exact match only
		if from == trivial_from && to == trivial_to {
			re := regexp.MustCompile("<rule\\s+from=\"[^\"]+\"[\\s\r\n]+to=\"[^\"]+\"\\s+/>")
			pxml := re.ReplaceAllString(string(xmlFile), "<rule from=\"^http:\" to=\"https:\" />")

			err := ioutil.WriteFile(filepath.Join(os.Args[1], file.Name()), []byte(pxml), 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
