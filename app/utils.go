package app

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func FormatDateToMySQL(date string) string {
	return date[:10] + " " + date[11:19]
}

func FormatMySQLDateToLocale(date string) string {
	parts := strings.Split(date, " ")
	datePart := parts[0]
	timePart := parts[1]
	dateParts := strings.Split(datePart, "-")
	year := dateParts[0]
	month := dateParts[1]
	day := dateParts[2]
	return fmt.Sprintf("%s/%s/%s %s", day, month, year, timePart)
}

func Md5String(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func ConvertBytesToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func SaveAttachmentToFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
