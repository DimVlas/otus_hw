package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	// var (
	// 	tt int64 = 1024
	// 	i  int64 = 0
	// )
	// pct := &pp.PctProgress{}
	// pct.Start(tt)
	// for ; i < tt; i += 10 {
	// 	pct.Increment(10)
	// 	time.Sleep(time.Second / 20)
	// }

	// pct.Finish()

	// time.Sleep(time.Millisecond * 500)

	if err := Copy(from, to, offset, limit); err != nil {
		fmt.Println(err)
		return
	}
}
