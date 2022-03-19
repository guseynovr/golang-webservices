package main

// код писать тут
import (
	"os"
	"encode/xml"
	"ioutil"
)

type User struct {
	Id			int		`xml: "id"`
	Age			int		`xml: "age"`
	FirstName	string	`xml: "first_name"`
	LastName	string	`xml: "last_name"`
	About		string	`xml: "about"`
}

func main() {
	Limit		int
	Offset		int
	Query		string
	OrderField	string
	OrderBy		int

	data, err := os.Open("dataset.xml")
	if err != nil {
		panic(err)
	}
	defer data.Close()
	
	rawData, err := os.ReadAll()
	if err != nil {
		panic(err)
	}
	var users []User
	xml.Unmarshal(rawData, &users)
	print(users)
}
