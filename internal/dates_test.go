package dates_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	dates "github.com/craigpaul/twiligo/internal"
)

type Date struct {
	Time dates.Rfc2822Time `json:"time"`
}

func TestCanUnmarshalProperRfc2822StringSuccessfully(t *testing.T) {
	jsonData := []byte(`{"time":"Thu, 30 Jul 2020 00:00:00 +0000"}`)

	var date Date
	err := json.Unmarshal(jsonData, &date)

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if date.Time.Day() != 30 {
		t.Logf("Incorrect time parsed, expecting [%d], but received [%d]", 30, date.Time.Day())
		t.Fail()
	}

	if date.Time.Weekday() != time.Thursday {
		t.Logf("Incorrect time parsed, expecting [%v], but received [%v]", time.Thursday, date.Time.Weekday())
		t.Fail()
	}
}

func TestWillHandleUnmarshallingIncorrectDateString(t *testing.T) {
	jsonData := []byte(`{"time":"Thu 30 Jul 00:00:00 +0000 2020"}`)

	var date Date
	err := json.Unmarshal(jsonData, &date)

	if err == nil {
		t.Log("Did not receive the expected error")
		t.Fail()
	}
}

func TestCanMarshalToJSONSuccessfully(t *testing.T) {
	jsonData := []byte(`{"time":"Thu, 30 Jul 2020 00:00:00 +0000"}`)

	var date Date
	err := json.Unmarshal(jsonData, &date)

	data, err := json.Marshal(date)

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	expected := "Thu, 30 Jul 2020 00:00:00 +0000"

	if strings.Contains(string(data), expected) == false {
		t.Logf("Incorrect time marshalled, expecting to see [%s] in [%s]", expected, string(data))
		t.Fail()
	}
}
