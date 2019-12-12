package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"time"
)

func main(){
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalln("ntp.QueryWithOptions() error: " + err.Error())
	}
	time := time.Now().Add(response.ClockOffset)
	fmt.Println(time.String())
}