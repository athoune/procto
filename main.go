package main

import (
	"fmt"
	hsperfdata "github.com/YaSuenag/hsbeat/hsperfdata"
	"log"
	"os"
	//"strconv"
	//"time"
	"path"
	"path/filepath"
)

type HeapData struct {
	Capacity int64
	Used     int64
}

type Data struct {
	Heap HeapData
}

func ReadHsperf(hsPerfDataPath string) (*Data, error) {
	data := &Data{}
	data.Heap = HeapData{}
	perfData := &hsperfdata.HSPerfData{}
	perfData.ForceCachedEntryName = make(map[string]int)

	entry := "sun/os/hrt/frequency"
	perfData.ForceCachedEntryName[entry] = 1

	f, err := os.Open(hsPerfDataPath)
	if err != nil {
		return data, err
	}
	defer f.Close()

	err = perfData.ReadPrologue(f)
	if err != nil {
		return data, err
	}

	f.Seek(int64(perfData.Prologue.EntryOffset), os.SEEK_SET)
	entries, err := perfData.ReadAllEntry(f)
	if err != nil {
		return data, err
	}
	for _, entry := range entries {
		match, err := path.Match("sun/gc/generation/*/capacity", entry.EntryName)
		if err != nil {
			return nil, err
		}
		if match {
			data.Heap.Capacity += entry.LongValue
		}
		match, err = path.Match("sun/gc/generation/*/space/*/used", entry.EntryName)
		if err != nil {
			return nil, err
		}
		if match {
			data.Heap.Used += entry.LongValue
		}
		//fmt.Println(entry.EntryName, entry.LongValue, entry.StringValue)

		// Threads
		// java/threads/daemon 26
		// java/threads/live 28
		// java/threads/livePeak 29
		// java/threads/started 40

		// Classes
		// java/cls/loadedClasses 9371
		// java/cls/sharedLoadedClasses 0
		// java/cls/sharedUnloadedClasses 0
		// java/cls/unloadedClasses 0
	}
	return data, nil
}

func getHSPerfData(user string) (string, error) {
	paths, err := filepath.Glob(filepath.Join(os.TempDir(), "hsperfdata_"+user, "*"))
	if err != nil {
		return "", err
	}
	return paths[0], nil
}

func main() {
	user := os.Args[1]
	hsPerfDataPath, err := getHSPerfData(user)
	if err != nil {
		log.Fatal(err)
	}
	hsperfdata, err := ReadHsperf(hsPerfDataPath)
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
