package main

import (
	ruleset "./httpse-lib"
	"bytes"
	"fmt"
	"regexp"
	"log"
	"io/ioutil"
	"os"
	"encoding/xml"
	"path/filepath"
	"strings"
)

// generic rule regex
var grr = "<rule\\s+from=\"[^\"]+\"[\\s\r\n]+to=\"[^\"]+\"\\s+/>"

// generic target regex
var gtr = "<target\\s+host=\"[^\"]+\"\\s+/>"

// trivial rule
var tr = "<rule from=\"^http:\" to=\"https:\" />"

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s PATH/https-everywhere PATH/trivialize-whitelist.txt\n", os.Args[0])
		os.Exit(1)
	}

	// read file from PATH/https-everywhere/rules
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// num_changes
	num_changes := 0


	// iterate through all '*.xml'
	for _, file := range files {
		fn := file.Name()
		fp := filepath.Join(os.Args[1], fn)

		// ignore file without '.xml' extension
		if strings.HasSuffix(fn, ".xml") == false {
			log.Printf("Skipping %s (.xml)", fn)
			continue
		}

		xmlContent, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Print(err)
			continue
		}

		// parse xml ruleset
		var r ruleset.Ruleset
		xml.Unmarshal(xmlContent, &r)

		// ignore default_off because they are likely problematic
		if len(r.Default_off) > 0 {
			continue
		}

		// rewrite trivial rule only as intended
		if len(r.Rules) != 1 || len(r.Targets) != 1 {
			continue
		}


		// dummy
		pxml := xmlContent

		// an array of func which rewrite trivial rulesets
		trivialize_type_n := [2]func([]byte, ruleset.Ruleset) []byte {
			trivialize_type_1,
			trivialize_type_2,
		}

		for _, f := range trivialize_type_n {
			// ruleset got rewritten
			if pxml = f(xmlContent, r); bytes.Compare(pxml, xmlContent) != 0 {
				// TODO print something here
				err := ioutil.WriteFile(fp, []byte(pxml), 0644)
				if err != nil {
					log.Print(err)
				} else {
					num_changes++
				}
			}
		}
	}

	log.Printf("Rewritten %d files", num_changes)
}

func trivialize_type_1(xmlContent []byte, r ruleset.Ruleset) []byte {
	target := r.Targets[0].Host
	from := r.Rules[0].From
	to   := r.Rules[0].To

	// no wildcard
	if strings.Contains(target, "*") {
		return xmlContent
	}

	tfrom := "^http://" + regexp.QuoteMeta(target) + "/"
	tto   := "https://" + target + "/"

	if from == tfrom && to == tto {
		re := regexp.MustCompile(grr)
		return re.ReplaceAll(xmlContent, []byte(tr))
	}
	return xmlContent
}

func trivialize_type_2(xmlContent []byte, r ruleset.Ruleset) []byte {
	target := r.Targets[0].Host
	from := r.Rules[0].From
	to   := r.Rules[0].To

	// only wildcard
	if strings.HasSuffix(target, "*.") {
		return xmlContent
	}

	// remove '*.'
	target = target[2:len(target)]

	// rule#from regex prefix
	rfrp := "^\\^http://\\((\\?:)?(\\w+\\|?)+\\)\\\\."
	rfrm := strings.Replace(regexp.QuoteMeta(target), "\\", "\\\\", -1)
	rfrs := "/$"

	if to != "https://$1." + target + "/" {
		return xmlContent
	}

	fr := regexp.MustCompile(rfrp + rfrm + rfrs)
	if fr.MatchString(from) == false {
		return xmlContent
	}

	start := strings.Index(from, "(") + 1
	end   := strings.Index(from, ")")

	t := ""

	// rewrite 'rule'
	rr1  := regexp.MustCompile(grr)
	pxml := rr1.ReplaceAll(xmlContent, []byte(tr))

	// rewrite 'target'
	foo := strings.Split(strings.Replace(from[start:end], "?:", "", -1), "|")
	for _, bar := range foo {
		tmp := "<target host=\"" + bar + "." + target + "\" />\n\t"
		t += tmp
	}

	// remove tailing '\n\t'
	t = t[0:len(t) - 2]

	// rewrite 'target'
	rr2 := regexp.MustCompile(gtr)
	pxml = rr2.ReplaceAll(pxml, []byte(t))

	return pxml
}
