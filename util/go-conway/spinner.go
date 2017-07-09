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
}

func (s *Spinner) Tick() {
	txt := s.Seq[s.Pos]
	s.Pos = (s.Pos + 1) % len(s.Seq)
	goterm.Printf("%s", txt)
	goterm.MoveCursorUp(1)
	goterm.Flush()
}

func (s Spinner) Done() {
	goterm.Print(strings.Repeat(" ", len(s.Seq[s.Pos])))
	goterm.MoveCursorUp(1)
	goterm.Flush()
}

func (s Spinner) Spin(done <-chan bool) {
	ticker := time.NewTicker(s.Delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Tick()
		case <-done:
			s.Done()
			return
		}
	}
}
