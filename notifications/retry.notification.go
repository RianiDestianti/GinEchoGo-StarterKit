package notifications

import (
	"net/url"
	"os"

	"gin-skeleton/repositories"
)

type retryNotification struct{}

var RetryNotification = retryNotification{}

func (notification retryNotification) Send(email string, name string, token string) {
	link := notification.url(name, email, token)
	button := notification.button()
	greeting := notification.greeting(name)
	content := notification.content(token)
	closing := notification.closing()

	subject := "Kode Verifikasi Baru Anda!"
	html := repositories.TemplateMail(greeting, content, button, link, closing)

	repositories.SendEmail(email, name, subject, html)
}

func (notification retryNotification) url(name string, email string, token string) string {
	return os.Getenv("FRONTEND_URL") + "/verification?name=" + url.QueryEscape(name) + "&email=" + url.QueryEscape(email) + "&token=" + url.QueryEscape(token)
}

func (notification retryNotification) button() string {
	return ""
}

func (notification retryNotification) greeting(firstname string) string {
	return "Halo " + firstname + "!"
}

func (notification retryNotification) content(otp string) string {
	return `Berikut kode verifikasi (OTP) baru Anda:<br><br>
        <center><strong style="font-size: 22px; letter-spacing: 2px;text-align: center;">` + otp + `</strong></center><br><br>
        Kode ini hanya berlaku selama 5 menit. Jangan bagikan kode ini kepada siapa pun.`
}

func (notification retryNotification) closing() string {
	return "Harap jangan ragu untuk menghubungi kami jika Anda memiliki pertanyaan.<br>"
}
