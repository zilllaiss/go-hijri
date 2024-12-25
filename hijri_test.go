package hijri_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/zilllaiss/go-hijri"
)

var hijriTestData []TestData

func init() {
	var err error
	hijriTestData, err = generateTestData("test/hijri.csv")
	if err != nil {
		panic(err)
	}
}

func Test_Hijri_ConvertDate(t *testing.T) {
	if len(hijriTestData) == 0 {
		t.Fatal("no tests available for Hijri")
	}

	for _, data := range hijriTestData {
		gregorianDate, _ := time.Parse("2006-01-02", data.Gregorian)
		hijriDate, _ := hijri.CreateHijriDate(gregorianDate, hijri.Default)
		strHijriDate := fmt.Sprintf("%04d-%02d-%02d",
			hijriDate.Year,
			hijriDate.Month,
			hijriDate.Day)

		if strHijriDate != data.Hijri {
			t.Errorf("%s: want %s got %s\n", data.Gregorian, data.Hijri, strHijriDate)
		}
	}
}

func TestNewHijri(t *testing.T) {
	// Must be identical
	for _, data := range hijriTestData {
		h := strings.Split(data.Hijri, "-")
		year, _ := strconv.ParseInt(h[0], 10, 64)
		month, _ := strconv.ParseInt(h[1], 10, 64)
		day, _ := strconv.ParseInt(h[2], 10, 64)

		hijriDate, err := hijri.NewHijriDate(year, month, day, hijri.Default)
		if err != nil {
			t.Error(err)
		}

		strHijriDate := fmt.Sprintf("%04d-%02d-%02d",
			hijriDate.Year,
			hijriDate.Month,
			hijriDate.Day)

		if strHijriDate != data.Hijri {
			t.Errorf("%s: want %s got %s\n", data.Gregorian, data.Hijri, strHijriDate)
		}
	}

	type hd struct {
		year, month, day int64
	}

	// Must not return error
	mustNotErr := []hd{
		// leap years
		{
			year:  1442,
			month: 12,
			day:   30,
		},
		{
			year:  1445,
			month: 12,
			day:   30,
		},
		{
			year:  1447,
			month: 12,
			day:   30,
		},
		{
			year:  1450,
			month: 12,
			day:   30,
		},
		{
			year:  1453,
			month: 12,
			day:   30,
		},
		{
			year:  1456,
			month: 12,
			day:   30,
		},
		{
			year:  1458,
			month: 12,
			day:   30,
		},
		{
			year:  1461,
			month: 12,
			day:   30,
		},
		{
			year:  1464,
			month: 12,
			day:   30,
		},
		{
			year:  1466,
			month: 12,
			day:   30,
		},
		{
			year:  1469,
			month: 12,
			day:   30,
		},
	}

	for _, h := range mustNotErr {
		hijridate, err := hijri.NewHijriDate(h.year, h.month, h.day, hijri.Default)
		if err != nil {
			t.Errorf("%#v: want no error, got err: %v\n", hijridate, err.Error())
		}
	}

	// Must return error
	mustErr := []hd{
		// zeroes
		{
			year:  0,
			month: 12,
			day:   12,
		},
		{
			year:  1,
			month: 0,
			day:   12,
		},
		{
			year:  1,
			month: 12,
			day:   0,
		},
		// maxLimits
		{
			year:  2,
			month: 13,
			day:   12,
		},
		{
			year:  2,
			month: 2,
			day:   30,
		},
		{
			year:  1443,
			month: 12,
			day:   30, // non leap years
		},
	}

	for _, h := range mustErr {
		hijridate, err := hijri.NewHijriDate(h.year, h.month, h.day, hijri.Default)
		if err == nil {
			t.Errorf("%#v: want error, got no error instead\n", hijridate)
		}
	}
}

func Test_Hijri_ToGregorian(t *testing.T) {
	if len(hijriTestData) == 0 {
		t.Fatal("no tests available for Hijri")
	}

	for _, data := range hijriTestData {
		var hijriDate hijri.HijriDate
		fmt.Sscanf(data.Hijri, "%d-%d-%d",
			&hijriDate.Year,
			&hijriDate.Month,
			&hijriDate.Day)

		result := hijriDate.ToGregorian().Format("2006-01-02")
		if result != data.Gregorian {
			t.Errorf("%s: want %s got %s\n", data.Hijri, data.Gregorian, result)
		}
	}
}

func Test_Hijri_Bidirectional(t *testing.T) {
	date := time.Date(622, 7, 16, 0, 0, 0, 0, time.UTC)
	for date.Year() <= 2120 {
		// Convert date to hijri
		hijriDate, err := hijri.CreateHijriDate(date, hijri.Default)
		if err != nil {
			date = date.AddDate(0, 0, 1)
			continue
		}

		// Convert back Hijri to Gregorian
		gregorianDate := hijriDate.ToGregorian()

		// Compare original and new gregorian
		strOriginal := date.Format("2006-01-02")
		strGregorian := gregorianDate.Format("2006-01-02")
		strHijri := fmt.Sprintf("%04d-%02d-%02d", hijriDate.Year, hijriDate.Month, hijriDate.Day)

		if strOriginal != strGregorian {
			t.Errorf("Original %s: Hijri %s, Gregorian %s\n",
				strOriginal, strHijri, strGregorian)
		}

		// Increase date
		date = date.AddDate(0, 0, 1)
	}
}
