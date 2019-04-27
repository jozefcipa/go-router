package router

var ContentType map[string]string

func init() {
	ContentType = map[string]string{
		"JSON":     "application/json",
		"text":     "text/plain",
		"formData": "multipart/form-data",
		"wwwForm":  "application/x-www-form-urlencoded",
	}
}
