package email

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func Email() {
	message := gomail.NewMessage()
	message.SetHeader("From", "femifalase228@gmail.com")
	message.SetHeader("To", "falasefemi31@gmail.com")
	message.SetHeader("Subject", "This is an email sent via Gomail and Gmail SMTP")

	message.SetBody("text/html", `
        <html>
            <body>
                <h1>This is a Test Email</h1>
                <p><b>Hello!</b> This is a test email with HTML formatting.</p>
                <p>Thanks,<br>Mailtrap</p>
            </body>
        </html>
    `)

	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"femifalase228@gmail.com",
		"urea rfjq kzig oetg",
	)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
