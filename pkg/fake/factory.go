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

// ValueClauseToQuestionMark convert ('message','23','2022-11-15') to (?,?,?)
func ValueClauseToQuestionMark(valueClause string) string {
	if len(valueClause) == 0 {
		return ""
	}

	re := regexp.MustCompile(`'\s*,\s*'`)
	matches := re.FindAllString(valueClause, -1)
	count := len(matches) + 1

	return "(" + strings.Repeat("?,", count)[0:count*2-1] + ")"
}

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

	randomPatterns := make([]RandomPattern, len(matches))

	for i, v := range matches {
		startArgument := strings.Index(v, "(")
		endArgument := strings.Index(v, ")")
		if startArgument+endArgument > 0 {
			randomPatterns[i].RandomFunc = v[2:startArgument]
			randomPatterns[i].Args = strings.Split(v[startArgument+1:endArgument], ",")
			for j := range randomPatterns[i].Args {
				randomPatterns[i].Args[j] = strings.TrimSpace(randomPatterns[i].Args[j])
			}
		} else {
			randomPatterns[i].RandomFunc = v[2 : len(v)-2]
			randomPatterns[i].Args = []string{}
		}
	}

	return randomPatterns, format, nil
}

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
