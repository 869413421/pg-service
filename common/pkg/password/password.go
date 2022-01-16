package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash hash加密
func Hash(password string) (string,error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "",err
	}

	return string(bytes),nil
}

//CheckHash 检查密码是否与hash值匹配
func CheckHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsHashed 检查是否已经加密过
func IsHashed(str string) bool {
	return len(str) == 60
}
