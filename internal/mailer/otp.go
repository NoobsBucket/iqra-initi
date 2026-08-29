package mailer


import (
	"crypto/rand"
	"fmt"
)

func GenerateOTP() (string, error) {
	// generates cryptographically secure 6 digit OTP
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// convert to 6 digit number
	otp := fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2]))
	return otp[len(otp)-6:], nil
}