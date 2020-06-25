package tour

import (
	"time"

	"github.com/teambition/rrule-go"

	"tourtoster/user"
)

type (
	Tour struct {
		ID                             int64          `json:"id"`
		Type                           Type           `json:"type"`
		Creator                        *user.User     `json:"creator"`
		Status                         Status         `json:"status"`
		Recurrence                     *rrule.Set     `json:"recurrence"`
		Title                          string         `json:"title"`
		Image                          string         `json:"image"`
		Description                    string         `json:"description"`
		Map                            string         `json:"map"`
		MaxPersons                     int            `json:"max_persons"`
		PricePerChildrenThreeSix       Price          `json:"price_per_children_three_six"`
		PricePerChildrenZeroSix        Price          `json:"price_per_children_zero_six"`
		PricePerChildrenSevenSeventeen Price          `json:"price_per_children_seven_seventeen"`
		PricePerAdults                 Price          `json:"price_per_adults"`
		TimeTable                      []TimeTableRow `json:"timetable"`
		Features                       []Feature      `json:"features"`
		Highlights                     []Highlight    `json:"highlights"`
		FAQs                           []FAQ          `json:"faqs"`
	}

	Highlight struct {
	}

	FAQ struct {
	}

	TimeTableRow struct {
		Time       time.Time `json:"time"`
		GroupCount int       `json:"group_count"`
	}
)
