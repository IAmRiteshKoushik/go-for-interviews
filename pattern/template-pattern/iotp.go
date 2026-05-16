package main

type IOtp interface {
	genRandomOTP(int) string
	saveOtpCache(string)
	getMessage(string) string
	sendNotification(string) error
}
