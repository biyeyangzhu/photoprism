package txt

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var DateRegexp = regexp.MustCompile("\\D\\d{4}[\\-_]\\d{2}[\\-_]\\d{2,}")
var DateCanonicalRegexp = regexp.MustCompile("\\D\\d{8}_\\d{6}_\\w+\\.")
var DatePathRegexp = regexp.MustCompile("\\D\\d{4}/\\d{1,2}/?\\d*")
var DateTimeRegexp = regexp.MustCompile("\\D\\d{4}[\\-_]\\d{2}[\\-_]\\d{2}.{1,4}\\d{2}\\D\\d{2}\\D\\d{2,}")
var DateIntRegexp = regexp.MustCompile("\\d{1,4}")
var YearRegexp = regexp.MustCompile("\\d{4,5}")

var (
	YearMin = 1990
	YearMax = time.Now().Year() + 3
)

const (
	MonthMin = 1
	MonthMax = 12
	DayMin   = 1
	DayMax   = 31
	HourMin  = 0
	HourMax  = 24
	MinMin   = 0
	MinMax   = 59
	SecMin   = 0
	SecMax   = 59
)

// Time returns a string as time or the zero time instant in case it can not be converted.
func Time(s string) (result time.Time) {
	defer func() {
		if err := recover(); err != nil {
			result = time.Time{}
		}
	}()

	if len(s) < 6 {
		return time.Time{}
	}

	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	b := []byte(s)

	if found := DateCanonicalRegexp.Find(b); len(found) > 0 { // Is it a canonical name like "20120727_093920_97425909.jpg"?
		if date, err := time.Parse("20060102_150405", string(found[1:16])); err == nil {
			result = date.Round(time.Second).UTC()
		}
	} else if found := DateTimeRegexp.Find(b); len(found) > 0 { // Is it a date with time like "2020-01-30_09-57-18"?
		n := DateIntRegexp.FindAll(found, -1)

		if len(n) != 6 {
			return result
		}

		year := Int(string(n[0]))
		month := Int(string(n[1]))
		day := Int(string(n[2]))
		hour := Int(string(n[3]))
		min := Int(string(n[4]))
		sec := Int(string(n[5]))

		if year < YearMin || year > YearMax || month < MonthMin || month > MonthMax || day < DayMin || day > DayMax {
			return result
		}

		if hour < HourMin || hour > HourMax || min < MinMin || min > MinMax || sec < SecMin || sec > SecMax {
			return result
		}

		result = time.Date(
			year,
			time.Month(month),
			day,
			hour,
			min,
			sec,
			0,
			time.UTC)

	} else if found := DateRegexp.Find(b); len(found) > 0 { // Is it a date only like "2020-01-30"?
		n := DateIntRegexp.FindAll(found, -1)

		if len(n) != 3 {
			return result
		}

		year := Int(string(n[0]))
		month := Int(string(n[1]))
		day := Int(string(n[2]))

		if year < YearMin || year > YearMax || month < MonthMin || month > MonthMax || day < DayMin || day > DayMax {
			return result
		}

		result = time.Date(
			year,
			time.Month(month),
			day,
			0,
			0,
			0,
			0,
			time.UTC)
	} else if found := DatePathRegexp.Find(b); len(found) > 0 { // Is it a date path like "2020/01/03"?
		n := DateIntRegexp.FindAll(found, -1)

		if len(n) < 2 || len(n) > 3 {
			return result
		}

		year := Int(string(n[0]))
		month := Int(string(n[1]))

		if year < YearMin || year > YearMax || month < MonthMin || month > MonthMax {
			return result
		}

		if len(n) == 2 {
			result = time.Date(
				year,
				time.Month(month),
				1,
				0,
				0,
				0,
				0,
				time.UTC)
		} else if day := Int(string(n[2])); day >= DayMin && day <= DayMax {
			result = time.Date(
				year,
				time.Month(month),
				day,
				0,
				0,
				0,
				0,
				time.UTC)
		}
	}

	return result.UTC()
}

// Int returns a string as int or 0 if it can not be converted.
func Int(s string) int {
	if s == "" {
		return 0
	}

	result, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return 0
	}

	return int(result)
}

// IsUInt returns true if a string only contains an unsigned integer.
func IsUInt(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r < 48 || r > 57 {
			return false
		}
	}

	return true
}

// CountryCode tries to find a matching country code for a given string e.g. from a file oder directory name.
func CountryCode(s string) string {
	if s == "zz" {
		return "zz"
	}

	r := strings.NewReplacer("--", " / ", "_", " ", "-", " ")
	s = r.Replace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "  ", " ")

	for keyword, code := range Countries {
		if strings.Contains(s, keyword) {
			return code
		}
	}

	return "zz"
}

// Year tries to find a matching year for a given string e.g. from a file oder directory name.
func Year(s string) int {
	b := []byte(s)

	found := YearRegexp.FindAll(b, -1)

	for _, match := range found {
		year := Int(string(match))

		if year > YearMin && year < YearMax {
			return year
		}
	}

	return 0
}
