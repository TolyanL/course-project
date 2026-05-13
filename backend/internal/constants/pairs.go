package constants

import "time"

type Pair struct {
	Number    int
	StartTime time.Time
	EndTime   time.Time
}

var Pairs = []Pair{
	{Number: 1, StartTime: toTime(8, 0), EndTime: toTime(9, 35)},
	{Number: 2, StartTime: toTime(9, 45), EndTime: toTime(11, 20)},
	{Number: 3, StartTime: toTime(12, 30), EndTime: toTime(14, 5)},
	{Number: 4, StartTime: toTime(14, 15), EndTime: toTime(15, 50)},
	{Number: 5, StartTime: toTime(16, 0), EndTime: toTime(17, 35)},
	{Number: 6, StartTime: toTime(17, 45), EndTime: toTime(19, 20)},
	{Number: 7, StartTime: toTime(19, 30), EndTime: toTime(21, 5)},
}

func toTime(hour, min int) time.Time {
	return time.Date(0, 1, 1, hour, min, 0, 0, time.UTC)
}

func GetPairByNumber(n int) *Pair {
	for i := range Pairs {
		if Pairs[i].Number == n {
			return &Pairs[i]
		}
	}
	return nil
}
