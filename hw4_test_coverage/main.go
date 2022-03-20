package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type UsersData struct {
	Users []struct {
		ID        int    `xml:"id"`
		Age       int    `xml:"age"`
		FirstName string `xml:"first_name"`
		LastName  string `xml:"last_name"`
		About     string `xml:"about"`
	} `xml:"row"`
}

func main() {
	// Limit := 50
	// Offset := 0
	// Query := "Boyd"
	OrderField := "Name"
	// OrderBy := 0

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

	// i := 0
	// for _, v := range users {

	// }

	var less func(int, int) bool
	switch strings.ToLower(OrderField) {
	case "id":
		less = func(i, j int) bool { return users[i].ID < users[j].ID }
	case "age":
		less = func(i, j int) bool { return users[i].Age < users[j].Age }
	case "name", "":
		less = func(i, j int) bool {
			return users[i].FirstName+users[i].LastName < users[j].FirstName+users[j].LastName
		}
	default:
		panic(fmt.Errorf("SearchServer: invalid OrderField: %s", OrderField))
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].FirstName+users[i].LastName < users[j].FirstName+users[j].LastName
	})
	sort.Slice(users, less)
	// fmt.Printf("%#v\n", users)
	for _, v := range users {
		// fmt.Printf("%#v\n", v.LastName)
		println(v.ID, v.FirstName, v.LastName)
	}
}
