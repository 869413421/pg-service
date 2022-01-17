package requests

import (
	"errors"
	"fmt"
	"github.com/869413421/pg-service/common/pkg/model"
	"github.com/869413421/pg-service/common/pkg/types"
	"github.com/thedevsaddam/govalidator"
	"strconv"
	"strings"
	"unicode/utf8"
)

func init() {
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {
		// 1.获取验证规则中的表名，和字段名称
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")
		//2.获取表名
		tableName := rng[0]
		dbFiled := rng[1]
		id, err := types.StringToInt(rng[2])
		if err != nil {
			return err
		}
		val := value.(string)

		//3.根据表名和字段名称获取记录总数
		type tempModel struct {
			ID uint64
		}
		var item tempModel
		model.GetDB().Table(tableName).Where(dbFiled+"=?", val).First(&item)

		//4.判断是否已经存在记录
		if item.ID != 0 && item.ID != uint64(id) {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 已经被使用", val)
		}

		return nil
	})

	// max_cn:8
	govalidator.AddCustomRule("max_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:")) //handle other error
		if valLength > l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能超过 %d 个字", l)
		}
		return nil
	})

	// min_cn:2
	govalidator.AddCustomRule("min_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:")) //handle other error
		if valLength < l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度需大于 %d 个字", l)
		}
		return nil
	})
}
