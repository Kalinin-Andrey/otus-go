//package hw06
package main

import(
	//"flag"
	//"github.com/cheggaaa/pb"
	"flag"
	"log"
)

var from string
var to string
var limit int
var offset int

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write into")
	flag.IntVar(&limit, "limit", 0, "limit")
	flag.IntVar(&offset, "offset", 0, "offset")
}

func main() {
	from    := "go.mod"
	to      := "go.mod.copy"
	limit   := -1
	offset  := 0

	flag.Parse()

	err := Copy(from, to, limit, offset)
	if err != nil {
		log.Fatalln("Copy error: " + err.Error())
	}
}

func Copy(from string, to string, limit int, offset int) error{




	return nil
}


