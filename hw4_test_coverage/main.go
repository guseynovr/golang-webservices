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
		Gender    string `xml:"gender"`
		About     string `xml:"about"`
	} `xml:"row"`
}

func main() {
	limit := 10
	offset := 0
	query := "culpa"
	orderField := "id"
	orderBy := -1

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
	i := 0
	for _, v := range users {
		if strings.Contains(v.About, query) ||
			strings.Contains(v.FirstName+v.LastName, query) {
			users[i] = v
			i++
		}
	}
	users = users[:i]

	if offset >= limit || offset > len(users) {
		offset = 0
		limit = 0
	} else if limit > len(users) {
		limit = len(users)
	}
	users = users[offset:limit]

	var less, more func(int, int) bool
	switch strings.ToLower(orderField) {
	case "id":
		less = func(i, j int) bool { return users[i].ID < users[j].ID }
		more = func(i, j int) bool { return users[i].ID > users[j].ID }
	case "age":
		less = func(i, j int) bool { return users[i].Age < users[j].Age }
		more = func(i, j int) bool { return users[i].Age > users[j].Age }
	case "name", "":
		less = func(i, j int) bool {
			return users[i].FirstName+users[i].LastName < users[j].FirstName+users[j].LastName
		}
		more = func(i, j int) bool {
			return users[i].FirstName+users[i].LastName > users[j].FirstName+users[j].LastName
		}
	default:
		panic(fmt.Errorf("SearchServer: invalid OrderField: %s", orderField))
	}
	if orderBy > 0 {
		sort.Slice(users, less)
	} else if orderBy < 0 {
		sort.Slice(users, more)
	}
	// fmt.Printf("%#v\n", users)
	for _, v := range users {
		// fmt.Printf("%#v\n", v.LastName)
		println(v.ID, v.FirstName, v.LastName)
	}
}
