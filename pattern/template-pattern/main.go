package main

import "fmt"

func main() {
	smsOtp := &Sms{}

	o := Otp{
		iOtp: smsOtp,
	}
	o.genAndSendOTP(4)

	fmt.Println()
	emailOtp := &Email{}
	o = Otp{
		iOtp: emailOtp,
	}
	o.genAndSendOTP(4)
}
