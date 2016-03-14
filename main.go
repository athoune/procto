package main

import (
	"fmt"
	"log"
	"os"
	"time"
	//"strconv"
	//"time"
	"./fd"
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
	fmt.Println("Heap used, capacity, ratio : ", hsperfdata.Heap.Used, hsperfdata.Heap.Capacity, float64(100*hsperfdata.Heap.Used)/float64(hsperfdata.Heap.Capacity))

	f, err := fd.NewFd(data.Pid)
	if err != nil {
		log.Fatal(err)
	}
	sockets, err := f.CountSockets()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sockets", sockets)
	pipes, err := f.CountPipes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pipes", pipes)

	stats := proc.NewTimeStatThreads(data.Pid)
	err = stats.Measures()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	values, processors, err := stats.MeasuresAndCalculate()
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range values {
		fmt.Println(k, v)
	}
	fmt.Println("Processors", processors)

	//interval, err := strconv.Atoi(os.Args[2])
	//if err != nil {
	//log.Fatal(err)
	//}
	//inter := time.Duration(interval) * time.Millisecond

}
