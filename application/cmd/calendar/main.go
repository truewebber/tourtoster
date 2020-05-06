package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func main() {
	srv, err := calendar.NewService(context.Background(), option.WithCredentialsFile("/Users/truewebber/Downloads/calendar.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	calendarID := "primary"

	//// INSERT
	//event := &calendar.Event{
	//	Summary:     "Another Sample event",
	//	Location:    "Sample location",
	//	Description: "This is a sample event.",
	//	Start: &calendar.EventDateTime{
	//		DateTime: "2020-05-22T10:00:00+03:00",
	//		TimeZone: "Europe/Moscow",
	//	},
	//	End: &calendar.EventDateTime{
	//		DateTime: "2020-05-22T11:00:00+03:00",
	//		TimeZone: "Europe/Moscow",
	//	},
	//	//Recurrence: []string{
	//	//	"RRULE:FREQ=DAILY;",
	//	//},
	//}
	//e, err := srv.Events.Insert(calendarID, event).Do()
	//if err != nil {
	//	log.Fatalf("Unable to create event. %v\n", err)
	//}
	//fmt.Printf("Event created: %s\n", e.HtmlLink)

	// DELETE
	//if err := srv.Events.Delete(calendarID, "8btrqbj1uieofo2uets522bhmk").Do(); err != nil {
	//	log.Fatalf("Unable to create event. %v\n", err)
	//}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
		return
	}

	fmt.Println("Upcoming events:")
	for _, item := range events.Items {
		date := item.Start.DateTime
		if date == "" {
			date = item.Start.Date
		}
		fmt.Printf("[%s][%s] %v (%v)\n", item.Id, item.RecurringEventId, item.Summary, date)
	}
}
