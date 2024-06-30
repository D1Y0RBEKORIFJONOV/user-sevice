package email

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
)

func SenSecretCode(to []string) (string, error) {
	from := "diyordev3@gmail.com"
	password := "phdh ielp mjoe nvsk"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	rand_int, err := rand.Int(rand.Reader, big.NewInt(100000000))

	message := []byte("Subject: Test Email\n" +
		"\n" +
		fmt.Sprintf("%d", rand_int))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("%d", rand_int), nil
}
