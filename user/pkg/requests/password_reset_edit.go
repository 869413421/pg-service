package requests

import (
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/thedevsaddam/govalidator"
)

func ValidatePasswordReset(data model.User) map[string][]string {
	rules := govalidator.MapData{
		"email": []string{
			"required",
			"email",
			"between:3,30",
		},
	}

	messages := govalidator.MapData{
		"email": []string{
			"required:Email必填",
			"email:邮件格式错误",
			"between:邮件长度在3到30字符之间",
		},
	}

	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}

	errs := govalidator.New(opts).ValidateStruct()

	return errs
}
