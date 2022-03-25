package main

// код писать тут
import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type TestCase struct {
	client		SearchClient
	Request     SearchRequest
	Result      *SearchResponse
	Error		error
}

func TestFindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	tc := []TestCase{
		{
			client: SearchClient{
				URL: ts.URL,
				AccessToken: "secret",
			},
			Request: SearchRequest{
				Limit:      50,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "name",
				OrderBy:    1,
			},
			Result: &SearchResponse{
				Users: []User{{
					Id:     0,
					Name:   "Boyd Wolf",
					Age:    22,
					About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
					Gender: "male",
				}},
			},
			Error: nil,
		},
		{
			client: SearchClient{
				URL: ts.URL,
				AccessToken: "secret",
			},
			Request: SearchRequest{
				Limit:      -5,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "name",
				OrderBy:    1,
			},
			Result:  nil,
			Error: fmt.Errorf("limit must be > 0"),
		},
		{
			client: SearchClient{
				URL: ts.URL,
				AccessToken: "secret",
			},
			Request: SearchRequest{
				Limit:      50,
				Offset:     -1,
				Query:      "Boyd",
				OrderField: "name",
				OrderBy:    1,
			},
			Result:  nil,
			Error: fmt.Errorf("offset must be > 0"),
		},
		{
			client: SearchClient{
				URL: ts.URL,
				AccessToken: "",
			},
			Request: SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "name",
				OrderBy:    1,
			},
			Result:  nil,
			Error: fmt.Errorf("Bad AccessToken"),
		},
		{
			client: SearchClient{
				URL: ts.URL,
				AccessToken: "secret",
			},
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "culpa",
				OrderField: "id",
				OrderBy:    1,
			},
			Result: &SearchResponse{
				Users: []User{{
					Id:     0,
					Name:   "Boyd Wolf",
					Age:    22,
					About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
					Gender: "male",
				}},
				NextPage: true,
			},
			Error: nil,
		},
		{
			client: SearchClient{
				URL: "",
				AccessToken: "",
			},
			Request: SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "name",
				OrderBy:    1,
			},
			Result:  nil,
			Error: fmt.Errorf("unknown error Get \"?limit=21&offset=0&order_by=1&order_field=name&query=Boyd\": unsupported protocol scheme \"\""),
		},
		{
			client: SearchClient{
				URL: ts.URL,
				AccessToken: "secret",
			},
			Request: SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "asd",
				OrderBy:    1,
			},
			Result:  nil,
			Error: fmt.Errorf("OrderFeld asd invalid"),
		},
		{
			client: SearchClient{
				URL: ts.URL,
				AccessToken: "secret",
			},
			Request: SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "",
				OrderBy:    10,
			},
			Result:  nil,
			Error: fmt.Errorf("unknown bad request error: ErrorBadOrderBy"),
		},
	}

	for testNum, testCase := range tc {
		result, err := testCase.client.FindUsers(testCase.Request)
		if err != nil && testCase.Error == nil {
			t.Errorf("unexpected error: %#v", err)
		}
		// if err == nil && testCase.Error != err {
		if testCase.Error != nil && err.Error() != testCase.Error.Error() {
			t.Errorf("expected %#v, got %#v", testCase.Error, err)
		}
		if !reflect.DeepEqual(result, testCase.Result) {
			t.Errorf("%d: incorrect result: expected: %#v\ngot: %#v\n", testNum, testCase.Result, result)
		}
	}
}

//Server errors
const (
	timeout = iota
	fatal
	unpackError
	unpackResult
)

var serverError int

func TestErrors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerErrors))
	defer ts.Close()

	tc := []int{timeout, fatal, unpackError, unpackResult}
	
	sc := SearchClient{AccessToken: "secret", URL: ts.URL}
	for _, srvError := range tc {
		serverError = srvError
		_, err := sc.FindUsers(SearchRequest{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	}
}

func SearchServerErrors(w http.ResponseWriter, r *http.Request) {
	switch serverError {
	case timeout:
		time.Sleep(6 * time.Second)
	case fatal:
		http.Error(w, "fatal", http.StatusInternalServerError)
	case unpackError:
		http.Error(w, "unpackError", http.StatusBadRequest)
	case unpackResult:
		w.Write([]byte("unpackResult"))
	}
}

type UserData struct {
	Id        int    `xml:"id"`
	Name      string `xml:"-"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name" json:"-"`
	LastName  string `xml:"last_name" json:"-"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

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
		http.Error(w, string(responseError(err.Error())), http.StatusBadRequest)
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

func responseError(e string) []byte {
	errResp := SearchErrorResponse {
		Error: e,
	}
	errJSON, err := json.Marshal(errResp)
	if err != nil {
		panic(err)
	}
	return errJSON
}
