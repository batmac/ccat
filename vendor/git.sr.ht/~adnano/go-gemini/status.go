package gemini

// Status represents a Gemini status code.
type Status int

// Gemini status codes.
const (
	StatusInput                    Status = 10
	StatusSensitiveInput           Status = 11
	StatusSuccess                  Status = 20
	StatusRedirect                 Status = 30
	StatusPermanentRedirect        Status = 31
	StatusTemporaryFailure         Status = 40
	StatusServerUnavailable        Status = 41
	StatusCGIError                 Status = 42
	StatusProxyError               Status = 43
	StatusSlowDown                 Status = 44
	StatusPermanentFailure         Status = 50
	StatusNotFound                 Status = 51
	StatusGone                     Status = 52
	StatusProxyRequestRefused      Status = 53
	StatusBadRequest               Status = 59
	StatusCertificateRequired      Status = 60
	StatusCertificateNotAuthorized Status = 61
	StatusCertificateNotValid      Status = 62
)

// Class returns the status class for the status code.
// 1x becomes 10, 2x becomes 20, and so on.
func (s Status) Class() Status {
	return (s / 10) * 10
}

// String returns a text for the status code.
// It returns the empty string if the status code is unknown.
func (s Status) String() string {
	switch s {
	case StatusInput:
		return "Input"
	case StatusSensitiveInput:
		return "Sensitive input"
	case StatusSuccess:
		return "Success"
	case StatusRedirect:
		return "Redirect"
	case StatusPermanentRedirect:
		return "Permanent redirect"
	case StatusTemporaryFailure:
		return "Temporary failure"
	case StatusServerUnavailable:
		return "Server unavailable"
	case StatusCGIError:
		return "CGI error"
	case StatusProxyError:
		return "Proxy error"
	case StatusSlowDown:
		return "Slow down"
	case StatusPermanentFailure:
		return "Permanent failure"
	case StatusNotFound:
		return "Not found"
	case StatusGone:
		return "Gone"
	case StatusProxyRequestRefused:
		return "Proxy request refused"
	case StatusBadRequest:
		return "Bad request"
	case StatusCertificateRequired:
		return "Certificate required"
	case StatusCertificateNotAuthorized:
		return "Certificate not authorized"
	case StatusCertificateNotValid:
		return "Certificate not valid"
	}
	return ""
}
