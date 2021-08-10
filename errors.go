package captchaAIO

import "errors"

// General errors
var (
	ErrNetwork = errors.New("captchaAIO: network error")
	ErrUnknown = errors.New("captchaAIO: could not identify server error")
)

// Possible errors for submitting task
var (
	ErrZeroBalance           = errors.New("captchaAIO: insufficient funds in account")
	ErrPageURL               = errors.New("captchaAIO: pageURL parameter is missing")
	ErrNoSlotAvailable       = errors.New("captchaAIO: You can receive this error in two cases:\n1. If you solve token-based captchas (reCAPTCHA, hCaptcha, ArkoseLabs FunCaptcha, GeeTest, etc): the queue of your captchas that are not distributed to workers is too long. Queue limit changes dynamically and depends on total amount of captchas awaiting solution and usually it’s between 50 and 100 captchas.\n2. If you solve Normal Captcha: your maximum rate for normal captchas is lower than current rate on the server")
	ErrZeroCaptchaFilesize   = errors.New("captchaAIO: image size is less than 100 bytes")
	ErrTooBigCaptcha         = errors.New("captchaAIO: image file size exceeds 100kB")
	ErrWrongFileExtension    = errors.New("captchaAIO: image file has unsupported extension. Accepted extensions: jpg, jpeg, gif, png")
	ErrImageTypeNotSupported = errors.New("captchaAIO: server can't recognize image file type")
	ErrUpload                = errors.New("captchaAIO: server couldn't get image from POST request")
	ErrIPNotAllowed          = errors.New("captchaAIO: request is sent from IP not on list of your allowed IPs")
	ErrIPBanned              = errors.New("captchaAIO: IP banned from service")
	ErrBadTokenOrPageURL     = errors.New("captchaAIO; when sending recaptcha V2, the paired page URL and sitekey are mismatched")
	ErrGoogleKey             = errors.New("captchaAIO: googlekey parameter is missing from request")
	ErrCaptchaImageBlocked   = errors.New("captchaAIO: sent image that is marked unrecognizable in service database")
	ErrTooManyBadImages      = errors.New("captchaAIO: sent too many unrecognizable images")
	ErrMaxUserTurn           = errors.New("captchaAIO: submitted too many captchas to service")
	ErrBadParameters         = errors.New("captchaAIO: required parameters are missing in request, or in incorrect format")
	ErrNoSuchCaptchaID       = errors.New("captchaAIO: Captcha you are requesting does not exist in your current captcha list or has been expired")
	ErrNoSuchMethod          = errors.New("captchaAIO: Request to API made with method which does not exist")
)

// Possible errors from results
var (
	ErrCaptchaNotReady   = errors.New("captchaAIO: Captcha not ready")
	ErrCaptchaUnsolvable = errors.New("captchaAIO: Captcha not solvable")
	ErrWrongUserKey      = errors.New("captchaAIO: provided key parameter value in incorrect format")
	ErrKeyDoesNotExist   = errors.New("captchaAIO: provided key does not exist")
	ErrWrongIDFormat     = errors.New("captchaAIO: provided captcha ID din wrong format")
	ErrWrongCaptchaID    = errors.New("captchaAIO: provided incorrect captcha ID")
	ErrBadDuplicates     = errors.New("captchaAIO: max numbers of tries is reached but min number of matches not found")
	ErrReportNotReported = errors.New("captchaAIO: already complained lots of correctly solved captchas (more than 40%) or more than 15 minutes passed after submitted captcha")
	ErrDuplicateReport   = errors.New("captchaAIO: reported the same captcha more than once")
	ErrTooManyRequests   = errors.New("captchaAIO: account is temporarily suspended, too many requests") // ERROR: NNNN
	// ErrIPAddress         = errors.New("captchaAIO: request is coming from an IP address that doesn't match the IP address of pingback IP or domain")
	ErrTokenExpired      = errors.New("captchaAIO: error code when sending GeeTest, challenge value provided is expired")
	ErrEmptyAction       = errors.New("captchaAIO: action parameter is missing or no value is provided")
	ErrProxyConnFail     = errors.New("captchaAIO: service could not connect to proxy")
)
