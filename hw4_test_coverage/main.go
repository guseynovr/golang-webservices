package main

import (
	"fmt"
	"os"
	"encoding/xml"
	"io/ioutil"
)

// type User struct {
// 	Id			int		`xml: "id"`
// 	Age			int		`xml: "age"`
// 	FirstName	string	`xml: "first_name"`
// 	LastName	string	`xml: "last_name"`
// 	About		string	`xml: "about"`
// }

type User struct {
	// Text          string `xml:",chardata"`
	ID            int `xml:"id"`
	// Guid          string `xml:"-"`
	// IsActive      string `xml:"-"`
	// Balance       string `xml:"-"`
	// Picture       string `xml:"picture"`
	Age           int `xml:"age"`
	// EyeColor      string `xml:"eyeColor"`
	FirstName     string `xml:"first_name"`
	LastName      string `xml:"last_name"`
	// Gender        string `xml:"gender"`
	// Company       string `xml:"company"`
	// Email         string `xml:"email"`
	// Phone         string `xml:"phone"`
	// Address       string `xml:"address"`
	About         string `xml:"about"`
	// Registered    string `xml:"registered"`
	// FavoriteFruit string `xml:"favoriteFruit"`
}

type UsersData struct {
	// XMLName xml.Name `xml:"root"`
	// Text    string   `xml:",chardata"`
	Users	[]User `xml:"row"`
} 

func main() {
	Limit := 50
	Offset := 0
	Query := ""
	OrderField := "id"
	OrderBy := 0

	data, err := os.Open("dataset.xml")
	if err != nil {
		panic(err)
	}
	defer data.Close()
	
	rawData, err := ioutil.ReadAll(data)
	if err != nil {
		panic(err)
	}
	var u UsersData
	err = xml.Unmarshal(rawData, &u)
	if err != nil {
		panic(err)
	}
	users := u.Users
	// fmt.Printf("%#v\n", users)
	// for _, v := range users {
	// 	fmt.Printf("%#v\n", v)
	// }
}
