package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var location string
	var ignore string
	var err error
	var outfile string
	var method string
	var exact bool
	flag.StringVar(&location,"loc", "","Scan location")
	flag.StringVar(&ignore,"ignore","", "Ignore locations separated by ;(semicolon)")
	flag.StringVar(&method,"method", "init","Method to search")
	flag.StringVar(&outfile, "outfile", "init_methods.log","Where to store method contents")
	flag.BoolVar(&exact,"exact",false,"Exact match")

	flag.Parse()
	f, err := os.Create(outfile)
	if err!=nil {
		fmt.Printf("Failed to create target file")
		os.Exit(2)
	}
	defer f.Close()
	scanner := &Scanner{Count:0,File:f,Method:method,Exact:exact}
	if location == "" {
		location, err = os.Getwd()
		if err!=nil {
			fmt.Printf("Error occurred:%v\n",err)
			os.Exit(1)
		}
	}
	scanner.SearchForFunction(location,ignore)

	fmt.Printf("Found init function in %v files\n",scanner.Count)
}
