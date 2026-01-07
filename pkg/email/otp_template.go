package email

import (
	"bytes"
	"html/template"
)

type OTPEmailTemplateData struct {
	OTP           string
	ExpiryMinutes int
	Year          int
	AppName       string
}

func RenderTemplate(path string, data any) (string, error) {
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
