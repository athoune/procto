package main

import (
	"fmt"
	"log"
	"os"
	//"strconv"
	//"time"
	"./hsperf"
	"./java"
	"./proc"
)

func main() {
	user := os.Args[1]
	data, err := java.FindAppData(user)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	hsperfdata, err := hsperf.ReadHsperf(data.Path())
	if err != nil {
		log.Fatal(err)
	}
	proc, err := proc.ReadStat(data.Pid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Heap used, capacity, ratio : ", hsperfdata.Heap.Used, hsperfdata.Heap.Capacity, float64(100*hsperfdata.Heap.Used)/float64(hsperfdata.Heap.Capacity))
	fmt.Println("CPU time User System : ", proc.Utime, proc.Stime)
	//interval, err := strconv.Atoi(os.Args[2])
	//if err != nil {
	//log.Fatal(err)
	//}
	//inter := time.Duration(interval) * time.Millisecond

}
