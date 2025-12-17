package providers

import "time"

// RealTimeProvider implements TimeProvider using actual system time
type RealTimeProvider struct{}

func (r *RealTimeProvider) Now() time.Time {
	return time.Now()
}

func (r *RealTimeProvider) Month() time.Month {
	return time.Now().Month()
}

func (r *RealTimeProvider) Day() int {
	return time.Now().Day()
}
