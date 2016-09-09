package main

import (
	"fmt"
	"net/smtp"
)

func SendEmailForPasswordRestore(email, id string) error {

	headers := fmt.Sprintf("From: %v\r\nTo: %v\r\nSubject: Восстановления пароля\r\nContent-Type: text/html; charset='UTF-8'\r\n\r\n", SMTPFrom, email)

	msg := fmt.Sprintf(`<html>
<header>
  <meta charset='utf-8' />
</header>
<body>
Здравствуйте!<br/><br/>
Для восстановления пароля перейдите по ссылке:<a href='http://yorkina.ru/restorepassword/?token=%v'>http://yorkina.ru/restorepassword/?token=%v</a> 
<br/>
<br/>
С уважением,<br/>
webmaster@yorkina.ru  
</body>
</html>`, id, id)

	emailText := headers + msg
	auth := smtp.PlainAuth("", SMTPUser, SMTPPass, SMTPServer)
	err := smtp.SendMail(SMTPServer+":"+SMTPPort, auth, SMTPFrom, []string{email}, []byte(emailText))

	return err
}
