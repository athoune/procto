package main

import (
	"fmt"
	"log"
	"os"
	//"strconv"
	//"time"
	"./hsperf"
)

func main() {
	user := os.Args[1]
	hsPerfDataPath, err := hsperf.FindHSPerfData(user)
	if err != nil {
		log.Fatal(err)
	}
	hsperfdata, err := hsperf.ReadHsperf(hsPerfDataPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hsperfdata.Heap.Used, hsperfdata.Heap.Capacity, float64(100*hsperfdata.Heap.Used)/float64(hsperfdata.Heap.Capacity))
	//interval, err := strconv.Atoi(os.Args[2])
	//if err != nil {
	//log.Fatal(err)
	//}
	//inter := time.Duration(interval) * time.Millisecond

}
