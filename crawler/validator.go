package crawler

import (
	"net/http"
	"regexp"
	"strings"
)

func IsValidURL(url string) bool {
	urlRegex := `^(https?:\/\/)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`
	re := regexp.MustCompile(urlRegex)
	return re.MatchString(url)
}

func NormalizeURL(url string) string {
	url = strings.TrimSpace(url)
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "https://" + url
}

func GetStatusCodeDescription(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "Bad Request - The server cannot process the request"
	case http.StatusUnauthorized:
		return "Unauthorized - Authentication is required"
	case http.StatusForbidden:
		return "Forbidden - Access to this resource is denied"
	case http.StatusNotFound:
		return "Not Found - The requested page does not exist"
	case http.StatusInternalServerError:
		return "Internal Server Error - The server encountered an error"
	default:
		if statusCode >= 400 && statusCode < 500 {
			return "Client Error - There's an issue with the request"
		} else if statusCode >= 500 {
			return "Server Error - The server encountered an error"
		}
		return "Unexpected status code"
	}
}
