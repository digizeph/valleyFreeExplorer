package main

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"os"
	"strings"
	"bytes"
	"compress/gzip"
	"strconv"
)

// The AS struct.
type As struct {
	// AS description
	Asn           int     // ASN for the AS

	// neighbor information
	Customers         map[int]*As       // list of customer ASes
	Peers             map[int]*As       // list of peer ASes
	Providers         map[int]*As       // list of provider ASes
}

// constructor
func NewAs(asn int) *As {
	as := new(As)

	as.Asn = asn
	as.Customers = map[int]*As{}
	as.Peers = map[int]*As{}
	as.Providers = map[int]*As{}

	return as
}

type Path []*As      // the path (list of pointers to ASes, destination on left)
type Prefix = string // define prefix as a string

// BuildTopology builds AS-level topology
func BuildTopology(relFilePath string) map[int]*As {
	var ll string
	asMap := map[int]*As{}

	inputFile, err := os.Open(relFilePath)
	check(err)
	defer inputFile.Close()
	scanner := bufio.NewScanner(bufio.NewReader(bzip2.NewReader(inputFile)))

	for scanner.Scan() {
		ll = scanner.Text()
		if strings.HasPrefix(ll, "#") || strings.HasPrefix(ll, " ") {
			continue
		}
		var (
			asn1  int
			asn2  int
			asrel int
		)
		fmt.Sscanf(ll, "%d|%d|%d", &asn1, &asn2, &asrel)

		// create ASes if not exists
		if _, exists := asMap[asn1]; !exists {
			asMap[asn1] = NewAs(asn1)
		}
		if _, exists := asMap[asn2]; !exists {
			asMap[asn2] = NewAs(asn2)
		}

		as1 := asMap[asn1]
		as2 := asMap[asn2]

		if asrel == 0 {
			as1.Peers[as2.Asn] = as2
			as2.Peers[as1.Asn] = as1
		} else {
			as1.Customers[asn2] = as2
			as2.Providers[as1.Asn] = as1
		}
	}

	return asMap
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// OutputPathsToGzip writes the paths content to a gzip file
func OutputPathsToGzip(fileName string, paths []Path) {
	_ = os.Mkdir("paths",os.ModePerm)
	f, err := os.Create("paths/"+fileName)
	check(err)
	defer f.Close()

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	for _, path := range paths {
		var segments []string
		for _, asptr := range path {
			segments = append(segments, strconv.Itoa(asptr.Asn))
		}
		bytes2 := []byte(strings.Join(segments, ",") + "\n")
		gz.Write(bytes2)
	}

	if err = gz.Flush(); err != nil {
		return
	}
	if err = gz.Close(); err != nil {
		return
	}

	f.Write(b.Bytes())
}
