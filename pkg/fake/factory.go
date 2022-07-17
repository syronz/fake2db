package fake

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type RandomFunc func() string
type Factory func() []string
type Faker func(args ...string) string
type FuncList struct {
	Method Faker
	Args   []string
}

// ValueClauseToQuestionMark parse config.query.values_clause for generating list of question marks for instance it
// converts ('message','23','2022-11-15') to (?,?,?)
// The output will be used inside statements, for bulk insertion.
func ValueClauseToQuestionMark(valueClause string) (string, int) {
	if len(valueClause) == 0 {
		return "", 0
	}

	re := regexp.MustCompile(`'\s*,\s*'`)
	matches := re.FindAllString(valueClause, -1)
	count := len(matches) + 1

	return "(" + strings.Repeat("?,", count)[0:count*2-1] + ")", count
}

// NewFactory parse config.query.values_clause to find columns and send proper data for each column for extracting
// mapped function
func NewFactory(valueClause string) (Factory, error) {
	re := regexp.MustCompile(`'[0-9A-Za-z_(),-:\s]*'`)
	matches := re.FindAllString(valueClause, -1)

	var funcList []RandomFunc

	for _, v := range matches {
		fn, err := NewFaker(v)
		if err != nil {
			log.Println("error in creating new faker for ", v, ". error is ", err)
			continue
		}
		funcList = append(funcList, fn)
	}

	return func() []string {
		var result []string

		for _, fn := range funcList {
			result = append(result, fn())
		}

		return result
	}, nil
}

type RandomPattern struct {
	RandomFunc string
	Args       []string
}

func extractRandomPatterns(pattern string) ([]RandomPattern, string, error) {
	format := pattern

	var matches []string
	var match string
	var isMatch bool
	str := []rune(pattern)
	for i := 0; i < len(str)-1; i++ {
		if str[i] == '_' && str[i+1] == '_' {
			isMatch = !isMatch
			if !isMatch {
				matches = append(matches, match+"__")
				match = ""
			}
		}
		if isMatch {
			match += string(str[i])
		}
	}

	for _, v := range matches {
		format = strings.Replace(format, v, "%v", 1)
	}

	fakePatterns := make([]RandomPattern, len(matches))

	for i, v := range matches {
		startArgument := strings.Index(v, "(")
		endArgument := strings.Index(v, ")")
		if startArgument+endArgument > 0 {
			fakePatterns[i].RandomFunc = v[2:startArgument]
			fakePatterns[i].Args = strings.Split(v[startArgument+1:endArgument], ",")
			for j := range fakePatterns[i].Args {
				fakePatterns[i].Args[j] = strings.TrimSpace(fakePatterns[i].Args[j])
			}
		} else {
			fakePatterns[i].RandomFunc = v[2 : len(v)-2]
			fakePatterns[i].Args = []string{}
		}
	}

	return fakePatterns, format, nil
}

// NewFaker find requested RandomFunc based on registered keywords that start with __, for instance it maps
// FirstName to __FIRST_NAME__ for generating random person's first name
func NewFaker(phrase string) (RandomFunc, error) {

	randomPatterns, format, _ := extractRandomPatterns(phrase)

	return func() string {
		var randoms []interface{}

		for _, v := range randomPatterns {
			fn, ok := Patterns[v.RandomFunc]
			if !ok {
				log.Fatalln("pattern not exist for", v.RandomFunc, "please update config file")
			}
			randoms = append(randoms, fn(v.Args...))
		}

		return strings.Trim(fmt.Sprintf(format, randoms...), " '")
	}, nil
}
