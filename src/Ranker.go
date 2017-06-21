package main

import (
	ruleset "./httpse-lib"

	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <top-1m.csv> <Example.com.xml>\n", os.Args[0])
		os.Exit(-1)
	}

	file, err := os.Open(os.Args[1])
	xmlFile, err := ioutil.ReadFile(os.Args[2])

	defer file.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	dmap := make(map[string]int)
	reader := csv.NewReader(file)

	for i := 1; ;{
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		dmap[record[1]] = i
		i = i + 1
	}

	var r ruleset.Ruleset
	xml.Unmarshal(xmlFile, &r)

	for _, target := range r.Targets {
		if val, ok := dmap[target.Host]; ok {
			os.Exit(val)
		}
	}
	os.Exit(0)
}
