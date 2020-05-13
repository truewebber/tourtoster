package main

import (
	"time"

	"github.com/teambition/rrule-go"
)

func main() {
	eventTimeLoc, locErr := time.LoadLocation("Europe/Moscow")
	if locErr != nil {
		println("error create load location", locErr.Error())
		return
	}

	r, rErr := rrule.NewRRule(rrule.ROption{
		Freq:      rrule.DAILY,
		Until:     time.Date(2021, time.May, 9, 23, 59, 59, 0, eventTimeLoc),
		Byweekday: []rrule.Weekday{rrule.WE, rrule.TH, rrule.FR, rrule.SA, rrule.SU},
		Wkst:      rrule.MO,
		RFC:       true,
	})
	if rErr != nil {
		println("error create RRule", rErr.Error())
		return
	}

	s := rrule.Set{}
	s.RRule(r)

	rruleStr := s.String()
	println(rruleStr)

	s2, s2Err := rrule.StrToRRuleSet(rruleStr)
	if s2Err != nil {
		println("error create RRule Set", s2Err.Error())
		return
	}

	cTime := time.Now().In(eventTimeLoc)

	dStart := time.Date(cTime.Year(), cTime.Month(), 1, 0, 0, 0, 0, cTime.Location())
	dEnd := time.Date(cTime.Year(), cTime.Month()+1, 1, 23, 59, 59, 0, cTime.Location())
	dEnd = dEnd.AddDate(0, 0, -1)

	println(dStart.Location().String())
	println(dStart.String())
	println(dEnd.String())

	all := s2.Between(dStart, dEnd, true)
	for i := 0; i < len(all); i++ {
		println(all[i].Format("[Mon] 02 Jan 2006"))
	}
}
