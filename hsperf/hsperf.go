package hsperf

import (
	hsperfdata "github.com/YaSuenag/hsbeat/hsperfdata"
	"os"
	//"strconv"
	//"time"
	"path"
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
