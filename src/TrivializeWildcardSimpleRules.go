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

		// #wildcard target == 1
		if len(r.Targets) != 1 || strings.HasPrefix(r.Targets[0].Host, "*.") == false {
			continue
		}

		wtarget := r.Targets[0].Host
		target  := wtarget[2:len(wtarget)]
		from   := r.Rules[0].From
		to     := r.Rules[0].To

		re := regexp.MustCompile("^\\^http://\\(((\\w+)\\|?)+\\)\\\\." + strings.Replace(regexp.QuoteMeta(target), "\\", "\\\\", -1) + "/$")
		if re.MatchString(from) == false {
			continue
		}


		if to != "https://$1." + target + "/" {
			continue
		}



		newTarget := ""

		foo := from[strings.Index(from, "(") + 1:strings.Index(from, ")")]

		for _, bar := range strings.Split(foo, "|") {
			newTarget += ("<target host=\"" + bar + "." + target + "\" />\n\t")
		}

		newTarget = newTarget[0:len(newTarget)-2]


		replaceTargetRegex := regexp.MustCompile("<target\\s+host=\"[^\"]+\"\\s+/>")
		pxml := replaceTargetRegex.ReplaceAllString(string(xmlFile), newTarget)

		replaceRuleRegex   := regexp.MustCompile("<rule\\s+from=\"[^\"]+\"[\\s\r\n]+to=\"[^\"]+\"\\s+/>")
		pxml = replaceRuleRegex.ReplaceAllString(string(pxml), "<rule from=\"^http:\" to=\"https:\" />")

		err = ioutil.WriteFile(filename, []byte(pxml), 0644)
		if err != nil {
			log.Print(err)
		}
	}
}
