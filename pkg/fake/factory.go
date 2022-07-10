package fake

import "fmt"

type Factory func() []string
type Faker func(args ...string) string
type FuncList struct {
	Method Faker
	Args   []string
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
