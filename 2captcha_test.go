package captchaAIO

import (
	"log"
	"os"
	"testing"
)

var api2CaptchaKey = os.Getenv("API2CAPTCHA")

func TestTwoCaptcha_Solve(t *testing.T) {
	solver := NewTwoCaptchaClient(api2CaptchaKey)
	solver.Debug = true
	c := ReCaptcha{
		SiteKey: "6LfW6wATAAAAAHLqO2pb8bDBahxlMxNdo9g947u9",
		PageUrl: "https://recaptcha-demo.appspot.com/recaptcha-v2-checkbox.php",
		Action: "verify",
	}
	res, err := solver.Solve(c, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println(res)
}