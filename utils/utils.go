package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// CheckError handles error
func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

type TimePoint struct {
	Week, DayOfWeek, Hour, Min int
}

func (t *TimePoint) AddMinutes(min int) {
	t.Min += min

	if t.Min < 60 {
		return
	}
	t.Min %= 60
	t.Hour++

	if t.Hour < 24 {
		return
	}
	t.Hour %= 24
	t.DayOfWeek++
	if t.DayOfWeek < 7 {
		return
	}
	t.DayOfWeek %= 7
	t.Week++

}

type Interval struct {
	Start, End TimePoint
}

func (in *Interval) String() string {
	return fmt.Sprintf("Start %d end %d", in.Start, in.End)
}

// GetIndex prints the index of an interval stepMinutesLength is the amount of minutes from one evalutaion point to another
func (t TimePoint) GetIndex(stepMinutesResolution int) int {
	return (t.Week*10080 + t.DayOfWeek*1440 + t.Hour*60 + t.Min) / stepMinutesResolution
}

// TimeToHM return hours and minutes from HH:mm
func TimeToHM(in string) (hours, minutes int) {
	hm := strings.Split(in, ":")

	hours, errH := strconv.Atoi(hm[0])
	if errH != nil {
		log.Fatal("could not convert hours")
	}

	if hours < 0 || hours > 23 {
		log.Fatal("Hour value is outside limit")
	}

	minutes, errM := strconv.Atoi(hm[1])
	if errM != nil {
		log.Fatal("could not convert minutes")
	}
	if minutes < 0 || minutes > 59 {
		log.Fatal("Hour value is outside limit")
	}
	return hours, minutes
}

func ExtractInterval(s string, weekOffset int, dayOffset int) (in Interval) {
	se := strings.Split(s, "-")
	startH, startM := TimeToHM(se[0])
	endH, endM := TimeToHM(se[1])
	return Interval{TimePoint{weekOffset, dayOffset, startH, startM}, TimePoint{weekOffset, dayOffset, endH, endM}}
}

func ExtractDailyIntervals(s string, weekOffset int, dayOffset int) (ins []Interval) {
	inDayIntervals := strings.Split(s, " e ")
	ret := make([]Interval, 0)
	for _, v := range inDayIntervals {
		// transform hour to interval index
		ret = append(ret, ExtractInterval(v, weekOffset, dayOffset))
	}
	return ret
}
