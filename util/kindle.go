package util

import (
	"bytes"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os/exec"

	"ebook/engines/book"

	"github.com/scorredoira/email"
)

func kindlegen(dir, html, mobi, locale string) {
	var args []string
	if locale == "" {
		args = []string{dir + html, "-o", mobi}
	} else {
		args = []string{dir + html, "-o", mobi, "-locale", locale}
	}
	fmt.Printf("kindlegen %s", args)
	cmd := exec.Command("kindlegen", args...)
	var in bytes.Buffer
	cmd.Stdin = &in
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	fmt.Println(string(out.Bytes()))
}

// Sendkindle 将图书发到指定邮箱
func Sendkindle(address string, bok book.Book) bool {

	localFile := bok.DownloadFile()

	if localFile != "" {
		// if strings.Contains(address, "163.com") {
		// 	return sendkindle163(address, bok.Name, localFile)
		// }
		return sendkindleqq(address, bok.Name, localFile)
	}

	return false
}

func sendkindle163(address, subject, locaFile string) bool {
	m := email.NewMessage(subject, subject)
	m.From = mail.Address{Name: "From", Address: "zzgmtv@163.com"}
	m.To = []string{address}

	// add attachments
	if err := m.Attach(locaFile); err != nil {
		log.Fatal(err)
	}

	// send it
	auth := smtp.PlainAuth("", "zzgmtv@163.com", "780227CNSCZDabc", "smtp.163.com")
	if err := email.Send("smtp.163.com:25", auth, m); err != nil {
		log.Fatal(err)
	}

	return true
}

func sendkindleqq(address, subject, locaFile string) bool {
	m := email.NewMessage(subject, subject)
	m.From = mail.Address{Name: "From", Address: "kindle@ifnfn.com"}
	m.To = []string{address}

	// add attachments
	if err := m.Attach(locaFile); err != nil {
		log.Print(err)
	}

	// send it
	auth := smtp.PlainAuth("", "294966", "780227CNSCZD", "smtp.qq.com")
	if err := email.Send("smtp.qq.com:25", auth, m); err != nil {
		log.Print(err)
	}

	return true
}
