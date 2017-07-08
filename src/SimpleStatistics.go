package main

import (
	ruleset "./httpse-lib"
	publicsuffix "golang.org/x/net/publicsuffix"

	"path/filepath"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"log"
)

func main() {
	// https://github.com/EFForg/https-everywhere/issues/10378
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <path/https-everywhere/rules\n")
		os.Exit(1)
	}

	files, errdirio := ioutil.ReadDir(os.Args[1])
	if errdirio != nil {
		log.Fatal(errdirio)
	}

	// map of 'default_off'
	dmap := make(map[string]int)

	// map of 'platform'
	pmap := make(map[string]int)

	// map of 'rule'
	rmap := make(map[string]int)

	// map of 'target'
	tmap := make(map[string]int)

	for _, file := range files {
		// full path of file
		filename := filepath.Join(os.Args[1], file.Name())

		// extension must be ".xml"
		if strings.HasSuffix(filename, ".xml") == false {
			log.Printf("skipping %s", filename)
			continue
		}

		xmlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Print(err)
			continue
		}

		// parse xml
		var r ruleset.Ruleset
		xml.Unmarshal(xmlFile, &r)

		// count 'default_off'
		if len(r.Default_off) > 0 {
			if _, ok := dmap[r.Default_off]; ok {
				dmap[r.Default_off]++
			} else {
				dmap[r.Default_off] = 1
			}
		}

		// count 'platform'
		if len(r.Platform) > 0 {
			if _, ok := pmap[r.Platform]; ok {
				pmap[r.Platform]++
			} else {
				pmap[r.Platform] = 1
			}
		}

		// count 'rule'
		if len(r.Rules) > 0 {
			for _, rule := range r.Rules {
				key := rule.From + rule.To
				if _, ok := rmap[key]; ok {
					rmap[key]++
				} else {
					rmap[key] = 1
				}
			}
		}

		// count 'target'
		if len(r.Targets) > 0 {
			for _, target := range r.Targets {
				d, err := publicsuffix.EffectiveTLDPlusOne(target.Host)
				if err != nil {
					log.Println(err)
					break
				}

				if _, ok := tmap[d]; ok {
					tmap[d]++
				} else {
					tmap[d] = 1
				}
			}
		}
	}

	fmt.Printf("| %d, %d ", MyMapSum(pmap), len(pmap))
	fmt.Printf("| %d, %d ", MyMapSum(dmap), len(dmap))
	fmt.Printf("| %d, %d ", MyMapSum(tmap), len(tmap))
	fmt.Printf("| %d, %d ", MyMapSum(rmap), len(rmap))
	fmt.Printf("|\n")
	fmt.Printf("| %d | %d | %d ", pmap["mixedcontent"], pmap["cacert"], pmap["cacert mixedcontent"])
	fmt.Printf("|\n")
}

func MyMapSum(pmap map[string]int) int {
	ret := 0
	for _, val := range pmap {
		ret += val
	}
	return ret
}
