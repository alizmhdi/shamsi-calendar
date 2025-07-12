package calendar

import (
	"time"
)

const (
	gregorianOffset = 621
	julianDayOffset = 79
	cycle33Years = 12053
	cycle4Years = 1461
	cycle400Years = 146097
	cycle100Years = 36524
	cycle4YearsIn100 = 1461
	firstHalfDays = 186 // 6 * 31
	esfandMonth = 12
	leapYearIndicator = 0
)

type JalaliDate struct {
	Year  int
	Month int
	Day   int
}

var daysInMonth = []int{31, 31, 31, 31, 31, 31, 30, 30, 30, 30, 30, 29}

// Calendar breaks for leap year calculations
// These years mark boundaries where the leap year pattern changes
var breaks = []int{-61, 9, 38, 199, 426, 686, 756, 818, 1111, 1181, 1210, 1635, 2060, 2097, 2192, 2262, 2324, 2394, 2456, 3178}

// Gregorian month day offsets (non-leap year)
var gregorianMonthOffsets = [...]int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}

// Gregorian month day offsets (leap year)
var gregorianMonthOffsetsLeap = [...]int{0, 31, 60, 91, 121, 152, 182, 213, 244, 274, 305, 335}

// jalCalResult holds the result of Jalali calendar calculations
type jalCalResult struct {
	leap  int // Leap year indicator (0 = leap year)
	gy    int // Gregorian year
	march int // March offset
}

// div returns integer division
func div(a, b int) int {
	return int(a / b)
}

// isGregorianLeapYear checks if a Gregorian year is a leap year
func isGregorianLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// jalCal calculates Jalali calendar parameters for a given Jalali year
func jalCal(jy int) jalCalResult {
	bl := len(breaks)
	gy := jy + gregorianOffset
	leapJ := -14
	jp := breaks[0]
	jump := 0
	leap := 0

	// Calculate leap years based on breaks
	for i := 1; i < bl; i++ {
		jm := breaks[i]
		jump = jm - jp
		if jy < jm {
			break
		}
		leapJ = leapJ + div(jump, 33)*8 + div((jump%33), 4)
		jp = jm
	}

	n := jy - jp
	leapJ = leapJ + div(n, 33)*8 + div((n%33)+3, 4)
	if (jump%33) == 4 && (jump-n) == 4 {
		leapJ++
	}

	leapG := div(gy, 4) - div((div(gy, 100)+1)*3, 4) - 150
	march := 20 + leapJ - leapG

	if jump-n < 6 {
		n = n - jump + div(jump+4, 33)*33
	}

	leap = (((n + 1) % 33) - 1) % 4
	if leap == -1 {
		leap = 4
	}

	return jalCalResult{leap: leap, gy: gy, march: march}
}

// calculateGregorianDayNumber calculates the Julian Day Number for a Gregorian date
func calculateGregorianDayNumber(gy, gm, gd int) int {
	gy2 := gy - 1600
	gm2 := gm - 1
	gd2 := gd - 1

	gDayNo := 365*gy2 + div(gy2+3, 4) - div(gy2+99, 100) + div(gy2+399, 400)
	gDayNo += gregorianMonthOffsets[gm2] + gd2

	if gm > 2 && isGregorianLeapYear(gy) {
		gDayNo++
	}

	return gDayNo
}

// GregorianToJalali converts Gregorian date to Jalali date
// This is an accurate port from jalaali-js and shams Rust project
func GregorianToJalali(gy, gm, gd int) JalaliDate {
	// Calculate the Julian Day Number for the Gregorian date
	gDayNo := calculateGregorianDayNumber(gy, gm, gd)

	// Convert to Jalali
	jDayNo := gDayNo - julianDayOffset
	jNp := div(jDayNo, cycle33Years) // 12053 = 33 years
	jDayNo = jDayNo % cycle33Years
	jy := 979 + 33*jNp + 4*div(jDayNo, cycle4Years)
	jDayNo = jDayNo % cycle4Years

	if jDayNo >= 366 {
		jy += div(jDayNo-1, 365)
		jDayNo = (jDayNo - 1) % 365
	}

	var jm, jd int
	if jDayNo < firstHalfDays {
		jm = 1 + div(jDayNo, 31)
		jd = 1 + (jDayNo % 31)
	} else {
		jm = 7 + div(jDayNo-firstHalfDays, 30)
		jd = 1 + ((jDayNo - firstHalfDays) % 30)
	}

	return JalaliDate{Year: jy, Month: jm, Day: jd}
}

// JalaliToGregorian converts Jalali date to Gregorian date
func JalaliToGregorian(jy, jm, jd int) (int, int, int) {
	gy := 0
	if jy > 979 {
		gy = 1600
		jy -= 979
	} else {
		gy = 621
	}

	jCal := jalCal(jy)

	// Calculate total days from Jalali epoch
	days := 365*jy + div(jy, 33)*8 + div((jy%33)+3, 4)

	// Add days for months before the current month
	for i := 0; i < jm-1; i++ {
		if i == esfandMonth-1 && IsJalaliLeapYear(jy) {
			days += 30 // Esfand in leap year has 30 days
		} else {
			days += daysInMonth[i]
		}
	}

	days += jd - 1
	gy += jCal.march + days

	// Convert back to Gregorian
	gDayNo := gy
	gy = 400 * div(gDayNo, cycle400Years)
	gDayNo = gDayNo % cycle400Years

	leap := true
	if gDayNo >= 36525 {
		gDayNo--
		gy += 100 * div(gDayNo, cycle100Years)
		gDayNo = gDayNo % cycle100Years
		if gDayNo >= 365 {
			gDayNo++
		} else {
			leap = false
		}
	}

	gy += 4 * div(gDayNo, cycle4YearsIn100)
	gDayNo = gDayNo % cycle4YearsIn100

	if gDayNo >= 366 {
		leap = false
		gDayNo--
		gy += div(gDayNo, 365)
		gDayNo = gDayNo % 365
	}

	// Find month and day
	gm, gd := 0, 0
	var monthOffsets []int
	if leap {
		monthOffsets = gregorianMonthOffsetsLeap[:]
	} else {
		monthOffsets = gregorianMonthOffsets[:]
	}

	for i := 11; i >= 0; i-- {
		if gDayNo >= monthOffsets[i] {
			gm = i + 1
			gd = gDayNo - monthOffsets[i] + 1
			break
		}
	}

	return gy, gm, gd
}

// IsJalaliLeapYear determines if a Jalali year is a leap year using the accurate algorithm
func IsJalaliLeapYear(jy int) bool {
	return jalCal(jy).leap == leapYearIndicator
}

// GetDaysInMonth returns the number of days in a given Jalali month
func GetDaysInMonth(year, month int) int {
	if month == esfandMonth && IsJalaliLeapYear(year) {
		return 30 // Esfand in leap year has 30 days
	}
	return daysInMonth[month-1]
}

// GetDayOfWeek returns the day of week (0=Sunday, 1=Monday, etc.)
func GetDayOfWeek(year, month, day int) int {
	gYear, gMonth, gDay := JalaliToGregorian(year, month, day)
	t := time.Date(gYear, time.Month(gMonth), gDay, 0, 0, 0, 0, time.UTC)
	return (int(t.Weekday()) + 6) % 7
}

// GetMonthCalendar returns a 2D array representing the calendar for a month
func GetMonthCalendar(year, month int) [][]int {
	daysInMonth := GetDaysInMonth(year, month)
	firstDayOfWeek := GetDayOfWeek(year, month, 1)

	// Calculate number of weeks needed
	weeks := (daysInMonth + firstDayOfWeek + 6) / 7

	calendar := make([][]int, weeks)
	for i := range calendar {
		calendar[i] = make([]int, 7)
	}

	day := 1
	for week := 0; week < weeks; week++ {
		for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
			if (week == 0 && dayOfWeek < firstDayOfWeek) || day > daysInMonth {
				calendar[week][dayOfWeek] = 0
			} else {
				calendar[week][dayOfWeek] = day
				day++
			}
		}
	}

	return calendar
}
