package repo

import (
	modelBase "github.com/869413421/pg-service/common/pkg/model"
	"github.com/869413421/pg-service/common/pkg/pagination"
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/jinzhu/gorm"
)

// UserRepositoryInterface 用户CURD仓库接口
type UserRepositoryInterface interface {
	GetFirst(where map[string]interface{}) (*model.User, error)
	GetByID(id uint64) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByPhone(phone string) (*model.User, error)
	Pagination(page uint64, perPage uint32) (users []model.User, viewData pagination.ViewData, err error)
}

// UserRepository 用户仓库
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository 初始化仓库
func NewUserRepository() UserRepositoryInterface {
	db := modelBase.GetDB()
	return &UserRepository{DB: db}
}

// GetByID 根据ID获取用户
func (repo UserRepository) GetByID(id uint64) (*model.User, error) {
	user := &model.User{}
	err := repo.DB.First(&user, id).Error
	return user, err
}

// Pagination 获取分页数据
func (repo UserRepository) Pagination(page uint64, perPage uint32) (users []model.User, viewData pagination.ViewData, err error) {
	//1.初始化分页实例
	DB := modelBase.GetDB()
	DB = DB.Model(model.User{}).Order("created_at desc")
	_pager := pagination.New(DB, page, perPage)

	// 2. 获取视图数据
	viewData = _pager.Paging()

	// 3. 获取数据
	_pager.Results(&users)

	return users, viewData, nil
}

// GetByEmail 根据email获取用户
func (repo UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := repo.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

// GetByPhone 根据电话获取用户
func (repo UserRepository) GetByPhone(phone string) (*model.User, error) {
	user := &model.User{}
	err := repo.DB.Where("phone = ?", phone).First(&user).Error
	return user, err
}

// GetFirst 根据条件获取用户
func (repo UserRepository) GetFirst(where map[string]interface{}) (*model.User, error) {
	user := &model.User{}
	for key, val := range where {
		repo.DB.Where(key+"=?", val)
	}
	err := repo.DB.First(&user).Error
	return user, err
}
