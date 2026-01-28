package clistat

/*
 TODO - add prefix support, issue start in New()
*/

import (
	"log"
	"time"
)

type CliStat struct {
	Start   int64
	Ltime   int64
	Cnt     int64
	Last    int64
	Timeout int64
}

func New(timeout int64) CliStat {
	now := time.Now().Unix()
	return CliStat{now, now, 0, 0, timeout}
}

func (s *CliStat) Hit() {
	s.Cnt++
	// fmt.Printf("Cnt=%d", s.Cnt)
	if s.Cnt&255 != 0 {
		return
	}
	now := time.Now().Unix()
	ellapsed := now - s.Ltime

	if ellapsed > s.Timeout {
		log.Printf("Cnt=%d K. HPS=%d K\n", s.Cnt/1000, ((s.Cnt-s.Last)/ellapsed)/1000)
		s.Last = s.Cnt
		s.Ltime = now
	}
}

func (s *CliStat) Finish() {
	ellapsed := time.Now().Unix() - s.Start
	log.Printf("DONE, Hits: %d, Seconds: %d\n", s.Cnt, ellapsed)
}
