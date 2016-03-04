//Stolen from https://github.com/uber-common/cpustat

package proc

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ProcStats struct {
	CaptureTime         time.Time
	PrevTime            time.Time
	Pid                 uint64
	Comm                string
	State               string
	Ppid                uint64
	Pgrp                int64
	Session             int64
	TtyNr               int64
	Tpgid               int64
	Flags               uint64
	Minflt              uint64
	Cminflt             uint64
	Majflt              uint64
	Cmajflt             uint64
	Utime               uint64
	Stime               uint64
	Cutime              uint64
	Cstime              uint64
	Priority            int64
	Nice                int64
	NumThreads          uint64
	StartTime           uint64
	Vsize               uint64
	Rss                 uint64
	Rsslim              uint64
	Processor           uint64
	RtPriority          uint64
	Policy              uint64
	DelayacctBlkioTicks uint64
	GuestTime           uint64
	CguestTime          uint64
}

// you might think that we could split on space, but due to what can at best be called
// a shortcoming of the /proc/pid/stat format, the comm field can have unescaped spaces, parens, etc.
// This may be a bit paranoid, because even many common tools like htop do not handle this case well.
func procPidStatSplit(line string) []string {
	line = strings.TrimSpace(line)
	parts := make([]string, 52)

	partnum := 0
	strpos := 0
	start := 0
	inword := false
	space := " "[0]
	open := "("[0]
	close := ")"[0]
	groupchar := space

	for ; strpos < len(line); strpos++ {
		if inword {
			if line[strpos] == space && (groupchar == space || line[strpos-1] == groupchar) {
				parts[partnum] = line[start:strpos]
				partnum++
				start = strpos
				inword = false
			}
		} else {
			if line[strpos] == open {
				groupchar = close
				inword = true
				start = strpos
				strpos = strings.LastIndex(line, ")") - 1
				if strpos <= start { // if we can't parse this insane field, skip to the end
					strpos = len(line)
					inword = false
				}
			} else if line[strpos] != space {
				groupchar = space
				inword = true
				start = strpos
			}
		}
	}

	if inword {
		parts[partnum] = line[start:strpos]
	}

	return parts
}

func ReadStat(pid int) (ProcStats, error) {

	lines, err := readFileLines(fmt.Sprintf("/proc/%d/stat", pid))
	// pid could have exited between when we scanned the dir and now
	if err != nil {
		return ProcStats{}, nil
	}
	// this format of this file is insane because comm can have split chars in it
	parts := procPidStatSplit(lines[0])

	stat := ProcStats{
		time.Now(), // this might be expensive. If so, can cache it. We only need 1ms resolution
		time.Time{},
		readUInt(parts[0]),                  // pid
		strings.Map(stripSpecial, parts[1]), // comm
		parts[2],            // state
		readUInt(parts[3]),  // ppid
		readInt(parts[4]),   // pgrp
		readInt(parts[5]),   // session
		readInt(parts[6]),   // tty_nr
		readInt(parts[7]),   // tpgid
		readUInt(parts[8]),  // flags
		readUInt(parts[9]),  // minflt
		readUInt(parts[10]), // cminflt
		readUInt(parts[11]), // majflt
		readUInt(parts[12]), // cmajflt
		readUInt(parts[13]), // utime
		readUInt(parts[14]), // stime
		readUInt(parts[15]), // cutime
		readUInt(parts[16]), // cstime
		readInt(parts[17]),  // priority
		readInt(parts[18]),  // nice
		readUInt(parts[19]), // num_threads
		// itrealvalue - not maintained
		readUInt(parts[21]), // starttime
		readUInt(parts[22]), // vsize
		readUInt(parts[23]), // rss
		readUInt(parts[24]), // rsslim
		// bunch of stuff about memory addresses
		readUInt(parts[38]), // processor
		readUInt(parts[39]), // rt_priority
		readUInt(parts[40]), // policy
		readUInt(parts[41]), // delayacct_blkio_ticks
		readUInt(parts[42]), // guest_time
		readUInt(parts[43]), // cguest_time
	}
	return stat, nil
}

// note that this is not thread safe
var buf *bytes.Buffer

// ReadSmallFile is like os.ReadFile but dangerously optimized for reading files from /proc.
// The file is not statted first, and the same buffer is used every time.
func ReadSmallFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		f.Close()
		return nil, err
	}

	if buf == nil {
		buf = bytes.NewBuffer(make([]byte, 0, 8192))
	} else {
		buf.Reset()
	}
	_, err = buf.ReadFrom(f)
	f.Close()
	return buf.Bytes(), err
}

// Read a small file and split on newline
func readFileLines(filename string) ([]string, error) {
	file, err := ReadSmallFile(filename)
	if err != nil {
		return nil, err
	}

	fileStr := strings.TrimSpace(string(file))
	return strings.Split(fileStr, "\n"), nil
}

// pull a float64 out of a string
func readFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)

	}
	return val
}

// pull a uint64 out of a string
func readUInt(str string) uint64 {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return val
}

// pull a int64 out of a string
func readInt(str string) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

// remove grouping characters that confuse the termui parser
func stripSpecial(r rune) rune {
	if r == '[' || r == ']' || r == '(' || r == ')' {
		return -1
	}
	return r
}
