package notifications

import (
	"net/url"
	"os"

	"gin-skeleton/repositories"
)

type registrationNotification struct{}

var RegistrationNotification = registrationNotification{}

func (notification registrationNotification) Send(email string, name string, token string) {
	link := notification.url(name, email, token)
	button := notification.button()
	greeting := notification.greeting(name)
	content := notification.content(token)
	closing := notification.closing()

	subject := "Kode Verifikasi Pendaftaran Anda!"
	html := repositories.TemplateMail(greeting, content, button, link, closing)

	repositories.SendEmail(email, name, subject, html)
}

func (notification registrationNotification) url(name string, email string, token string) string {
	return os.Getenv("FRONTEND_URL") + "/verification?name=" + url.QueryEscape(name) + "&email=" + url.QueryEscape(email) + "&token=" + url.QueryEscape(token)
}

func (notification registrationNotification) button() string {
	return ""
}

func (notification registrationNotification) greeting(firstname string) string {
	return "Halo " + firstname + "!"
}

func (notification registrationNotification) content(otp string) string {
	return `Terima kasih telah mendaftar. <br><br>
        Untuk menyelesaikan proses registrasi, silakan masukkan kode verifikasi (OTP) berikut:<br><br>
        <center><strong style="font-size: 22px; letter-spacing: 2px;text-align: center;">` + otp + `</strong></center><br><br>
        Kode ini hanya berlaku selama 5 menit. Jangan bagikan kode ini kepada siapa pun.`
}

func (notification registrationNotification) closing() string {
	return "Harap jangan ragu untuk menghubungi kami jika Anda memiliki pertanyaan.<br>"
}
