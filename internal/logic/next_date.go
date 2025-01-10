package logic

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var maxDayInMonth = map[int]int{
	1:  31,
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

func isLeapYear(year int) bool {
	if year%4 != 0 {
		return false
	}
	if year%100 != 0 {
		return true
	}
	if year%400 != 0 {
		return false
	}
	return true
}
func sortDays(days []int) []int {
	sort.Slice(days, func(i, j int) bool {
		// Если один из чисел положительный, сортируем их первыми
		if days[i] >= 0 && days[j] < 0 {
			return true
		}
		if days[i] < 0 && days[j] >= 0 {
			return false
		}
		// В остальном сортируем по возрастанию
		return days[i] < days[j]
	})
	return days
}

func nextDateDay(now, nextDate time.Time, repeat string) (string, error) {
	var count int
	step, err := strconv.Atoi(repeat[2:])
	if step > 400 || err != nil {
		return "", fmt.Errorf("repeat is invalid")
	}
	for now.After(nextDate) || now.Equal(nextDate) {
		nextDate = nextDate.AddDate(0, 0, step)
		count++
	}
	if count == 0 {
		nextDate = nextDate.AddDate(0, 0, step)
	}
	return nextDate.Format(DateFormat), nil
}

func nextDateWeek(now, nextDate time.Time, repeat string) (string, error) {
	nowWeekday := int(now.Weekday())
	daysRaw := strings.Split(repeat[2:], ",")
	if nextDate.Before(now) {
		nextDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}
	for _, c := range daysRaw {
		day, err := strconv.Atoi(string(c))
		if err != nil || day > 7 || day < 1 {
			return "", fmt.Errorf("repeat is invalid")
		}
		if nowWeekday < day {
			return nextDate.AddDate(0, 0, day-nowWeekday).Format(DateFormat), nil
		}
	}
	firstDay, err := strconv.Atoi(daysRaw[0])
	if err != nil {
		return "", fmt.Errorf("repeat is invalid")
	}
	return nextDate.AddDate(0, 0, 7-nowWeekday+firstDay).Format(DateFormat), nil
}

func getDaysInMonth(repeatRaw string) ([]int, error) {
	daysRawStr := strings.Split(repeatRaw, ",")
	var daysRaw []int
	for _, c := range daysRawStr {
		day, err := strconv.Atoi(c)
		if err != nil || day > 31 || day < -2 {
			return []int{}, fmt.Errorf("repeat is invalid")
		}
		daysRaw = append(daysRaw, day)
	}
	daysRaw = sortDays(daysRaw)
	return daysRaw, nil
}

func getMonths(repeatRaw string) ([]int, error) {
	var monthsRaw []int
	monthsRawStr := strings.Split(repeatRaw, ",")
	for _, m := range monthsRawStr {
		month, err := strconv.Atoi(m)
		if err != nil || month > 12 || month < 1 {
			return []int{}, fmt.Errorf("repeat is invalid")
		}
		monthsRaw = append(monthsRaw, month)
	}
	sort.Ints(monthsRaw)
	return monthsRaw, nil
}

func nextDateMonth(now, nextDate time.Time, repeat string) (string, error) {
	repeatRaw := strings.Split(repeat[2:], " ")

	daysRaw, err := getDaysInMonth(repeatRaw[0])
	if err != nil {
		return "", err
	}

	var monthsRaw []int
	if len(repeatRaw) == 2 {
		monthsRaw, err = getMonths(repeatRaw[1])
		if err != nil {
			return "", err
		}
	}

	if isLeapYear(now.Year()) {
		maxDayInMonth[2] = 29
	}

	var day int
	if monthsRaw == nil {
		if nextDate.Before(now) {
			nextDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		}
		for {
			for _, day = range daysRaw {
				if day == -1 {
					day = maxDayInMonth[int(nextDate.Month())]
				}
				if day == -2 {
					day = maxDayInMonth[int(nextDate.Month())] - 1
				}
				if day > maxDayInMonth[int(nextDate.Month())] {
					continue
				}
				nextDate = time.Date(nextDate.Year(), nextDate.Month(), day, 0, 0, 0, 0, time.UTC)
				if nextDate.After(now) {
					return nextDate.Format(DateFormat), nil
				}
			}
			nextDate = nextDate.AddDate(0, 1, 0)
		}
	} else {
		var monthIndex int
		for i, month := range monthsRaw {
			if month >= int(now.Month()) {
				if nextDate.Before(now) {
					nextDate = time.Date(now.Year(), time.Month(month), daysRaw[0], 0, 0, 0, 0, time.UTC)
				} else {
					nextDate = time.Date(nextDate.Year(), time.Month(month), daysRaw[0], 0, 0, 0, 0, time.UTC)
				}
				monthIndex = i
				break
			}
		}

		for i := monthIndex; i < len(monthsRaw); i++ {
			for _, day = range daysRaw {
				if day == -1 {
					day = maxDayInMonth[int(nextDate.Month())]
				}
				if day == -2 {
					day = maxDayInMonth[int(nextDate.Month())] - 1
				}
				if day > maxDayInMonth[int(nextDate.Month())] {
					continue
				}
				nextDate = time.Date(nextDate.Year(), time.Month(monthsRaw[i]), day, 0, 0, 0, 0, time.UTC)
				if nextDate.After(now) {
					return nextDate.Format(DateFormat), nil
				}
			}
		}
		return time.Date(nextDate.Year()+1, time.Month(monthsRaw[0]), daysRaw[0], 0, 0, 0, 0, time.UTC).Format(DateFormat), nil
	}
}

func nextDateYear(now, nextDate time.Time) (string, error) {
	var count int
	for now.After(nextDate) || now.Equal(nextDate) {
		nextDate = nextDate.AddDate(1, 0, 0)
		count++
	}
	if count == 0 {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format(DateFormat), nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", fmt.Errorf("repeat is empty")
	}

	nextDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("date is invalid")
	}

	switch repeat[0] {
	case 'd':
		if len(repeat) < 3 {
			return "", fmt.Errorf("repeat is invalid")
		}
		return nextDateDay(now, nextDate, repeat)
	case 'w':
		if len(repeat) < 3 {
			return "", fmt.Errorf("repeat is invalid")
		}
		return nextDateWeek(now, nextDate, repeat)
	case 'm':
		if len(repeat) < 3 {
			return "", fmt.Errorf("repeat is invalid")
		}
		return nextDateMonth(now, nextDate, repeat)
	case 'y':
		return nextDateYear(now, nextDate)
	}
	return "", fmt.Errorf("repeat is invalid")
}
