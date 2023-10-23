package types

import (
	"strings"
	"time"
)

type TimeRangeArray string

func (t *TimeRangeArray) StartAt() *Time {
	arr := strings.Split(string(*t), ",")
	if len(arr) < 1 {
		t, _ := time.Parse("2006-01-02 15:04:05", "0000-00-00 00:00:00")
		t1 := Time(t) //返回最小时间
		return &t1
	}
	ts := Time{}
	err := ts.UnmarshalJSON([]byte(arr[0]))
	if err != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", "0000-00-00 00:00:00")
		t1 := Time(t) //返回最小时间
		return &t1
	}
	return &ts
}

func (t *TimeRangeArray) EndAt() *Time {
	arr := strings.Split(string(*t), ",")
	if len(arr) < 2 {
		t, _ := time.Parse("2006-01-02 15:04:05", "9999-12-31 23:59:59")
		t1 := Time(t) //返回最大时间
		return &t1
	}
	ts := Time{}
	err := ts.UnmarshalJSON([]byte(arr[1]))
	if err != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", "9999-12-31 23:59:59")
		t1 := Time(t) //返回最大时间
		return &t1
	}
	return &ts
}
