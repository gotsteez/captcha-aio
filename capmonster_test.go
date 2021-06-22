package captchaAIO

import (
	"os"
	"testing"
)



var capMonsterKey = os.Getenv("APICAPMONSTER")

// TestTwoCaptcha_Type tests whether the the 2captcha client implements
// to the Client interface
func TestCApMonster_Type(t *testing.T) {
	var _ Client = &CapMonster{}
	var _ Client = NewTwoCaptchaClient("")
}
//TestCapMonster_ReCaptchaSolve tests Recpatcha v2 and logs solution
func TestCapMonster_ReCaptchaSolve(t *testing.T) {
	solver := NewCapMonsterClient(capMonsterKey)
	solver.Debug = true
	c := ReCaptcha{
		SiteKey: "6LfW6wATAAAAAHLqO2pb8bDBahxlMxNdo9g947u9",
		PageUrl: "https://recaptcha-demo.appspot.com/recaptcha-v2-checkbox.php",
		Version: "2",
	}
	res, err := solver.Solve(c, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("Solution: %s", res)
}
//TestCapMonster_GetBalance Gets current balance from cap monster
func TestCapMonster_GetBalance(t *testing.T){
	solver := NewCapMonsterClient(capMonsterKey)
	balance, err := solver.GetBalance()
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("Your Balance: %f", balance)
}
//TestCapMonster_HCaptchaSolve  tests Hcaptcha and logs solution 
func TestCapMonster_HCaptchaSolve(t *testing.T){
	solver := NewCapMonsterClient(capMonsterKey)
	solver.Debug = true
	c := HCaptcha{
		SiteKey: "51829642-2cda-4b09-896c-594f89d700cc",
		PageUrl: "http://democaptcha.com/demo-form-eng/hcaptcha.html",
	}
	res, err := solver.Solve(c, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("Solution: %s", res)

}
