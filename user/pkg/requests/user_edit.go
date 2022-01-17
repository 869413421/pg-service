package requests

import (
	"github.com/869413421/pg-service/common/pkg/types"
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/thedevsaddam/govalidator"
)

func ValidateUserEdit(data model.User) map[string][]string {
	rules := govalidator.MapData{
		"name": []string{
			"required",
			"alpha_num",
			"between:3,30",
			"not_exists:users,name," + types.UInt64ToString(data.ID),
		},
		"email": []string{
			"required",
			"email",
			"between:3,30",
			"not_exists:users,email," + types.UInt64ToString(data.ID),
		},
		"phone": []string{
			"required",
			"max_cn:15",
			"not_exists:users,phone," + types.UInt64ToString(data.ID),
		},
		"password": []string{
			"required",
			"alpha_num",
			"between:6,30",
		},
		"avatar": []string{
			"between:6,100",
		},
	}

	if data.ID > 0 {
		delete(rules, "password")
	}

	messages := govalidator.MapData{
		"name": []string{
			"required：用户名为必填选项",
			"alpha_num:只允许数字和英文",
			"between:用户名在3到30个字符之间",
			"not_exists：用户名已经存在",
		},
		"email": []string{
			"required:Email必填",
			"email:邮件格式错误",
			"between:邮件长度在3到30字符之间",
			"not_exists：邮箱已经存在",
		},
		"password": []string{
			"required:密码必填",
			"between:密码在3到30个字符之间",
		},
		"passwordComfirm": []string{
			"required:确认密码框为必填项",
		},
		"avatar": []string{
			"between:图片路径长度不符合",
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
