package fake

import "testing"
import "github.com/stretchr/testify/assert"

func TestValueClauseToQuestionMark(t *testing.T) {
	samples := []struct {
		message string
		in      string
		out     string
	}{
		{
			message: "one-line-simple",
			in:      "('name','23','2022-11-15')",
			out:     "(?,?,?)",
		},
		{
			message: "one-line-comma-inside-value",
			in:      "('\"name,lastname','23','2022-11-15')",
			out:     "(?,?,?)",
		},
		{
			message: "one-line-comma-inside-value",
			in:      "",
			out:     "",
		},
		{
			message: "multi-line-real-scenario",
			in: `VALUES(
				'__FIRST_NAME__ __LAST_NAME__',
				'__RANDOM_STRING(male,female,other)__',
				'__NUMBER(100, 999)__',
				'__DATE(1950-01-01,2010-01-01)__' ,
				'__ADDRESS__',
				'__DATETIME(2020-01-01 00:00:00,2022-07-10 23:59:59)__')`,
			out: "(?,?,?,?,?,?)",
		},
	}

	for _, v := range samples {
		r, _ := ValueClauseToQuestionMark(v.in)
		assert.Equal(t, v.out, r, v.message)
	}
}

func TestExtractRandomPatterns(t *testing.T) {
	samples := []struct {
		message string
		in      string
		out     []RandomPattern
		format  string
	}{
		{
			message: "simple-two-pattern",
			in:      "__FIRST_NAME__ __LAST_NAME__",
			out: []RandomPattern{
				{
					RandomFunc: "FIRST_NAME",
					Args:       []string{},
				},
				{
					RandomFunc: "LAST_NAME",
					Args:       []string{},
				},
			},
			format: "%v %v",
		},
		{
			message: "simple-two-pattern-with-argument",
			in:      "__ENUM(male,female,other)__ / __RANDOM_NUMBER(1, 100)__",
			out: []RandomPattern{
				{
					RandomFunc: "ENUM",
					Args:       []string{"male", "female", "other"},
				},
				{
					RandomFunc: "RANDOM_NUMBER",
					Args:       []string{"1", "100"},
				},
			},
			format: "%v / %v",
		},
		{
			message: "address, contain city,zipcode street name and street number",
			in:      "__CITY__ __ZIP__, __STREET_NAME__ __RANDOM_NUMBER(20, 90)__",
			out: []RandomPattern{
				{
					RandomFunc: "CITY",
					Args:       []string{},
				},
				{
					RandomFunc: "ZIP",
					Args:       []string{},
				},
				{
					RandomFunc: "STREET_NAME",
					Args:       []string{},
				},
				{
					RandomFunc: "RANDOM_NUMBER",
					Args:       []string{"20", "90"},
				},
			},
			format: "%v %v, %v %v",
		},
	}

	for _, v := range samples {
		result, format, _ := extractRandomPatterns(v.in)
		assert.Equal(t, v.out, result, v.message)
		assert.Equal(t, v.format, format, v.message)
	}
}
