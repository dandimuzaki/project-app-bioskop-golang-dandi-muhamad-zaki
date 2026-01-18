package utils

import (
	"fmt"
	"io"

	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"gopkg.in/mail.v2"
)

func SendEmail(data dto.EmailRequest, config Configuration, log *zap.Logger) error {
	m := mail.NewMessage()
	m.SetHeader("From", config.SMTP.Email)
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", data.Body)

	// Attach QR from memory
	for _, att := range data.Attachments {
		a := att

		m.Attach(
			a.FileName,
			mail.SetCopyFunc(func(w io.Writer) error {
					_, err := w.Write(a.FileByte)
					return err
			}),
			mail.SetHeader(map[string][]string{
					"Content-Type": {a.ContentType},
			}),
		)
	}
	
	d := mail.NewDialer("smtp.gmail.com", config.SMTP.Port, config.SMTP.Email, config.SMTP.Password)
	if err := d.DialAndSend(m); err != nil {
		log.Error("Error send email: ", zap.Error(err))
	}
	log.Info("Email send successfully")
	return nil
}

func GenerateQR(qrToken string, config Configuration, log *zap.Logger) ([]byte, error) {
	url := fmt.Sprintf("%s/tickets/verify?token=%s", config.BaseURL, qrToken)
 	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		log.Error("Error generate QR: ", zap.Error(err))
		return nil, err
	}
	return png, err
}