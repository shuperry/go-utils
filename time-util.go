package utils

import (
	"log"
	"strconv"
	"time"
)

const (
	HOUR_FMT               = "2006-01-02 15h"
	DATE_FMT               = "2006-01-02"
	DATE_TIME_FMT          = "2006-01-02 15:04:05"
	MONTH_FMT              = "2006-01"
	YEAR_FMT               = "2006"
	DATE_TIME_MILL_SEC_FMT = "2006-01-02 15:04:05.000"
)

func ParseUTCTimestampToDatetimeStr(ms string) string {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		log.Panic(err)
	}

	tm := time.Unix(0, msInt*int64(time.Millisecond))

	return tm.Format(DATE_TIME_MILL_SEC_FMT)
}

func ParseUTCTimestampToTime(ms string) time.Time {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		log.Panic(err)
	}

	tm := time.Unix(0, msInt*int64(time.Millisecond))

	return tm
}

func ParseUnixTimestampToDatetimeStr(s string) string {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Panic(err)
	}

	tm := time.Unix(i, 0)

	return tm.Format(DATE_TIME_MILL_SEC_FMT)
}

func ParseUnixTimestampToTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Panic(err)
	}

	tm := time.Unix(i, 0)

	return tm
}

func GetNowDateStr() string {
	return time.Now().Format(DATE_FMT)
}

func GetNowDateTimeStr() string {
	return time.Now().Format(DATE_TIME_FMT)
}

func GetDateStrByTime(t time.Time) string {
	return t.Format(DATE_FMT)
}

func GetDateTimeStrByTime(t time.Time) string {
	return t.Format(DATE_TIME_FMT)
}

/**
return now time utc-timestamp(ms).
*/
func GetNowUTCTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

/**
return now time unix-timestamp(s).
*/
func GetNowUnixTimestamp() int64 {
	return time.Now().Unix()
}
