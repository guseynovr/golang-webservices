package main

// код писать тут
import (
	// "encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestCase struct {
	AccessToken string
	Request     SearchRequest
	Result      *SearchResponse
	IsError     bool
}

func TestFindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	tc := []TestCase{
		TestCase{
			AccessToken: "secret",
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
				NextPage: false,
			},
			IsError: false,
		},
		TestCase{
			AccessToken: "secret",
			Request: SearchRequest{
				Limit:      -5,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "name",
				OrderBy:    1,
			},
			Result:  nil,
			IsError: true,
		},
		TestCase{
			AccessToken: "secret",
			Request: SearchRequest{
				Limit:      50,
				Offset:     -1,
				Query:      "Boyd",
				OrderField: "name",
				OrderBy:    1,
			},
			Result:  nil,
			IsError: true,
		},
		TestCase{
			AccessToken: "",
			Request: SearchRequest{
				Limit:      20,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "name",
				OrderBy:    1,
			},
			Result:  nil,
			IsError: true,
		},
	}

	sc := SearchClient{URL: ts.URL}
	for testNum, testCase := range tc {
		sc.AccessToken = testCase.AccessToken
		result, err := sc.FindUsers(testCase.Request)
		if err != nil && !testCase.IsError {
			t.Errorf("unexpected error: %#v", err)
		}
		if err == nil && testCase.IsError {
			t.Error("expected error, got nil")
		}
		// got, _ := json.Marshal(result)
		// expected, _ := json.Marshal(testCase.Result)
		// if string(got) != string(expected) {
		if !reflect.DeepEqual(result, testCase.Result) {
			t.Errorf("%d: incorrect result: expected: %#v\ngot: %#v\n", testNum, testCase.Result, result)
		}
	}
}
