package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type NullMarshalInterface interface {
	MarshalJSON() ([]byte, error)
}

type TimeRange struct {
	StartAt *Time `form:"startAt" json:"startAt"`
	EndAt   *Time `form:"endAt" json:"endAt"`
}

func (v TimeRange) MarshalJSON() ([]byte, error) {
	s := "\"" + time.Time(*v.StartAt).Format(TimeFormat) + "\""
	e := "\"" + time.Time(*v.StartAt).Format(TimeFormat) + "\""
	formatted := s + ":" + e
	return []byte(formatted), nil

}

func (v *TimeRange) UnmarshalJSON(data []byte) error {
	arr := strings.Split(string(data), ",")
	if len(arr) < 2 {
		return fmt.Errorf("expect a slice with 2 element,but got %v", len(arr))
	}
	err := v.StartAt.UnmarshalJSON([]byte(arr[0]))
	if err != nil {
		return err
	}
	err = v.EndAt.UnmarshalJSON([]byte(arr[1]))
	if err != nil {
		return err
	}
	return nil
}

// 因为json 对time.Time解释有bug,这里自定义Time类型序列化和反序列化处理
type Time time.Time

func (v Time) MarshalJSON() ([]byte, error) {
	formatted := "\"" + time.Time(v).Format(TimeFormat) + "\""
	return []byte(formatted), nil

}

var unixTimeStampePattern = regexp.MustCompile("^[0-9]{13}$")

func (v *Time) UnmarshalJSON(data []byte) error {
	//1695733583563
	var t time.Time
	if unixTimeStampePattern.Match(data) {
		timestamp, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		t = time.Unix(timestamp/1000, 0)
	} else {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		var err error
		t, err = time.ParseInLocation(TimeFormat, strings.ReplaceAll(string(data), "\"", ""), loc)
		if err != nil {
			return err
		}
	}

	b := Time(t)
	*v = b
	return nil
}

type NullTime sql.NullTime

func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var s *time.Time

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Time = *s
	} else {
		v.Valid = false
	}

	return nil
}

type NullString sql.NullString

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

type NullInt16 sql.NullInt16

func (v NullInt16) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int16)
	} else {
		return json.Marshal(nil)
	}
}

type NullInt32 sql.NullInt32

func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	} else {
		return json.Marshal(nil)
	}
}

type NullInt64 sql.NullInt64

func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

type NullBool sql.NullBool

func (v NullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	} else {
		return json.Marshal(nil)
	}
}

type NullByte sql.NullByte

func (v NullByte) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Byte)
	} else {
		return json.Marshal(nil)
	}
}

type NullFloat64 sql.NullFloat64

func (v NullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	} else {
		return json.Marshal(nil)
	}
}
