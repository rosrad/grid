//
package testlib

import (
	"log"
	"net/smtp"
)

func SendGmail(subject, msg string) error {
	to := "justnow.ren@gmail.com"
	from := "s145011"
	pwd := "ren23wan"
	auth := smtp.PlainAuth("", from, pwd, "stn.nagaokaut.ac.jp")
	body := "To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + msg
	err := smtp.SendMail("stn.nagaokaut.ac.jp:25", auth, from,
		[]string{to}, []byte(body))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
