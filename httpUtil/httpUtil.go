package httpUtil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	GET_METHOD    = "GET"
	POST_METHOD   = "POST"
	SENDTYPE_FORM = "from"
	SENDTYPE_JSON = "json"
)

type HttpSend struct {
	Link     string
	SendType string
	Header   map[string]string
	Body     interface{}
	sync.RWMutex
}

func NewHttpSend(link string) *HttpSend {
	return &HttpSend{
		Link:     link,
		SendType: SENDTYPE_FORM,
	}
}

func (h *HttpSend) SetBody(body interface{}) {
	h.Lock()
	defer h.Unlock()
	h.Body = body
}
func (h *HttpSend) SetHeader(header map[string]string) {
	h.Lock()
	defer h.Unlock()
	h.Header = header
}
func (h *HttpSend) SetSendType(sendType string) {
	h.Lock()
	defer h.Unlock()
	h.SendType = sendType
}
func (h *HttpSend) Send(methond string) ([]byte, error) {
	var (
		req       *http.Request
		resp      *http.Response
		client    http.Client
		send_data string
		err       error
	)
	if h.Body != nil {
		if strings.ToLower(h.SendType) == SENDTYPE_JSON {
			send_body, json_err := json.Marshal(h.Body)
			if json_err != nil {
				return nil, json_err
			}
			send_data = string(send_body)
		} else {
			send_body := http.Request{}
			send_body.ParseForm()
			for k, v := range h.Body.(map[string]string) {
				send_body.Form.Add(k, v)
			}
			send_data = send_body.Form.Encode()
		}
	}
	req, err = http.NewRequest(methond, h.Link, strings.NewReader(send_data))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	//设置默认header
	if len(h.Header) == 0 {
		if strings.ToLower(h.SendType) == SENDTYPE_JSON {
			h.Header = map[string]string{
				"Contnet-Type": "application/json; charset=utf-8",
			}
		} else {
			h.Header = map[string]string{
				"Contnet-Type": "application/x-www-form-urlencoded",
			}
		}
	}
	for k, v := range h.Header {
		if strings.ToLower(k) == "host" {
			req.Host = v
		} else {
			req.Header.Add(k, v)
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (h *HttpSend) Get() ([]byte, error) {
	return h.Send(GET_METHOD)
}

func (h *HttpSend) Post() ([]byte, error) {
	return h.Send(POST_METHOD)
}

func GetUrlBuild(link string, data map[string]string) string {
	u, _ := url.Parse(link)
	q := u.Query()
	for k, v := range data {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
