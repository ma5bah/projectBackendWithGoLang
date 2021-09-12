package common

import (
	"regexp"
)
//@TODO add a white list of valid email provider

func EmailVerifier(emailData string) bool {

	if len(emailData) > 50 {
		return false
	}
	//re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	re := regexp.MustCompile("^[a-zA-Z0-9.!$&*+=_|~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	//re := regexp.MustCompile("^[a-zA-Z0-9.!$&*+=_|~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(emailData)

	//if strings.Count(emailData, "@") != 1 {
	//	return false
	//}
	//input := "example<script>alert('Injected!');</script>@domain.com"

	//fmt.Printf("Pattern: %v\n", re.String()) // print pattern
	//fmt.Printf("\nEmail: %v :%v\n", input, re.MatchString(input))
	//result, err := regexp.MatchString("[@]", emailData)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(result)
}
