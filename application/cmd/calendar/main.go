package main

import (
	"time"

	"github.com/teambition/rrule-go"
)

func main() {
	r, rErr := rrule.NewRRule(rrule.ROption{
		RFC:       true,
		Byweekday: []rrule.Weekday{rrule.WE, rrule.SU},
	})
	if rErr != nil {
		println("error create RRule", rErr.Error())
		return
	}

	exR, exRErr := rrule.NewRRule(rrule.ROption{
		Bymonth:    []int{7},
		Bymonthday: []int{19},
	})
	if exRErr != nil {
		println("error create RRule", exRErr.Error())
		return
	}

	s := rrule.Set{}
	s.RRule(r)
	s.ExRule(exR)

	rruleStr := s.String()
	println(rruleStr)

	s2, s2Err := rrule.StrToRRuleSet(rruleStr)
	//s2, s2Err := rrule.StrToRRuleSet("RRULE:FREQ=WEEKLY;COUNT=30;INTERVAL=1;WKST=MO")
	if s2Err != nil {
		println("error create RRule Set", s2Err.Error())
		return
	}

	cTime := time.Now()
	dStart := time.Date(cTime.Year(), cTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	dEnd := time.Date(cTime.Year(), cTime.Month()+3, 1, 23, 59, 59, 0, time.UTC)
	dEnd = dEnd.AddDate(0, 0, -1)

	println(dStart.String())
	println(dEnd.String())

	all := s2.Between(dStart, dEnd, true)
	for i := 0; i < len(all); i++ {
		println(all[i].Format("[Mon] 02 Jan 2006"))
	}
}
