package captchaAIO

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func NewTwoCaptchaClient(key string) *TwoCaptcha {
	return &TwoCaptcha{
		Key:  key,
		http: &http.Client{},
	}
}

type TwoCaptcha struct {
	Key    string
	SoftID string
	Debug bool
	http   *http.Client
}

func (tc *TwoCaptcha) logf(f string, v ...interface{}) {
	if debug {
		log.Printf(f, v...)
	}
}

func (tc *TwoCaptcha) SetTimeout(t time.Duration) {
	tc.http.Timeout = t
}

var ErrUnsupportedCaptcha = errors.New("captchaAIO: could not identify given captcha type")

func (tc *TwoCaptcha) Solve(captcha interface{}, proxy string) (string, error) {
	tc.logf("Starting")
	id, err := tc.Send(captcha, proxy)
	if err != nil {
		return "", err
	}
	tc.logf("Solving task %v", id)
	time.Sleep(time.Duration(20) * time.Second)
	for {
		result, err := tc.GetRes(id)
		if err != nil {
			if errors.Is(err, ErrCaptchaNotReady) {
				tc.logf("captcha not ready, waiting 5 seconds")
				time.Sleep(time.Duration(5) * time.Second)
				continue
			}
			tc.logf(err.Error())
			return "", err
		}
		tc.logf("Result: %v", result)
		return result[3:], nil
	}
}

// Send sends a captcha task to be solved by twocaptcha, will return the id of the task
// or an error
func (tc *TwoCaptcha) Send(captcha interface{}, proxy string) (string, error) {
	var req TwoCaptchaRequest
	switch v := captcha.(type) {
	case ReCaptcha:
		req = tc.reCaptcha(v)
	default:
		return "", ErrUnsupportedCaptcha
	}
	req.Params["key"] = tc.Key
	tc.logf("%v", req)

	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		req.Params["proxytype"] = u.Scheme
		if pass, ok := u.User.Password(); ok {
			req.Params["proxy"] = fmt.Sprintf("%v:%v@%v", u.User.Username(), pass, u.Host)
		} else {
			req.Params["proxy"] = fmt.Sprintf("%v@%v", u.User.Username(), u.Host)
		}
	}
	if tc.SoftID != "" {
		req.Params["soft_id"] = tc.SoftID
	}

	var resp *http.Response
	if req.Files != nil && len(req.Files) > 0 {
		body := &bytes.Buffer{}
		w := multipart.NewWriter(body)
		for name, path := range req.Files {
			file, err := os.Open(path)
			if err != nil {
				return "", err
			}

			part, err := w.CreateFormFile(name, filepath.Base(path))
			if err != nil {
				return "", err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return "", err
			}
			file.Close()
		}
		for k, v := range req.Params {
			err := w.WriteField(k, v)
			if err != nil {
				return "", err
			}
		}

		req, err := http.NewRequest("POST", "https://2captcha.com/in.php", body)
		if err != nil {
			return "", ErrNetwork
		}
		req.Header.Set("content-type", w.FormDataContentType())
		resp, err = tc.http.Do(req)
		if err != nil {
			return "", err
		}
	} else {
		form := url.Values{}
		for k, v := range req.Params {
			form.Add(k, v)
		}

		r, err := http.NewRequest("POST", "https://2captcha.com/in.php", nil)
		if err != nil {
			return "", err
		}
		r.URL.RawQuery = form.Encode()
		resp, err = tc.http.Do(r)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	data := string(body)
	tc.logf("task submission request returned data %v", data)
	if strings.HasPrefix(data, "ERROR") {
		switch data {
		case "ERROR_WRONG_USER_KEY":
			return "", ErrWrongUserKey
		case "ERROR_KEY_DOES_NOT_EXIST":
			return "", ErrKeyDoesNotExist
		case "ERROR_ZERO_BALANCE":
			return "", ErrZeroBalance
		case "ERROR_PAGEURL":
			return "", ErrPageURL
		case "ERROR_NO_SLOT_AVAILABLE":
			return "", ErrNoSlotAvailable
		case "ERROR_ZERO_CAPTCHA_FILESIZE":
			return "", ErrZeroCaptchaFilesize
		case "ERROR_TOO_BIG_CAPTCHA_FILESIZE":
			return "", ErrTooBigCaptcha
		case "ERROR_WRONG_FILE_EXTENSION":
			return "", ErrWrongFileExtension
		case "ERROR_IMAGE_TYPE_NOT_SUPPORTED":
			return "", ErrImageTypeNotSupported
		case "ERROR_UPLOAD":
			return "", ErrUpload
		case "ERROR_IP_NOT_ALLOWED":
			return "", ErrIPNotAllowed
		case "IP_BANNED":
			return "", ErrIPBanned
		case "ERROR_BAD_TOKEN_OR_PAGEURL":
			return "", ErrBadTokenOrPageURL
		case "ERROR_GOOGLEKEY":
			return "", ErrGoogleKey
		case "ERROR_WRONG_GOOGLEKEY":
			return "", ErrGoogleKey
		case "ERROR_CAPTCHAIMAGE_BLOCKED":
			return "", ErrCaptchaImageBlocked
		case "TOO_MANY_BAD_IMAGES":
			return "", ErrTooManyBadImages
		case "MAX_USER_TURN":
			return "", ErrMaxUserTurn
		case "ERROR: NNNN":
			return "", ErrTooManyRequests
		case "ERROR_BAD_PARAMETERS":
			return "", ErrBadParameters
		case "ERROR_BAD_PROXY":
			return "", ErrProxyConnFail
		default:
			return "", ErrUnknown
		}
	}

	return data[3:], nil
}

func (tc *TwoCaptcha) GetRes(id string) (string, error) {
	req, err := http.NewRequest("GET", "https://2captcha.com/res.php", nil)
	if err != nil {
		return "", err
	}
	q := url.Values{}
	q.Add("key", tc.Key)
	q.Add("action", "get")
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()
	resp, err := tc.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	tc.logf("res got response %v", string(body))
	switch string(body) {
	case "CAPCHA_NOT_READY":
		return "", ErrCaptchaNotReady
	case "ERROR_CAPTCHA_UNSOLVABLE":
		return "", ErrCaptchaUnsolvable
	case "ERROR_KEY_DOES_NOT_EXIST":
		return "", ErrKeyDoesNotExist
	case "ERROR_WRONG_ID_FORMAT":
		return "", ErrWrongIDFormat
	case "ERROR_WRONG_CAPTCHA_ID":
		return "", ErrWrongCaptchaID
	case "ERROR_BAD_DUPLICATES":
		return "", ErrBadDuplicates
	case "ERROR: NNNN":
		return "", ErrTooManyRequests
	case "ERROR_IP_ADDRES":
		return "", ErrIPNotAllowed
	case "ERROR_TOKEN_EXPIRED":
		return "", ErrTokenExpired
	case "ERROR_EMPTY_ACTION":
		return "", ErrEmptyAction
	case "ERROR_PROXY_CONNECTION_FAILED":
		return "", ErrProxyConnFail
	default:
		return string(body), nil
	}
}

func (tc *TwoCaptcha) Report(id string, correct bool) error {
	req, err := http.NewRequest("GET", "https://2captcha.com/res.php", nil)
	if err != nil {
		return err
	}
	q := url.Values{}
	if correct {
		q.Add("action", "reportgood")
	} else {
		q.Add("action", "reportbad")
	}
	q.Add("id", id)
	resp, err := tc.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch string(body) {
	case "ERROR_EMPTY_ACTION":
		return ErrEmptyAction
	case "REPORT_NOT_RECORDED":
		return ErrReportNotReported
	case "ERROR_DUPLICATE_REPORT":
		return ErrDuplicateReport
	default:
		return nil
	}
}

func (tc *TwoCaptcha) GetBalance() (float64, error) {
	req, err := http.NewRequest("GET", "https://2captcha.com/res.php", nil)
	if err != nil {
		return 0, err
	}
	q := url.Values{}
	q.Add("key", tc.Key)
	q.Add("action", "getbalance")
	req.URL.RawQuery = q.Encode()
	resp, err := tc.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	bal, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return 0, err
	}
	return bal, nil
}

type TwoCaptchaRequest struct {
	Params map[string]string
	Files  map[string]string
}

func (*TwoCaptcha) reCaptcha(c ReCaptcha) TwoCaptchaRequest {
	req := TwoCaptchaRequest{
		Params: map[string]string{
			"method": "userrecaptcha",
		},
	}
	if c.SiteKey != "" {
		req.Params["googlekey"] = c.SiteKey
	}
	if c.PageUrl != "" {
		req.Params["pageurl"] = c.PageUrl
	}
	if c.Invisible {
		req.Params["invisible"] = "1"
	}
	if c.Version != "" {
		req.Params["version"] = c.Version
	}
	if c.Action != "" {
		req.Params["action"] = c.Action
	}
	if c.Score != 0 {
		req.Params["min_score"] = strconv.FormatFloat(c.Score, 'f', -1, 64)
	}

	return req
}

func (*TwoCaptcha) hCaptcha(c HCaptcha) TwoCaptchaRequest {
	req := TwoCaptchaRequest{
		Params: map[string]string{
			"method": "hcaptcha",
		},
	}
	if c.SiteKey != "" {
		req.Params["googlekey"] = c.SiteKey
	}
	if c.PageUrl != "" {
		req.Params["pageurl"] = c.PageUrl
	}
	return req
}