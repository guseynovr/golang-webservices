package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	json "encoding/json"

	easyjson "github.com/mailru/easyjson"

	jlexer "github.com/mailru/easyjson/jlexer"

	jwriter "github.com/mailru/easyjson/jwriter"
)

//easyjson:json
type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"-"`
	Country  string   `json:"-"`
	Email    string   `json:"email"`
	Job      string   `json:"-"`
	Name     string   `json:"name"`
	Phone    string   `json:"-"`
}

// const filePath string = "./data/users.txt"

// func main() {
// 	FastSearch(os.Stdout)
// }

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(out, "found users:")
	
	seenBrowsers := make(map[string]bool)
	uniqueBrowsers := 0
	reader := bufio.NewReader(file)

	line, _, err := reader.ReadLine()
	for i := 0; err == nil && len(line) > 0; i++ {
		user := User{}
		err := easyjson.Unmarshal(line, &user)
		if err != nil {
			fmt.Println("Err line:", line)
			panic(err)
		}

		isAndroid := false
		isMSIE := false
		
		browsers := user.Browsers

		for _, browser := range browsers {
			thisAndroid := strings.Contains(browser, "Android")
			thisMSIE := strings.Contains(browser, "MSIE")
			if thisAndroid {
				isAndroid = true
			}
			if thisMSIE {
				isMSIE = true
			}
			if (thisAndroid || thisMSIE) && !seenBrowsers[browser] {
				seenBrowsers[browser] = true
				uniqueBrowsers++	
			}
		}

		if isAndroid && isMSIE {
			// log.Println("Android and MSIE user:", user["name"], user["email"])
			email := strings.Replace(user.Email, "@", " [at] ", 1)
			fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
		}

		line, _, err = reader.ReadLine()
	}
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeModels(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeModels(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix[1:])
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeModels(l, v)
}
