package kaiheila

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (c *Client) request(method string, version int, path string, values interface{}, v interface{}) (err error) {
	resp, err := c.do(method, version, path, struct2values(values))
	if err != nil {
		return fmt.Errorf("[kaiheila] %s > %s", path, err)
	}
	url := resp.Request.URL.String()
	defer resp.Body.Close()

	// Status check
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[kaiheila] %s > %v", url, resp.StatusCode)
	}

	// Decode msg
	msg := &httpMsg{}
	err = json.NewDecoder(resp.Body).Decode(msg)
	if err != nil {
		return fmt.Errorf("[kaiheila] %s > %s", url, err)
	}
	if msg.Code != 0 {
		return fmt.Errorf("[kaiheila] %s > %d %s", url, msg.Code, msg.Message)
	}
	if v != nil {
		err = json.Unmarshal(msg.Data, v)
		if err != nil {
			return fmt.Errorf("[kaiheila] %s > %s", url, err)
		}
	}
	return nil
}

func (c *Client) do(method string, version int, path string, values url.Values) (resp *http.Response, err error) {
	http.DefaultClient.Timeout = time.Duration(c.Timeout) * time.Second
	// Build request
	var req *http.Request
	var body io.Reader
	var header = http.Header{}
	url := c.URL + "/v" + strconv.Itoa(version) + "/" + path
	if values != nil {
		switch method {
		case "GET":
			url += "?" + values.Encode()
		case "POST":
			body = strings.NewReader(values.Encode())
			header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	header.Add("Authorization", c.TokenType+" "+c.Token)
	req, err = http.NewRequest(method, url, body)
	req.Header = header
	if err != nil {
		return
	}

	// Do request
	return http.DefaultClient.Do(req)
}

func struct2values(v interface{}) (values url.Values) {
	if v == nil {
		return nil
	}
	values = url.Values{}
	iVal := reflect.ValueOf(v).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// tags
		name, options := parseTag(typ.Field(i).Tag.Get("json"))
		if options.Contains("omitempty") && f.IsZero() {
			continue
		}
		// value
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		values.Set(name, v)
	}
	return
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, ""
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
