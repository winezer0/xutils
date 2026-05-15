package timeutils

import "time"

// TimeObjAddDays 增加天数 (n 可以是负数表示减少)
func TimeObjAddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// TimeObjAddMonths 增加月份
func TimeObjAddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}
