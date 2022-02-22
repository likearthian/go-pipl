package pipl

import (
	"time"
)

type executionStat struct {
	executionsCounter  int
	totalExecutionTime float64
	avgExecutionTime   float64
}

func (s *executionStat) recordExecution(foo func()) {
	s.executionsCounter++
	st := time.Now()
	foo()
	s.totalExecutionTime += time.Now().Sub(st).Seconds()
}

func (s *executionStat) calculate() {
	if s.executionsCounter > 0 {
		s.avgExecutionTime = (s.totalExecutionTime / float64(s.executionsCounter))
	}
}
