package buaaclock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	loginHeader = map[string]string{
		"Accept":           `application/json, text/javascript, */*; q=0.01`,
		"Accept-Encoding":  `gzip, deflate, br`,
		"Accept-Language":  `zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7`,
		"Cache-Control":    `no-cache`,
		"Connection":       `keep-alive`,
		"Content-Length":   `43`,
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
)

type loginReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRespBody struct {
	E int    `json:"e"`
	M string `json:"m"`
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
	cfg := overwriteConfig(configs...)
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
	url, err := url.Parse(c.loginURL)
	if err != nil {
		return err
	}

	// set body
	data := loginReqBody{
		Username: c.username,
		Password: c.password,
	}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(jsonBody)),
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
		return nil
	}

	if body.E != 0 {
		return errors.New(body.M)
	}

	return nil
}

func (c *Clock) info() error {
	return nil
}

func (c *Clock) save() error {
	return nil
}

func (c *Clock) Exec() error {
	count, success := 0, false

	for {
		// return when success
		if success {
			return nil
		}

		// return when fail count larger than retry count
		if c.retry != 0 && count >= c.retry {
			return errors.New("can't clock, please try again")
		}

		if err := c.login(); err != nil {
			count++
			continue
		}

		if err := c.info(); err != nil {
			count++
			continue
		}

		if err := c.save(); err != nil {
			count++
			continue
		}

		count++
		success = true
	}
}
