package captchaAIO

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func NewCapMonsterClient(key string) *CapMonster {
	return &CapMonster{
		Key:  key,
		http: &http.Client{},
	}
}

type CapMonster struct {
	Key   string
	Debug bool
	http  *http.Client
}

func (cm *CapMonster) logf(f string, v ...interface{}) {
	if debug {
		log.Printf(f, v...)
	}
}

func (cm *CapMonster) SetTimeout(t time.Duration) {
	cm.http.Timeout = t
}

func (cm *CapMonster) Solve(captcha interface{}, proxy string) (string, error) {
	cm.logf("Starting CapMonster")
	id, err := cm.Send(captcha, proxy)
	if err != nil {
		return "", err
	}
	cm.logf("Solving task: %v", id)
	time.Sleep(time.Duration(5) * time.Second)
	for {
		solution, err := cm.GetRes(id)
		if err != nil {
			if errors.Is(err, ErrCaptchaNotReady) {
				cm.logf("captcha not ready, waiting 5 seconds")
				time.Sleep(time.Duration(5) * time.Second)
				continue
			}
			cm.logf(err.Error())
			return "", err
		}
		return solution, nil
	}
}
func (cm *CapMonster) GetRes(id string) (string, error) {
	payload := fmt.Sprintf(`{ "clientKey" : "%s", taskId: %s}`, cm.Key, id)
	req, err := http.NewRequest("POST", "https://api.capmonster.cloud/getTaskResult", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", ErrNetwork
	}
	resp, err := cm.http.Do(req)
	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	data := string(resBody)

	type cmSolveRes struct {
		Status   string `json:"status"`
		Solution struct {
			GRecaptchaResponse string `json:"gRecaptchaResponse"`
		} `json:"solution"`
	}
	var cmRes cmSolveRes
	err = json.Unmarshal([]byte(data), &cmRes)
	if cmRes.Status != "ready" {
		return "", ErrCaptchaNotReady
	}
	return cmRes.Solution.GRecaptchaResponse, nil

}

func (cm *CapMonster) Send(captcha interface{}, proxy string) (string, error) {
	var cmReq CapMonsterRequest

	switch t := captcha.(type) {
	case ReCaptcha:
		cmReq = cm.reCaptcha(t)
		if t.Version == "3" {
			cmReq.Task["task"] = "RecaptchaV3TaskProxyless"
		}
		if t.Version == "2" && proxy != "" {
			cmReq.Task["task"] = "NoCaptchaTask"
		} else {
			if t.Version != "2" && t.Version != "3" {
				return "", ErrUnsupportedCaptcha
			}
			cmReq.Task["task"] = "NoCaptchaTaskProxyless"

		}
	case HCaptcha:
		cmReq = cm.hCaptcha(t)
		if proxy != "" {
			cmReq.Task["task"] = "HCaptchaTaskProxyless"
		} else {
			cmReq.Task["task"] = "HCaptchaTask"
		}
	default:
		return "", ErrUnsupportedCaptcha
	}

	cmReq.ClientKey = cm.Key
	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		cmReq.Task["proxyType"] = u.Scheme
		cmReq.Task["proxyAddress"] = u.Hostname()
		cmReq.Task["proxyPort"] = u.Port()
		cmReq.Task["proxyLogin"] = u.User.Username()
		if pass, ok := u.User.Password(); ok {
			cmReq.Task["proxyPassword"] = pass
		}
	}

	body, _ := formatJSON(&cmReq)

	req, err := http.NewRequest("POST", "https://api.capmonster.cloud/createTask", bytes.NewBuffer([]byte(body)))
	if err != nil {

		return "", ErrNetwork
	}
	//	req.Header.Set("content-type", "application/json")
	resp, err := cm.http.Do(req)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {

		return "", err
	}
	data := string(resBody)

	type cmCreateTaskRes struct {
		TaskId    int    `json:"taskId"`
		ErrorCode string `json:"errorCode"`
	}

	var cmRes cmCreateTaskRes
	err = json.Unmarshal([]byte(data), &cmRes)
	if err != nil {
		return "", err
	}

	if cmRes.ErrorCode != "" {
		switch cmRes.ErrorCode {
		case "CAPCHA_NOT_READY":
			return "", ErrCaptchaNotReady
		case "ERROR_KEY_DOES_NOT_EXIST":
			return "", ErrKeyDoesNotExist
		case "ERROR_WRONG_IP_NOT_ALLOWED":
			return "", ErrIPNotAllowed
		case "ERROR_TOO_BIG_CAPTCHA_FILESIZE":
			return "", ErrTooBigCaptcha
		case "ERROR_ZERO_BALANCE":
			return "", ErrZeroBalance
		case "ERROR_CAPTCHA_UNSOLVABLE":
			return "", ErrCaptchaUnsolvable
		case "ERROR_NO_SUCH_CAPCHA_ID":
			return "", ErrNoSuchCaptchaID
		case "WRONG_CAPTCHA_ID":
			return "", ErrWrongCaptchaID
		case "ERROR_IP_BANNED":
			return "", ErrIPBanned
		case "ERROR_NO_SUCH_METHOD":
			return "", ErrNoSuchMethod
		case "ERROR_TOO_MANY_REQUESTS":
			return "", ErrTooManyRequests
		}
	}

	return strconv.Itoa(cmRes.TaskId), nil

}
func formatJSON(cmr *CapMonsterRequest) (string, error) {

	formattedJSON := fmt.Sprintf(`{
					"clientKey": "%s",
                    "task": {
						"type": "%s",
						"websiteURL": "%s",
       					"websiteKey": "%s",
                        "proxyType": "%s",
						"proxyAddress": "%s",
						"proxyPort": "%s",
						"proxyLogin": "%s",
						"proxyPassword": "%s",
						"userAgent": "%s",
						"minScore": "%s",
						"pageAction": "%s"

					 }
 


                   }`, cmr.ClientKey, cmr.Task["task"], cmr.Task["websiteURL"], cmr.Task["websiteKey"], cmr.Task["proxyType"], cmr.Task["proxyAddress"],

		cmr.Task["proxyPort"], cmr.Task["proxyLogin"], cmr.Task["proxyPassword"], cmr.Task["userAgent"], cmr.Task["minScore"], cmr.Task["pageAction"])
	return formattedJSON, nil
}

type CapMonsterRequest struct {
	ClientKey string            `json:"clientkey"`
	Task      map[string]string `json:"task"`
}

func (cm *CapMonster) reCaptcha(c ReCaptcha) CapMonsterRequest {
	req := CapMonsterRequest{}
	req.Task = make(map[string]string)

	if c.PageUrl != "" {
		req.Task["websiteURL"] = c.PageUrl
	}

	if c.SiteKey != "" {
		req.Task["websiteKey"] = c.SiteKey
	}

	if c.Score != 0 {
		req.Task["minScore"] = strconv.FormatFloat(c.Score, 'f', -1, 64)
	}
	if c.Action != "" {
		req.Task["pageAction"] = c.Action
	}
	if c.UserAgent != "" {
		req.Task["userAgent"] = c.UserAgent
	}
	return req
}

func (cm *CapMonster) hCaptcha(c HCaptcha) CapMonsterRequest {

	req := CapMonsterRequest{}
	req.Task = make(map[string]string)
	req.Task["websiteURL"] = c.PageUrl
	req.Task["websiteKey"] = c.SiteKey
	if c.UserAgent != "" {
		req.Task["userAgent"] = c.UserAgent
	}

	return req
}
func (cm *CapMonster) GetBalance() (float64, error) {
	payload := fmt.Sprintf(`{ "clientKey": "%s"  }`, cm.Key)
	req, err := http.NewRequest("POST", "https://api.capmonster.cloud/getBalance", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return 0, err
	}
	resp, err := cm.http.Do(req)
	if err != nil {
		return 0, ErrNetwork
	}
	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {

		return 0, err
	}
	data := string(resBody)
	type cmGetBalanceRes struct {
		ErrorCode string  `json:"errorCode"`
		Balance   float64 `json:"balance"`
	}
	var cmBalance cmGetBalanceRes

	err = json.Unmarshal([]byte(data), &cmBalance)
	if cmBalance.ErrorCode != "" {
		switch cmBalance.ErrorCode {
		case "CAPCHA_NOT_READY":
			return 0, ErrCaptchaNotReady
		case "ERROR_KEY_DOES_NOT_EXIST":
			return 0, ErrKeyDoesNotExist
		case "ERROR_WRONG_IP_NOT_ALLOWED":
			return 0, ErrIPNotAllowed
		case "ERROR_TOO_BIG_CAPTCHA_FILESIZE":
			return 0, ErrTooBigCaptcha
		case "ERROR_ZERO_BALANCE":
			return 0, ErrZeroBalance
		case "ERROR_CAPTCHA_UNSOLVABLE":
			return 0, ErrCaptchaUnsolvable
		case "ERROR_NO_SUCH_CAPCHA_ID":
			return 0, ErrNoSuchCaptchaID
		case "WRONG_CAPTCHA_ID":
			return 0, ErrWrongCaptchaID
		case "ERROR_IP_BANNED":
			return 0, ErrIPBanned
		case "ERROR_NO_SUCH_METHOD":
			return 0, ErrNoSuchMethod
		case "ERROR_TOO_MANY_REQUESTS":
			return 0, ErrTooManyRequests
		}
	}
	return cmBalance.Balance, nil
}
