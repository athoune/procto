package main

import (
	"fmt"
	hsperfdata "github.com/YaSuenag/hsbeat/hsperfdata"
	"log"
	"os"
	//"strconv"
	//"time"
)

func main() {
	pid := os.Args[1]
	hsPerfDataPath, err := hsperfdata.GetHSPerfDataPath(pid)
	if err != nil {
		log.Fatal(err)
	}
	perfData := &hsperfdata.HSPerfData{}
	perfData.ForceCachedEntryName = make(map[string]int)

	entry := "sun/os/hrt/frequency"
	perfData.ForceCachedEntryName[entry] = 1

	f, err := os.Open(hsPerfDataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = perfData.ReadPrologue(f)
	if err != nil {
		log.Fatal(err)
	}

	f.Seek(int64(perfData.Prologue.EntryOffset), os.SEEK_SET)
	result, err := perfData.ReadAllEntry(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Sprintf("hsperfdata : %+v", result)

	//interval, err := strconv.Atoi(os.Args[2])
	//if err != nil {
	//log.Fatal(err)
	//}
	//inter := time.Duration(interval) * time.Millisecond

}
