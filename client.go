package captchaAIO

import (
	"os"
	"strings"
)

var debug bool

func init() {
	e := os.Getenv("GODEBUG")
	if strings.Contains(e, "captchaaiodebug=1") {
		debug = true
	}
}

type (
	Canvas struct {
		File            string
		Base64          string
		PreviousId      int
		CanSkip         bool
		Lang            string
		HintText        string
		HintImageBase64 string
		HintImageFile   string
	}

	Capy struct {
		SiteKey   string
		Url       string
		ApiServer string
	}

	Coordinates struct {
		File            string
		Base64          string
		Lang            string
		HintText        string
		HintImageBase64 string
		HintImageFile   string
	}

	FunCaptcha struct {
		SiteKey   string
		Url       string
		Surl      string
		UserAgent string
		Data      map[string]string
	}

	GeeTest struct {
		GT        string
		Challenge string
		Url       string
		ApiServer string
	}

	Grid struct {
		File            string
		Base64          string
		Rows            int
		Cols            int
		PreviousId      int
		CanSkip         bool
		Lang            string
		HintText        string
		HintImageBase64 string
		HintImageFile   string
	}

	HCaptcha struct {
		SiteKey string
		PageUrl     string
		UserAgent string
	}

	KeyCaptcha struct {
		UserId         int
		SessionId      string
		WebServerSign  string
		WebServerSign2 string
		Url            string
	}

	Normal struct {
		File            string
		Base64          string
		Phrase          bool
		CaseSensitive   bool
		Calc            bool
		Numberic        int
		MinLen          int
		MaxLen          int
		Lang            string
		HintText        string
		HintImageBase64 string
		HintImageFile   string
	}

	ReCaptcha struct {
		SiteKey   string
		PageUrl       string
		Invisible bool
		Version   string
		Action    string
		Score     float64
		UserAgent string
	}

	Rotate struct {
		File            string
		Files           []string
		Angle           int
		Lang            string
		HintText        string
		HintImageBase64 string
		HintImageFile   string
	}

	Text struct {
		Text string
		Lang string
	}
)

type Client interface {
	Solve(captcha interface{}, proxy string) (string, error)
	GetBalance() (float64, error)
}
