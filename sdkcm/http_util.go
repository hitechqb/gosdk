package sdkcm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type SDKHttpUtil interface {
	Get(url string, data interface{}) ([]byte, error)
	DoPost(url string, params url.Values) ([]byte, error)
	DoPostJSON(url string, data interface{}) ([]byte, error)
}

type httpUtil struct {
	Timeout *time.Duration
	Header  http.Header
}

func NewHttpUtil(timeout *time.Duration, header http.Header) *httpUtil {
	return &httpUtil{
		Timeout: timeout,
		Header:  header,
	}
}

func do(req *http.Request, timeout time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}

// Method Get
func (p *httpUtil) Get(url string, data interface{}) ([]byte, error) {
	jsonData, _ := json.Marshal(&data)
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	if p.Timeout == nil {
		t := time.Duration(5)
		p.Timeout = &t
	}

	if p.Header != nil {
		req.Header = p.Header
	}

	return do(req, *p.Timeout)
}

// Method POST with params (form)
func (p *httpUtil) DoPost(url string, params url.Values) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if p.Timeout == nil {
		t := time.Duration(5)
		p.Timeout = &t
	}

	return do(req, *p.Timeout)
}

// Method POST with body json
func (p *httpUtil) DoPostJSON(url string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if p.Header != nil {
		req.Header = p.Header
	}

	if p.Timeout == nil {
		t := time.Duration(5)
		p.Timeout = &t
	}

	return do(req, *p.Timeout)
}

// Map struct to param - useful for method DoPost
func StructToParams(in interface{}) url.Values {
	params := make(url.Values)
	v := reflect.ValueOf(in)

	for i := 0; i < v.NumField(); i++ {
		name := strings.ToLower(v.Type().Field(i).Name)
		value := v.Field(i).Interface()
		if str, ok := value.(string); ok {
			params.Add(name, str)
		}
	}

	return params
}

func ParamsToMap(params url.Values) map[string]interface{} {
	m := make(map[string]interface{})

	for k, v := range params {
		m[k] = v[0]
	}

	return m
}

func MapToParams(m map[string]string) url.Values {
	params := make(url.Values)

	for k, v := range m {
		params.Add(k, v)
	}

	return params
}
