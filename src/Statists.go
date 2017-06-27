package main

import (
	ruleset "./httpse-lib"

	"encoding/xml"
	"fmt"
	"sort"
	"io/ioutil"
	"log"
	"os"
)


type Pair struct {
	Key string
	Val int64
}

type PairList []Pair

func mapToSortedPairList(mmap map[string]int64) PairList {
	pl := make(PairList, len(mmap))

	i := 0
	for k, v := range mmap {
		pl[i] = Pair{k, v}
		i++
	}

	sort.Sort(sort.Reverse(pl))
	return pl
}

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	if p[i].Val == p[j].Val {
		return p[i].Key < p[j].Key
	}
	return p[i].Val < p[j].Val
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}


func main() {
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// total number of rulesets
	total := len(files)

	// default_off count
	cdo := 0

	// platform count
	cpf := 0

	// default_off 
	dmap := make(map[string]int64)

	// platform
	pmap := make(map[string]int64)

	for _, file := range files {
		xmlFile, err := ioutil.ReadFile(os.Args[1] + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
			break
		}

		var r ruleset.Ruleset
		xml.Unmarshal(xmlFile, &r)

		if len(r.Default_off) > 0 {
			if val, exist := dmap[r.Default_off]; exist {
				dmap[r.Default_off] = val + 1
			} else {
				dmap[r.Default_off] = 1
			}
			cdo++
		}


		if len(r.Platform) > 0 {
			if val, exist := pmap[r.Platform]; exist {
				pmap[r.Platform] = val + 1
			} else {
				pmap[r.Platform] = 1
			}
			cpf++
		}
	}

	xdmap := mapToSortedPairList(dmap)
	xpmap := mapToSortedPairList(pmap)

	// version, commit, total, (default_off, #), (platform, #)
	fmt.Printf("| %d | %d, %d | %d, %d |\n", total, cpf, len(xpmap), cdo, len(xdmap))

	for ind, val := range xdmap {
		if ind > 5 {
			break
		}
		fmt.Printf("| %s, %d ", val.Key, val.Val)
	}
	fmt.Printf("|\n")

	for ind, val := range xpmap {
		if ind > 5 {
			break
		}
		fmt.Printf("| %s, %d ", val.Key, val.Val)
	}
	fmt.Printf("|\n")
}
