package buaaclock

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	loginHeader = map[string]string{
		"Accept":           `application/json, text/javascript, */*; q=0.01`,
		"Accept-Encoding":  `gzip, deflate, br`,
		"Accept-Language":  `zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7`,
		"Cache-Control":    `no-cache`,
		"Connection":       `keep-alive`,
		"Content-Type":     `application/x-www-form-urlencoded; charset=UTF-8`,
		"Host":             `app.buaa.edu.cn`,
		"Origin":           `https://app.buaa.edu.cn`,
		"Pragma":           `no-cache`,
		"Referer":          `https://app.buaa.edu.cn/uc/wap/login?redirect=https%3A%2F%2Fapp.buaa.edu.cn%2Fsite%2FbuaaStudentNcov%2Findex`,
		"sec-ch-ua":        `"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`,
		"sec-ch-ua-mobile": `?0`,
		"Sec-Fetch-Dest":   `empty`,
		"Sec-Fetch-Mode":   `cors`,
		"Sec-Fetch-Site":   `same-origin`,
		"User-Agent":       `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36`,
		"X-Requested-With": `XMLHttpRequest`,
	}

	infoHeader = map[string]string{
		"Accept":           `application/json, text/plain, */*`,
		"Accept-Encoding":  `gzip, deflate, br`,
		"Accept-Language":  `zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7`,
		"Cache-Control":    `no-cache`,
		"Connection":       `keep-alive`,
		"Host":             `app.buaa.edu.cn`,
		"Pragma":           `no-cache`,
		"Referer":          `https://app.buaa.edu.cn/site/buaaStudentNcov/index`,
		"sec-ch-ua":        `"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`,
		"sec-ch-ua-mobile": `?0`,
		"Sec-Fetch-Dest":   `empty`,
		"Sec-Fetch-Mode":   `cors`,
		"Sec-Fetch-Site":   `same-origin`,
		"User-Agent":       `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36`,
		"X-Requested-With": `XMLHttpRequest`,
	}
)

type loginRespBody struct {
	E int    `json:"e"`
	M string `json:"m"`
}

type infoDUInfoRole struct {
	Number string `json:"number"`
}

type infoDUinfo struct {
	Realname string         `json:"realname"`
	Role     infoDUInfoRole `json:"role"`
}

type infoD struct {
	UInfo infoDUinfo     `json:"uinfo"`
	Info  map[string]any `json:"info"`
}

type infoRespBody struct {
	E int    `json:"e"`
	M string `json:"m"`
	D infoD  `json:"d"`
}

type Clock struct {
	loginURL string
	infoURL  string
	saveURL  string

	retry int

	username string
	password string

	client http.Client
}

func NewClock(configs ...Config) *Clock {
	cfg := configDefault(configs...)
	jar, _ := cookiejar.New(nil)

	return &Clock{
		loginURL: cfg.LoginURL,
		infoURL:  cfg.InfoURL,
		saveURL:  cfg.SaveURL,
		retry:    cfg.Retry,
		username: cfg.UserName,
		password: cfg.Password,
		client: http.Client{
			Jar: jar,
		},
	}
}

func (c *Clock) login() error {
	path, err := url.Parse(c.loginURL)
	if err != nil {
		return err
	}

	// set body
	reqVal := url.Values{}
	reqVal.Add("username", c.username)
	reqVal.Add("password", c.password)

	req := &http.Request{
		Method: http.MethodPost,
		URL:    path,
		Body:   io.NopCloser(strings.NewReader(reqVal.Encode())),
	}

	// set header
	for k, v := range loginHeader {
		req.Header.Add(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("login request error with status " + resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	body := loginRespBody{}
	err = json.Unmarshal(respBody, &body)
	if err != nil {
		return err
	}

	if body.E != 0 {
		return errors.New(body.M)
	}

	return nil
}

func (c *Clock) info() (*infoRespBody, error) {
	url, err := url.Parse(c.infoURL)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
	}

	// set header
	for k, v := range infoHeader {
		req.Header.Add(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("info request error with status " + resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	body := infoRespBody{}
	err = json.Unmarshal(respBody, &body)
	if err != nil {
		return nil, err
	}

	if body.E != 0 {
		return nil, errors.New(body.M)
	}

	return &body, nil
}

func (c *Clock) save(info *infoRespBody) error {
	return nil
}

func (c *Clock) exec() error {
	if err := c.login(); err != nil {
		return err
	}

	infoBody, err := c.info()
	if err != nil {
		return err
	}

	if err := c.save(infoBody); err != nil {
		return err
	}

	return nil
}

func (c *Clock) Exec() error {
	if c.retry != 0 {
		retry := NewExponentialBackoff(RetryConfig{
			MaxRetryCount: c.retry,
		})

		return retry.Retry(c.exec)
	}

	for {
		if err := c.exec(); err == nil {
			return nil
		}
	}
}
