package gigasecond

import "time"

// My solution
func AddGigasecond(t time.Time) time.Time {
	return t.Add(time.Second * 1000000000)
}

// Solution prefered 1
func AddGigasecond(t time.Time) time.Time {
   return t.Add(time.Second * 1e9)
}

// Solution prefered 2
const Gigasecond time.Duration = 1e18

func AddGigasecond(t time.Time) time.Time {
   return t.Add(Gigasecond)
}
