package dsl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strings"

	lua "github.com/yuin/gopher-lua"
)


type dslHttp struct {

	typeOfRequest	string
	url 		string
	basicAuthUser   string
	basicAuthPasswd string
	contentType 	string
	headers 	string
	useragent 	string
	body 		string
	client  	*http.Client

}

func (h *dslHttp) buildRequest() (req *http.Request, err error) {


	switch h.typeOfRequest {
	case "GET":

		req, err = http.NewRequest("GET", h.url, nil);  
		if err != nil {
			return nil, err
		}


	case "POST":

		buf := bytes.NewBuffer([]byte(h.body))

		req, err = http.NewRequest("POST", h.url, buf); 
		req.Header.Set("Content-Type", h.contentType)
		if err != nil {
			return nil, err
		}


	default:
		return nil, fmt.Errorf("[ERROR] failed to create request with unsupported type: %s\n", h.typeOfRequest)

	}

	if h.useragent != "nil" {
		req.Header.Set("User-Agent", h.useragent)
	}

	if h.basicAuthUser != "nil" {
		req.SetBasicAuth(h.basicAuthUser, h.basicAuthPasswd)
	}

	if h.headers != "nil" {
		headers := strings.Split(h.headers, ";")

		for _, header := range headers {
			slice := strings.Split(header, ":")

			req.Header.Set(strings.TrimSpace(slice[0]), strings.TrimSpace(slice[1]))
		}
	}


	return req, nil
}

func (d *dslState) dslHttpRequest(L *lua.LState) int {
	args := L.CheckTable(1)

	url := args.RawGetString("url").String()
	typeOfRequest := args.RawGetString("type").String()
	ua := args.RawGetString("useragent").String()	
	contentType := args.RawGetString("contenttype").String()
	body := args.RawGetString("body").String()
	user := args.RawGetString("user").String()
	password := args.RawGetString("password").String()
	headers := args.RawGetString("headers").String()

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	t := &dslHttp { url: url, typeOfRequest: typeOfRequest, useragent: ua, client: client, contentType: contentType, body: body,
		basicAuthUser: user, basicAuthPasswd: password, headers: headers }

	req, err := t.buildRequest()	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("http create request: %s\n", err.Error())))
		return 2
		}

	response, err := t.doRequest(req)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("http send request error: %s\n", err.Error())))
		return 2
	}
	//fmt.Printf("%+v", response)

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("http read responce error: %s\n", err.Error())))
		return 2
	defer response.Body.Close()
	}
	result := L.NewTable()
	L.SetField(result, "code", lua.LNumber(response.StatusCode))
	L.SetField(result, "body", lua.LString(string(data)))
	L.Push(result)
	return 1

}

func (h *dslHttp) doRequest(req *http.Request)( *http.Response, error)  {
	response, err := h.client.Do(req)
	if err != nil {
		return nil, err 
	}
	return response, nil
}

