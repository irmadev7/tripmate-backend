package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type EmailRequest struct {
	Sender struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"to"`
	Subject     string `json:"subject"`
	HtmlContent string `json:"htmlContent"`
}

func SendEmail(toEmail, toName string) error {
	url := os.Getenv("BREVO_EMAIL_URL")

	body := EmailRequest{}
	body.Sender.Email = os.Getenv("BREVO_EMAIL_SENDER")
	body.Sender.Name = os.Getenv("BREVO_EMAIL_SENDER_NAME")

	body.To = []struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		{
			Email: toEmail,
			Name:  toName,
		},
	}

	body.Subject = "Welcome to Tripmate!"
	body.HtmlContent = "<h1>Welcome!</h1><p>Thanks for registering.</p>"

	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", os.Getenv("BREVO_API_KEY"))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
