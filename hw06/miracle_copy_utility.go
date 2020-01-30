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

const approxCount = 30

var from    string
var to      string
var limit   uint
var offset  uint

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write into")
	flag.UintVar(&limit, "limit", 0, "limit")
	flag.UintVar(&offset, "offset", 0, "offset")
}

func main() {

	flag.Parse()

	err := Copy(from, to, limit, offset)
	if err != nil {
		log.Fatalln("Copy error: " + err.Error())
	}
}

// Copy is the miracle function for copying one file to another
func Copy(from string, to string, limit uint, offset uint) error{
	fileInfo, err := os.Stat(from)
	if err != nil {
		return errors.Wrapf(err, "Error get stat of file %v", from)
	}

	if fileInfo.Size() < int64(offset) {
		return errors.New("offset is more than file size")
	}

	fileFrom, err := os.OpenFile(from, os.O_RDONLY, 0644)
	defer fileFrom.Close()
	if err != nil {
		return errors.Wrapf(err, "Can't open file %v", from)
	}
	fileTo, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE, 0644)
	defer fileTo.Close()
	if err != nil {
		return errors.Wrapf(err, "Can't open file %v", from)
	}
	var sizeToCopy uint = uint(fileInfo.Size()) - offset

	if limit > 0 {
		if sizeToCopy > limit {
			sizeToCopy = limit
		}
	}

	var count uint = approxCount
	var len uint

	if sizeToCopy > count {
		len = sizeToCopy / count
		count = sizeToCopy / len
		if sizeToCopy % len > 0 {
			count ++
		}
	} else {
		count = sizeToCopy
		len = 1
	}

	bar := pb.StartNew(int(count))
	fileFrom.Seek(int64(offset), io.SeekStart)

	for i := 0; i < int(count); i++ {

		if i == int(count) - 1 && sizeToCopy % len > 0 {
			len = sizeToCopy % len
		}
		n, err := io.CopyN(fileTo, fileFrom, int64(len))
		if uint(n) != len && err != nil && err != io.EOF {
			errors.Wrapf(err, "Error io.CopyN from=%v; to=%v", from, to)
		}
		bar.Increment()
		//time.Sleep(0.5 * time.Second)
	}
	bar.Finish()

	fmt.Println(count)
	return nil
}


