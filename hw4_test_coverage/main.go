package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type UserXML struct {
	ID        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
}

type UsersData struct {
	Users []UserXML `xml:"row"`
	// Users []struct {
	// 	ID        int    `xml:"id"`
	// 	Age       int    `xml:"age"`
	// 	FirstName string `xml:"first_name"`
	// 	LastName  string `xml:"last_name"`
	// 	Gender    string `xml:"gender"`
	// 	About     string `xml:"about"`
	// } `xml:"row"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {

	vals := r.URL.Query()

	// fmt.Printf("%#v\n", vals.Get("limit"))
	limit, _ := strconv.Atoi(vals.Get("limit"))
	println("limit = ", limit)
	offset, _ := strconv.Atoi(vals.Get("offset"))
	query := vals.Get("query")
	orderField := vals.Get("order_field")
	orderBy, _ := strconv.Atoi(vals.Get("order_by"))
	data, _ := os.Open("dataset.xml")
	// if err != nil {
	// 	panic(err)
	// }
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

	sortByParams(users, orderField, orderBy)
	// fmt.Printf("%#v\n", users)
	fmt.Fprintf(w, "%#v", users)
	// for _, v := range users {
	// 	// fmt.Printf("%#v\n", v.LastName)
	// 	println(v.ID, v.FirstName, v.LastName)
	// }
}

func sortByParams(users []UserXML, orderField string, orderBy int) {
	if orderBy == 0 {
		return
	}
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
}

func main() {
	err := http.ListenAndServe(":5000", http.HandlerFunc(SearchServer))
	log.Print(err)
	// sc := SearchClient{
	// 	URL: "http://127.0.0.1:5000",
	// }
	// sr := SearchRequest{
	// 	Limit:      50,
	// 	Offset:     0,
	// 	Query:      "culpa",
	// 	OrderField: "name",
	// 	OrderBy:    1,
	// }
	// sResp, _ := sc.FindUsers(sr)
	// println(sResp.Users)
}
