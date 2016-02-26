package main

import (
	"fmt"
	hsperfdata "github.com/YaSuenag/hsbeat/hsperfdata"
	"log"
	"os"
	//"strconv"
	//"time"
	"path"
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
	entries, err := perfData.ReadAllEntry(f)
	if err != nil {
		log.Fatal(err)
	}
	capacity := int64(0)
	used := int64(0)
	for _, entry := range entries {
		match, err := path.Match("sun/gc/generation/*/capacity", entry.EntryName)
		if err != nil {
			log.Fatal(err)
		}
		if match {
			capacity += entry.LongValue
		}
		match, err = path.Match("sun/gc/generation/*/space/*/used", entry.EntryName)
		if err != nil {
			log.Fatal(err)
		}
		if match {
			used += entry.LongValue
		}
		//fmt.Println(entry.EntryName, entry.LongValue, entry.StringValue)
	}
	fmt.Println(used, capacity, float64(100*(capacity-used))/float64(capacity))
	//interval, err := strconv.Atoi(os.Args[2])
	//if err != nil {
	//log.Fatal(err)
	//}
	//inter := time.Duration(interval) * time.Millisecond

}
