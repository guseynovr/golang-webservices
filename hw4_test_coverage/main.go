package main

/* package main

import (
	"encoding/json"
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

type UserData struct {
	Id        int    `xml:"id"`
	Name      string `xml:"-"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name" json:"-"`
	LastName  string `xml:"last_name" json:"-"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

// type UsersData struct {
// 	Users []UserData `xml:"row"`
// }

func SearchServer(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("AccessToken") != "secret" {
		http.Error(w, fmt.Sprint("Bad AccessToken"), http.StatusUnauthorized)
		return
	}
	vals := r.URL.Query()
	limit, err := strconv.Atoi(vals.Get("limit"))
	offset, err := strconv.Atoi(vals.Get("offset"))
	query := vals.Get("query")
	orderField := vals.Get("order_field")
	orderBy, err := strconv.Atoi(vals.Get("order_by"))
	if err != nil {
		http.Error(w, responseError(err.Error()), http.StatusBadRequest)
		return
	}
	data, _ := os.Open("dataset.xml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer data.Close()

	rawData, err := ioutil.ReadAll(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var u struct {
		Users []UserData `xml:"row"`
	}
	err = xml.Unmarshal(rawData, &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users := u.Users
	i := 0
	for _, v := range users {
		if strings.Contains(v.About, query) ||
			strings.Contains(v.FirstName+v.LastName, query) {
			users[i] = v
			users[i].Name = v.FirstName + " " + v.LastName
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

	err = sortByParams(users, orderField, orderBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", result)
}

func sortByParams(users []UserData, orderField string, orderBy int) error {
	if orderBy == 0 {
		return nil
	}
	if orderBy < -1 || orderBy > 1 {
		return fmt.Errorf("%s", responseError("ErrorBadOrderBy"))
	}
	var less, more func(int, int) bool
	switch strings.ToLower(orderField) {
	case "id":
		less = func(i, j int) bool { return users[i].Id < users[j].Id }
		more = func(i, j int) bool { return users[i].Id > users[j].Id }
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
		return fmt.Errorf("%s", responseError("ErrorBadOrderField"))
	}
	if orderBy > 0 {
		sort.Slice(users, less)
	} else if orderBy < 0 {
		sort.Slice(users, more)
	}
	return nil
}

func responseError(e string) string {
	errResp := SearchErrorResponse {
		Error: e,
	}
	errJSON, err := json.Marshal(errResp)
	if err != nil {
		panic(err)
	}
	return string(errJSON)
}

func main() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", SearchServer)
	err := http.ListenAndServe(":5000", nil)
	log.Print(err)
}
*/
