package services

import (
	"go-shop/dao"
	"go-shop/models"
	"go-shop/utils"

	"xorm.io/xorm"
)

/**
 * 用户服务
 * 标准的开发模式将每个实体的提供的功能以接口标准的形式定义,供控制层进行调用。
 */
type UserService interface {
	GetByUserNameAndPassword(username, password string) (*models.User, error) //通过用户名+密码 获取用户实体
	AddUser(*models.User) error                                               //添加用户,注册用户
	SelectUser(*models.User) (*utils.ListAndCount, error)                      //查询用户信息
	DeleteUser(int64) error                                                     //删除用户
	UpdateUser(*models.User) error                                             //更新用户
}

// 用户服务实现结构体
type userService struct {
	db *xorm.Engine
}

// 初始化函数
func NewUserService(db *xorm.Engine) UserService {
	return &userService{db: db}
}

// 通过用户名和密码查询用户
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

// 注册用户
func (u *userService) AddUser(user *models.User) (err error) {
	pwdByte, err := utils.GeneratePassword(user.Pwd) //加密密码
	if err != nil {
		return
	}
	user.Pwd = string(pwdByte)
	return dao.NewUserDao(u.db).AddUser(user)
}

// 删除用户
func (u *userService) DeleteUser(userID int64) (err error) {
	isOk, err := dao.NewUserDao(u.db).DeleteUserByID(userID)
	if err != nil {
		return
	}

	if isOk {
		utils.SugarLogger.Infof("删除商品成功,ID为:%d", userID)
	} else {
		utils.SugarLogger.Infof("删除商品失败,ID为:%d", userID)
	}

	return
}

// 更新用户
func (u *userService) UpdateUser(User *models.User) (err error) {
	id := User.ID
	User.ID = 0 //清空主键
	err = dao.NewUserDao(u.db).UpdateUserByID(id, User)
	return
}

//查询用户信息
func (u *userService) SelectUser(user *models.User) (UserListandCount *utils.ListAndCount, err error) {
	if user.ID == 0 {
		UserListandCount, err = dao.NewUserDao(u.db).GetUsers(user)
		return
	} else {
		UserListandCount, err = dao.NewUserDao(u.db).GetUserByID(user.ID)
	}
	return
}
