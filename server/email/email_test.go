package email

import (
	"fmt"
	"testing"
)

func TestMail(t *testing.T) {
	msg, id, err := Send("cdutwhu@outlook.com", "Fancy subject!", "Hello from Mailgun Go!")
	if err == nil {
		fmt.Println("message:", msg)
		fmt.Println("id:", id)
	} else {
		fmt.Println(err)
	}
}
