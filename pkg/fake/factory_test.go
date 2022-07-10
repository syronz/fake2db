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
			message: "one-line-comma-inside-value",
			in: `VALUES(
				'__PERSON_FIRST_NAME__ __PERSON_LAST_NAME__',
				'__ENUM(male,female,other)__',
				'__INT(100,999)__',
				'__DATE(1950-01-01,2010-01-01)__' ,
				'__ADDRESS__',
				'__DATETIME(2020-01-01 00:00:00,2022-07-10 23:59:59)__')`,
			out: "(?,?,?,?,?,?)",
		},
	}

	for _, v := range samples {
		r := ValueClauseToQuestionMark(v.in)
		assert.Equal(t, v.out, r, v.message)
	}
}
