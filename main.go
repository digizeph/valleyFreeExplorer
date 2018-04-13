package main

import (
	"fmt"
	"flag"
	"strconv"
	"os"
	"log"
)

func main() {

	flagFile := flag.String("data", "20161102.as-rel.txt.bz2", "CAIDA AS-relationship data file")

	flag.Parse()

	flagAses := flag.Args()

	var announceAses []int
	for _, asnStr := range flagAses{
		asn,err := strconv.Atoi(asnStr)
		check(err)
		announceAses = append(announceAses, asn)
		log.Printf("flag: %d", asn)
	}

	if len(announceAses) == 0 {
		log.Fatal("USAGE: valleyFreeExplorer [-data DATAFILE] ASN1 ASN2 ASN3 ...")
		os.Exit(1)
	}


	ases := BuildTopology(*flagFile)
	count := 0

	for _, asn := range announceAses {
		paths := DepthFirstSearch(ases[asn], GoUp, Path{}, map[int]bool{})
		fileName := fmt.Sprintf("%d-paths.txt.gz", asn)
		OutputPathsToGzip(fileName, paths)
		fmt.Printf("%d\tAS%d: %d paths in total\n", count, asn, len(paths))
		count++
	}
}

