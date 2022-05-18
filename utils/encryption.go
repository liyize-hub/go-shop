package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ValidatePassword(password string, hashed string) (err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil { //比对成功时返回nil,失败时返回error
		Logger.Info("密码比对错误！")
		return errors.New("密码输入错误，请重新输入！")
	}
	return
}
