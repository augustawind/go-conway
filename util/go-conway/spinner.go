package util

import (
	"strings"
	"time"

	"github.com/buger/goterm"
)

type Spinner struct {
	Seq   []string
	Delay time.Duration
	Pos   int
	done  chan bool
}

func (s *Spinner) Tick() {
	txt := s.Seq[s.Pos]
	s.Pos = (s.Pos + 1) % len(s.Seq)
	goterm.Printf("%s", txt)
	goterm.MoveCursorUp(1)
	goterm.Flush()
}

func (s Spinner) Done() {
	s.done <- true
	goterm.Print(strings.Repeat(" ", len(s.Seq[s.Pos])))
	goterm.MoveCursorUp(1)
	goterm.Flush()
}

func (s Spinner) Spin() {
	defer s.Done()

	ticker := time.NewTicker(s.Delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Tick()
		case <-s.done:
			return
		}
	}
}
