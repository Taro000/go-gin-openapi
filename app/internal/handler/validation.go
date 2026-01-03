package handler

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"go-gin-webapi/schemas"
)

const todoDueDatetimeLayout = "2006/01/02 15:04" // yyyy/mm/dd hh:mm

func runeLen(s string) int { return utf8.RuneCountInString(s) }

func validateMaxRunes(value, field string, max int) error {
	if runeLen(value) > max {
		return errors.New(field + " must be <= " + strconvItoa(max) + " chars")
	}
	return nil
}

// strconvItoa: tiny helper to avoid pulling fmt for simple integer conversion.
func strconvItoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var b [32]byte
	n := 0
	for i > 0 {
		b[n] = byte('0' + (i % 10))
		n++
		i /= 10
	}
	out := make([]byte, 0, n+1)
	if neg {
		out = append(out, '-')
	}
	for j := n - 1; j >= 0; j-- {
		out = append(out, b[j])
	}
	return string(out)
}

func parseTodoDueDatetime(in schemas.TodoDueDatetime) (time.Time, error) {
	s := strings.TrimSpace(string(in))
	if s == "" {
		return time.Time{}, errors.New("due_datetime is empty")
	}
	// Store without timezone; interpret as local time.
	t, err := time.ParseInLocation(todoDueDatetimeLayout, s, time.Local)
	if err != nil {
		return time.Time{}, errors.New("due_datetime must be yyyy/mm/dd hh:mm")
	}
	// Ensure canonical zero-padded form (reject e.g. 2026/1/3 9:3).
	if t.Format(todoDueDatetimeLayout) != s {
		return time.Time{}, errors.New("due_datetime must be yyyy/mm/dd hh:mm")
	}
	return t, nil
}

func formatTodoDueDatetime(t time.Time) schemas.TodoDueDatetime {
	return schemas.TodoDueDatetime(t.Format(todoDueDatetimeLayout))
}

var todoStatusToDBCode = map[schemas.TodoStatus]string{
	schemas.TodoStatus("未着手"): "00",
	schemas.TodoStatus("進行中"): "01",
	schemas.TodoStatus("完了"):  "02",
	schemas.TodoStatus("保留"):  "03",
}

var todoDBCodeToStatus = map[string]schemas.TodoStatus{
	"00": schemas.TodoStatus("未着手"),
	"01": schemas.TodoStatus("進行中"),
	"02": schemas.TodoStatus("完了"),
	"03": schemas.TodoStatus("保留"),
}

func todoStatusToCode(s schemas.TodoStatus) (string, bool) {
	code, ok := todoStatusToDBCode[s]
	return code, ok
}

func todoCodeToStatus(code string) (schemas.TodoStatus, bool) {
	s, ok := todoDBCodeToStatus[code]
	return s, ok
}

