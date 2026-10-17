package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// SendEmail mengirim email lewat SendGrid di background (tidak menahan response, seperti di Bun).
func SendEmail(email string, name string, subject string, html string) {
	go func() {
		if err := send(email, subject, html); err != nil {
			log.Println("[mail] Gagal kirim email:", err)
			return
		}
		log.Println("[mail] Email berhasil dikirim!")
	}()
}

func send(email string, subject string, html string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("SENDGRID_API_KEY belum diisi")
	}

	payload, _ := json.Marshal(map[string]any{
		"personalizations": []map[string]any{{"to": []map[string]string{{"email": email}}}},
		"from":             map[string]string{"email": os.Getenv("MAIL_FROM")},
		"subject":          subject,
		"content":          []map[string]string{{"type": "text/html", "value": html}},
	})

	request, err := http.NewRequest(http.MethodPost, "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{Timeout: 15 * time.Second}).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("%s: %s", response.Status, body)
	}

	return nil
}

func TemplateMail(greeting string, content string, button string, url string, closing string) string {
	buttonHTML := ""
	if button != "" {
		buttonHTML = `
                                    <tr>
                                        <td style="padding: 40px 35px; text-align: center;">
                                            <a href="` + url + `" style="border-radius: 4px; color: #fff; padding: 8px 18px; text-decoration: none; background-color: #2d3748;"
                                            target="_blank" rel="noopener">` + button + `</a>
                                        </td>
                                    </tr>`
	}

	return `
        <table cellspacing="0" border="0" cellpadding="0" width="100%" bgcolor="#f2f3f8" style="font-family: Open Sans, sans-serif;">
            <tr>
                <td>
                    <table style="background-color: #f2f3f8; max-width:670px; margin:0 auto;" width="100%" border="0"
                        align="center" cellpadding="0" cellspacing="0">
                        <tr>
                            <td style="height:50px;">&nbsp;</td>
                        </tr>
                        <tr>
                            <td>
                                <table width="95%" border="0" cellpadding="0" cellspacing="0"
                                    style="max-width:670px; background:#fff; border-radius:3px; box-shadow:0 6px 18px 0 rgba(0,0,0,.06);">
                                    <tr>
                                        <td style="height:40px;">&nbsp;</td>
                                    </tr>
                                    <tr>
                                        <td style="padding:0 35px;">
                                            <p style="font-size:15px; color:#455056; margin:8px 0 0; line-height:24px;"><b>` + greeting + `</b></p>
                                        </td>
                                    </tr>
                                    <tr>
                                        <td style="padding:0 35px;">
                                            <p style="font-size:15px; color:#455056; margin:8px 0 0; line-height:24px;">` + content + `</p>
                                        </td>
                                    </tr>` + buttonHTML + `
                                    <tr>
                                        <td style="padding:0 35px;">
                                            <p style="font-size:15px; color:#455056; margin:8px 0 0; line-height:24px;">` + closing + `</p>
                                        </td>
                                    </tr>
                                    <tr>
                                        <td style="height:40px;">&nbsp;</td>
                                    </tr>
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="height:80px;">&nbsp;</td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>`
}
