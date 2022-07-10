package fake

import (
	"fmt"
	"regexp"
	"strings"
)

type Factory func() []string
type Faker func(args ...string) string
type FuncList struct {
	Method Faker
	Args   []string
}

// ValueClauseToQuestionMark convert ('message',23,'2022-11-15') to (?,?,?)
func ValueClauseToQuestionMark(valueClause string) string {
	if len(valueClause) == 0 {
		return ""
	}

	var count int
	//var inVariable rune
	//runeArr := []rune(valueClause)
	//for i := 0; i < len(runeArr); i++ {
	//	v := runeArr[i]
	//
	//	fmt.Printf("i: %v, v: %q, count: %v, inVariable: %v\n", i, v, count, inVariable)
	//
	//	switch v {
	//	case '"':
	//		fallthrough
	//	case '\'':
	//		if inVariable == 0 {
	//			inVariable = v
	//		} else if inVariable == v {
	//			inVariable = 0
	//		}
	//
	//	case '\\':
	//		// skip two characters
	//		i++
	//		continue
	//	}
	//
	//	if inVariable == 0 && v == ',' {
	//		count++
	//	}
	//}

	re := regexp.MustCompile(`'\s*,\s*'`)
	matches := re.FindAllString(valueClause, -1)
	count = len(matches) + 1

	return "(" + strings.Repeat("?,", count)[0:count*2-1] + ")"
}

func NewFactory(valueClause string) (Factory, error) {

	funcList := []FuncList{
		{
			Method: PersonFirstName,
		},
		{
			Method: PersonFirstName,
		},
	}

	fmt.Println(funcList)

	return func() []string {
		var rands []string
		for _, v := range funcList {
			rands = append(rands, v.Method(v.Args...))
		}
		return rands
	}, nil
}

func PersonFirstName(args ...string) string {
	return "diako"
}
