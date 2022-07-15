package fake

import (
	"log"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/syronz/fake2db/pkg"
)

var Patterns map[string]Faker

type Bridge struct{}

func init() {
	Patterns = make(map[string]Faker)
	bridge := Bridge{}

	faker := gofakeit.New(0)
	gofakeit.SetGlobalFaker(faker)

	Patterns["FIRST_NAME"] = bridge.PersonFirstName
	Patterns["LAST_NAME"] = bridge.PersonLastName
	Patterns["NUMBER"] = bridge.Number
	Patterns["RANDOM_STRING"] = bridge.RandomString
	Patterns["DATE_RANGE"] = bridge.DateRange
	Patterns["DATE"] = bridge.Date
	Patterns["CITY"] = bridge.City
	Patterns["ZIP"] = bridge.Zip
	Patterns["STREET_NAME"] = bridge.StreetName
	Patterns["STREET_NUMBER"] = bridge.StreetNumber
	//Patterns["DATETIME"] = bridge.RandomDateTime
}

func (b *Bridge) PersonFirstName(args ...string) string {
	return gofakeit.FirstName()
}

func (b *Bridge) PersonLastName(args ...string) string {
	return gofakeit.LastName()
}

func (b *Bridge) Number(args ...string) string {
	var min, max int

	var err error
	if len(args) == 2 {
		if min, err = strconv.Atoi(args[0]); err != nil {
			log.Fatalln("first argument in __NUMBER__ is wrong", err)
		}
		if max, err = strconv.Atoi(args[1]); err != nil {
			log.Fatalln("second argument in __NUMBER__ is wrong", err)
		}
	} else {
		log.Fatalln("__NUMBER__ needs two argument")
	}

	return strconv.Itoa(gofakeit.Number(min, max))
}

func (b *Bridge) RandomString(args ...string) string {
	return gofakeit.RandomString(args)
}

func dateGenerator(args ...string) time.Time {
	var start, end time.Time

	var err error
	if len(args) == 2 {
		if start, err = time.Parse(pkg.DateLayout, args[0]); err != nil {
			log.Fatalln("first argument in __DATE_RANGE__ is wrong", err)
		}
		if end, err = time.Parse(pkg.DateLayout, args[1]); err != nil {
			log.Fatalln("second argument in __DATE_RANGE__ is wrong", err)
		}
	} else {
		log.Fatalln("__DATE_RANGE__ needs two argument")
	}

	return gofakeit.DateRange(start, end)
}

func (b *Bridge) DateRange(args ...string) string {
	d := dateGenerator(args...)

	return d.Format("2006-01-02")
}

func (b *Bridge) Date(args ...string) string {
	return gofakeit.Date().Format("2006-01-02 15:04:05")
}

func (b *Bridge) City(args ...string) string {
	return gofakeit.City()
}

func (b *Bridge) Zip(args ...string) string {
	return gofakeit.Zip()
}

func (b *Bridge) StreetName(args ...string) string {
	return gofakeit.StreetName()
}

func (b *Bridge) StreetNumber(args ...string) string {
	return gofakeit.StreetNumber()
}

//
//func (b *Bridge) Address(args ...string) string {
//	return b.internalFaker.Address()
//}
