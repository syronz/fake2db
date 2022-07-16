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
	//faker := gofakeit.NewUnlocked(0)
	gofakeit.SetGlobalFaker(faker)

	Patterns["FIRST_NAME"] = bridge.FirstName
	Patterns["LAST_NAME"] = bridge.LastName
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

// FirstName return random first name
func (b *Bridge) FirstName(args ...string) string {
	return gofakeit.FirstName()
}

// LastName return random first name
func (b *Bridge) LastName(args ...string) string {
	return gofakeit.LastName()
}

// Number generate a random number between selected range
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

// RandomString returns random text that could be used inside descriptions and details
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

// DateRange returns a date based on two date as arguments
func (b *Bridge) DateRange(args ...string) string {
	d := dateGenerator(args...)

	return d.Format("2006-01-02")
}

// Date generate random date-time
func (b *Bridge) Date(args ...string) string {
	return gofakeit.Date().Format("2006-01-02 15:04:05")
}

// City returns random city name
func (b *Bridge) City(args ...string) string {
	return gofakeit.City()
}

// Zip generate random zip number, could be used inside address
func (b *Bridge) Zip(args ...string) string {
	return gofakeit.Zip()
}

// StreetName returns random street name, could be combined with Zip and StreetNumber
func (b *Bridge) StreetName(args ...string) string {
	return gofakeit.StreetName()
}

// StreetNumber returns random street number, could be combined with Zip and StreetName
func (b *Bridge) StreetNumber(args ...string) string {
	return gofakeit.StreetNumber()
}

//
//func (b *Bridge) Address(args ...string) string {
//	return b.internalFaker.Address()
//}
