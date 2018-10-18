package dsl

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
	"strconv"
	"regexp"


	lua "github.com/yuin/gopher-lua"
)

type dslZabbix struct {
	url             string
	user            string
	passwd          string
	proxy           string
	basicAuthUser   string
	basicAuthPasswd string
	ignoreSsl       bool
	// internal
	client  *http.Client
	id      int
	auth    string
	version string
}

func (d *dslZabbix) toString() string {
	return fmt.Sprintf("(%s:%s) %s [%s]", d.user, d.passwd, d.url, d.version)
}

type dslZabbixJsonRPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	Id      int         `json:"id"`
}

type dslZabbixJsonRPCResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	Error   dslZabbixError `json:"error"`
	Result  interface{}    `json:"result"`
	Id      int            `json:"id"`
}

type dslZabbixError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type dslZabbixGroup struct {
	GroupID  string `json:"groupid"`
	Name     string `json:"name"`
	Internal string `json:"internal"`
}

type dslZabbixHost struct {
	HostID string `json:"hostid"`
	Name   string `json:"name"`
	Host   string `json:"host"`
}

type dslZabbixItem struct {
	DataType   string      `json:"data_type"`
	HostID     string      `json:"hostid"`
	ItemID     string      `json:"itemid"`
	LastValue  interface{} `json:"lastvalue"`
	LastClock  string      `json:"lastclock"`
	Name       string      `json:"name"`
	TemplateID string      `json:"templateid"`
	Units      string      `json:"units"`
	Type       string      `json:"type"`
	ValueType  string      `json:"value_type"`
}

type dslZabbixTrigger struct {
	TriggerID   string           `json:"triggerid"`
	Description string           `json:"description"`
	Priority    string           `json:"priority"`
	LastChange  string           `json:"lastchange"`
	Groups      []dslZabbixGroup `json:"groups"`
	Hosts       []dslZabbixHost  `json:"hosts"`
	Items       []dslZabbixItem  `json:"items"`
}


func (z *dslZabbixError) Error() string {
	return z.Data
}

func (d *dslZabbix) request(method string, data interface{}) (dslZabbixJsonRPCResponse, error) {

	id := d.id
	d.id = id + 1
	jsonobj := dslZabbixJsonRPCRequest{Jsonrpc: "2.0", Method: method, Params: data, Auth: d.auth, Id: id}
	if method == `APIInfo.version` {
		jsonobj.Auth = ``
	}
	encoded, err := json.Marshal(jsonobj)
	if err != nil {
		log.Printf("[ERROR] process request zabbix[url: `%s`, method: `%s`, id: `%d`]: %s\n", d.url, method, d.id, err.Error())
		return dslZabbixJsonRPCResponse{Error: dslZabbixError{Code: -1, Data: err.Error()}}, err
	}
	log.Printf("[DEBUG] body: %s\n", encoded)
	request, err := http.NewRequest("POST", d.url, bytes.NewBuffer(encoded))
	if err != nil {
		log.Printf("[ERROR] process request zabbix[url: `%s`, method: `%s`, id: `%d`]: %s\n", d.url, method, d.id, err.Error())
		return dslZabbixJsonRPCResponse{Error: dslZabbixError{Code: -1, Data: err.Error()}}, err
	}
	d.updateReq(request)
	if d.client == nil {
		panic("empty client")
	}
	response, err := d.client.Do(request)
	if err != nil {
		log.Printf("[ERROR] process request zabbix[url: `%s`, method: `%s`, id: `%d`]: %s\n", d.url, method, d.id, err.Error())
		return dslZabbixJsonRPCResponse{Error: dslZabbixError{Code: -1, Data: err.Error()}}, err
	}
	defer response.Body.Close()
	var result dslZabbixJsonRPCResponse
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		log.Printf("[ERROR] process request zabbix[url: `%s`, method: `%s`, id: `%d`]: %s\n", d.url, method, d.id, err.Error())
		return dslZabbixJsonRPCResponse{Error: dslZabbixError{Code: -1, Data: err.Error()}}, err
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		log.Printf("[ERROR] process request zabbix[url: `%s`, method: `%s`, id: `%d`]: %s\n", d.url, method, d.id, err.Error())
		return dslZabbixJsonRPCResponse{Error: dslZabbixError{Code: -1, Data: err.Error()}}, err
	}
	log.Printf("[INFO] processed request zabbix[url: `%s`, method: `%s`, id: `%d`]\n", d.url, method, id)
	return result, nil
}

func (d *dslState) dslZabbixLogin(L *lua.LState) int {
	params := make(map[string]string, 0)
	params["user"] = L.CheckString(2)
	params["password"] = L.CheckString(3)

	z := &dslZabbix{ user: params["user"], passwd: params["password"], url: L.CheckString(1)}

	err := z.buildHttpClient()
	if err != nil {
		log.Println(err)
	}

	response, err := z.request("user.login", params)
	if err != nil {
		log.Printf("[ERROR] ", err)
	        L.Push(lua.LNil)
                L.Push(lua.LString(err.Error()))
		return 2
	}
	if response.Error.Code != 0 {
		log.Printf("[ERROR] %+v", &response.Error) 
	        L.Push(lua.LNil)
                L.Push(lua.LString(response.Error.Code))
		return 2
	}
	z.auth = response.Result.(string)

	ud := L.NewUserData()
        ud.Value = z
        L.SetMetatable(ud, L.GetTypeMetatable("zabbix"))
        L.Push(ud)
        log.Printf("[INFO] New zabbix connection to `%s`\n", z.url)
        return 1

}

func (d *dslState) dslZabbixLogout(L *lua.LState) int {
	z := checkZabbixConn(L)
	response, err := z.request("user.logout", make(map[string]string, 0))
	if err != nil {
		return 1
	}
	if response.Error.Code != 0 {
		return 1
	}
	return 0
}

func (d *dslZabbix) apiVersion() (string, error) {
	response, err := d.request("APIInfo.version", make(map[string]string, 0))
	if err != nil {
		return "", err
	}
	if response.Error.Code != 0 {
		return "", &response.Error
	}
	version := response.Result.(string)
	d.version = version
	return version, nil
}

func (d *dslState) dslZabbixGetTriggers(L *lua.LState) int {
	z := checkZabbixConn(L)
	args := L.CheckTable(2)
	pattern := args.RawGetString("pattern").String()
	duration := args.RawGetString("duration").String()
	minSeverity := args.RawGetString("severity").String()

	durationInt, _ := strconv.Atoi(duration)

	triggerUntil := time.Now().Unix() - int64(durationInt * 60)

	params := make(map[string]interface{}, 0)
	params["output"] = "extend"
	params["sortfield"] = "priority"
	params["sortorder"] = "DESC"

	filter := make(map[string]string)
	filter["value"] = "1"
	filter["status"] = "0"
	params["filter"] = filter
	params["min_severity"] = minSeverity
	params["lastChangeTill"] = triggerUntil

	params["expandData"] = "1"
	params["expandDescription"] = "1"
	params["skipDependent"] = "1"
	params["withLastEventUnacknowledged"] = "1"
	params["selectGroups"] = "extend"
	params["selectTriggers"] = "extend"
	params["selectHosts"] = "extend"
	params["selectItems"] = "extend"
	params["limit"] = 200

	response, err := z.request(`trigger.get`, params)
	if err != nil {
		log.Printf("[ERROR] ", err)
	        L.Push(lua.LNil)
                L.Push(lua.LString(err.Error()))
		return 2
	}
	if response.Error.Code != 0 {
		log.Printf("[ERROR] %+v", &response.Error) 
	        L.Push(lua.LNil)
                L.Push(lua.LString(response.Error.Code))
		return 2
	}

	data, err := json.Marshal(response.Result)
	if err != nil {
	        L.Push(lua.LNil)
                L.Push(lua.LString(err.Error()))
		return 2
	}
//	log.Printf("[DEBUG] data: %s\n", data)
	result := make([]dslZabbixTrigger, 0)
	if err := json.Unmarshal(data, &result); err != nil {
		log.Printf("[ERROR]: Unmarshal failed: ", err)
		return 2
	}
	log.Printf("[INFO] zabbix[url: `%s`, method: `%s`, id: `%d`] returned %d element(s)\n", z.url, `trigger.get`, z.id-1, len(result))

	slice := L.CreateTable(0, len(result))
	for key, item := range result {
		log.Printf("[INFO] searching pattern `%s` in description `%s` ", pattern, item.Description)
		re, _ := regexp.Compile(pattern)
		match := re.FindStringSubmatch(item.Description)

		if match != nil {

			log.Printf("[INFO] add element (host - `%s`) to result", item.Hosts[0].Name)

			tbl := L.CreateTable(0, 10)
			slice.RawSetH(lua.LString(key), tbl)

			tbl.RawSetH(lua.LString("description"), lua.LString(item.Description))
			tbl.RawSetH(lua.LString("host"), lua.LString(item.Hosts[0].Name))
			tbl.RawSetH(lua.LString("priority"), lua.LString(item.Priority))

			lastValue := item.Items[0].LastValue

			if str, ok := lastValue.(string); ok {
				tbl.RawSetH(lua.LString("lastvalue"), lua.LString(str))
			}
		}
	}
	L.Push(slice)
	return 1
}

func checkZabbixConn(L *lua.LState) *dslZabbix {
        ud := L.CheckUserData(1)
        if v, ok := ud.Value.(*dslZabbix); ok {
                return v
        }
        L.ArgError(1, "It is not a zabbix connection")
        return nil
}

func (d *dslZabbix) buildHttpClient() error {

	client := &http.Client{}
	transport := &http.Transport{}
	if d.proxy != `` {
		proxyUrl, err := url.Parse(d.proxy)
		if err != nil {
			return err
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	if d.ignoreSsl {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client.Transport = transport

	d.client = client
	return nil
}

func (d *dslZabbix) updateReq(req *http.Request) {
	if d.basicAuthUser != `` {
		req.SetBasicAuth(d.basicAuthUser, d.basicAuthPasswd)
	}
	req.Header.Set("Content-Type", `application/json-rpc`)
}
