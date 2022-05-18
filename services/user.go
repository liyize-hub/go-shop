package services

import (
	"go-shop/dao"
	"go-shop/models"
	"go-shop/utils"

	"xorm.io/xorm"
)

/**
 * 管理员服务
 * 标准的开发模式将每个实体的提供的功能以接口标准的形式定义,供控制层进行调用。
 */
type UserService interface {
	GetByUserNameAndPassword(username, password string) (*models.User, error) //通过用户名+密码 获取管理员实体
	AddUser(*models.User) error                                               //添加管理员
}

// 管理员服务实现结构体
type userService struct {
	db *xorm.Engine
}

// 初始化函数
func NewUserService(db *xorm.Engine) UserService {
	return &userService{db: db}
}

// 通过用户名和密码查询管理员
func (u *userService) GetByUserNameAndPassword(username, password string) (*models.User, error) {
	User, err := dao.NewUserDao(u.db).GetUser(&models.User{Name: username})
	if err != nil {
		return &models.User{}, err
	}

	err = utils.ValidatePassword(password, User.Pwd) //验证密码
	if err != nil {
		return &models.User{}, err
	}

	return User, nil
}

// 注册管理员
func (u *userService) AddUser(User *models.User) (err error) {
	pwdByte, err := utils.GeneratePassword(User.Pwd) //加密密码
	if err != nil {
		return
	}
	User.Pwd = string(pwdByte)
	return dao.NewUserDao(u.db).AddUser(User)
}
