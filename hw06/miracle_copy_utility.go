package main

import (
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
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

	flag.Parse()

	err := Copy(from, to, limit, offset)
	if err != nil {
		log.Fatalln("Copy error: " + err.Error())
	}
}

// Copy is the miracle function for copying one file to another
func Copy(from string, to string, limit int, offset int) error{
	fileInfo, err := os.Stat(from)
	if err != nil {
		return errors.Wrapf(err, "Error get stat of file %v", from)
	}
	fileFrom, err := os.OpenFile(from, os.O_RDONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "Can't open file %v", from)
	}
	fileTo, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrapf(err, "Can't open file %v", from)
	}

	count := 100
	len := fileInfo.Size() / int64(count)
	if fileInfo.Size() % 100 > 0 {
		len ++
	}
	bar := pb.StartNew(count)

	for i := 0; i < count; i++ {
		n, err := io.CopyN(fileTo, fileFrom, len)
		if n != len && err != nil && err != io.EOF {
			errors.Wrapf(err, "Error io.CopyN from=%v; to=%v", from, to)
		}
		bar.Increment()
		//time.Sleep(0.5 * time.Second)
	}
	bar.Finish()

	fmt.Println(count)
	return nil
}


