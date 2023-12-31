package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"time"
)

type TemplateContent struct {
	Year int
}
type CountdownResponse struct {
	Body string
}

const (
	TIMEZONE     = "America/New_York"
	MYRTLE_DAYS  = 6
	HOURS_IN_DAY = 24
)

var templates = template.Must(template.ParseFiles("web/index.html", "web/style.css", "web/query.js"))

func getMyrtleLength() time.Duration {
	duration, err := time.ParseDuration(fmt.Sprintf("%dh", MYRTLE_DAYS*HOURS_IN_DAY))
	if err != nil {
		log.Fatalf("Failed to parse duration %d", duration)
	}
	return duration
}

func getLocationInTimezone() *time.Location {
	loc, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		log.Fatalf("Failed to load location of %s", TIMEZONE)
	}
	return loc
}

func getMyrtleTimes() (time.Time, time.Time) {
	loc := getLocationInTimezone()
	year, month, _ := time.Now().In(loc).Date()
	if month > time.May {
		year += 1
	}
	may_first_day := time.Date(year, time.May, 1, 0, 0, 0, 0, loc)
	days_till_myrtle := (7 - ((int)(may_first_day.Weekday() - time.Sunday))) * HOURS_IN_DAY
	days_in_duration, err := time.ParseDuration(fmt.Sprintf("%dh", days_till_myrtle))
	if err != nil {
		log.Fatalf("Failed to convert days (%d) to Duration", days_till_myrtle)
	}
	first_day_myrtle := may_first_day.Add(days_in_duration)
	last_day_myrtle := first_day_myrtle.Add(getMyrtleLength())
	return first_day_myrtle, last_day_myrtle
}

// Shamelessly stolen from time.go
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

func durationToHumanString(duration time.Duration) string {
	var buf [64]byte
	w := len(buf)
	u := (uint64)(math.Round(duration.Seconds()))
	w--
	buf[w] = 's'
	w = fmtInt(buf[:w], u%60)
	u /= 60

	w--
	buf[w] = ':'
	w--
	buf[w] = 'm'
	w = fmtInt(buf[:w], u%60)
	u /= 60

	w--
	buf[w] = ':'
	w--
	buf[w] = 'h'
	w = fmtInt(buf[:w], u%HOURS_IN_DAY)
	u /= HOURS_IN_DAY

	w--
	buf[w] = ':'
	w--
	buf[w] = 'd'
	w = fmtInt(buf[:w], u)

	return string(buf[w:])
}

func countdownHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().In(getLocationInTimezone())
	myrtleStart, myrtleEnd := getMyrtleTimes()
	start_difference := myrtleStart.Sub(now)
	end_difference := myrtleEnd.Sub(now)
	var body string
	if start_difference > 0 {
		body = fmt.Sprintf("Countdown to Myrtle: %s", durationToHumanString(start_difference))
	} else if end_difference > 0 {
		body = "It'S MYRTLE TIME!!!"
	} else {
		body = "See y'all next year!"
	}

	resp := CountdownResponse{body}

	enc := json.NewEncoder(w)
	err := enc.Encode(resp)
	if err != nil {
		log.Fatalf("Failed to write HTTP response: %s", err.Error())
	}
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().In(getLocationInTimezone())
	var year = now.Year()
	if now.Month() > time.May {
		year += 1
	}
	rendered := TemplateContent{year}
	err := templates.Execute(w, rendered)
	if err != nil {
		log.Fatalf("Failed to render html")
	}
}

func main() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	http.HandleFunc("/", htmlHandler)
	http.HandleFunc("/countdown", countdownHandler)
	log.Println("Starting server on port 8080 and listening")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
