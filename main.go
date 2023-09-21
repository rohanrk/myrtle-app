package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"time"
)

type Content struct {
	Year int
	Body string
}

const (
	TIMEZONE     = "America/New_York"
	MYRTLE_DAYS  = 6
	HOURS_IN_DAY = 24
)

var templates = template.Must(template.ParseFiles("index.html"))

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
		log.Fatal("Failed to calculate time to Myrtle")
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
	fmt.Printf("secs: %d\nremainder: %d\n", u, u%60)
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

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().In(getLocationInTimezone())
	myrtleStart, myrtleEnd := getMyrtleTimes()
	log.Printf("now: %s\nmyrtle start: %s\nmyrtle_end: %s\n", now.String(), myrtleStart.String(), myrtleEnd.String())
	start_difference := myrtleStart.Sub(now)
	end_difference := myrtleEnd.Sub(now)
	log.Printf("start: %d\nend: %d", start_difference, end_difference)
	var body string
	if start_difference > 0 {
		body = fmt.Sprintf("Countdown to Myrtle: %s", durationToHumanString(start_difference))
	} else if end_difference > 0 {
		body = "It'S MYRTLE TIME!!!"
	} else {
		body = "See y'all next year!"
	}

	log.Printf("Got body: %s\n", body)

	rendered := Content{now.Year(), body}
	err := templates.Execute(w, rendered)
	if err != nil {
		log.Fatalf("Failed to build html template with %s", err.Error())
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Starting server on port 8080 and listening")
	log.Fatal(http.ListenAndServe(":8080", nil))
	// myrtleStart, _ := getMyrtleTimes()
	// start_difference := myrtleStart.Sub(time.Now())
	// log.Println(durationToHumanString(start_difference))
	// log.Println(start_difference.String())
}
