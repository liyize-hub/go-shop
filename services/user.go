package services

import (
	"go-shop/dao"
	"go-shop/models"
	"go-shop/utils"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

/**
 * 用户服务
 * 标准的开发模式将每个实体的提供的功能以接口标准的形式定义,供控制层进行调用。
 */
type UserService interface {
	AddUser(*models.User) error                            //添加用户,注册用户
	DeleteUser(int64) error                                //删除用户
	UpdateUser(*models.User) error                         //更新用户
	SelectUsers(*models.User) (*utils.ListAndCount, error) //批量查询用户信息
	SelectUser(*models.User) (*models.User, bool, error)   //查询单个用户
	SetToken(string, interface{}, time.Duration) error     //设置用户token到redis
	GetToken(string) (string, error)                       //从redis获取用户token数据
}

// 用户服务实现结构体
type userService struct {
	db  *xorm.Engine
	rdb *redis.Client
}

// 初始化函数
func NewUserService(db *xorm.Engine, redis *redis.Client) UserService {
	return &userService{db: db, rdb: redis}
}

// 注册用户
func (u *userService) AddUser(user *models.User) (err error) {
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
func (u *userService) UpdateUser(user *models.User) (err error) {
	id := user.ID
	user.ID = 0 //清空主键
	err = dao.NewUserDao(u.db).UpdateUserByID(id, user)
	return
}

//查询用户信息
func (u *userService) SelectUsers(user *models.User) (UserListandCount *utils.ListAndCount, err error) {
	UserListandCount, err = dao.NewUserDao(u.db).GetUsers(user)
	return
}

//查询用户信息
func (u *userService) SelectUser(user *models.User) (*models.User, bool, error) {
	user, exist, err := dao.NewUserDao(u.db).GetUser(&models.User{OpenID: user.OpenID})
	return user, exist, err
}

//设置token Redis
func (u *userService) SetToken(key string, value interface{}, time time.Duration) error {
	cmd := u.rdb.Set(key, value, time)
	if cmd.Err() != nil {
		utils.Logger.Info("Set Redis error", zap.String("key", key), zap.Any("value", value))
		return cmd.Err()
	}
	return nil
}

//获取token Redis
func (u *userService) GetToken(key string) (string, error) {
	cmd := u.rdb.Get(key)
	if cmd.Err() != nil {
		utils.Logger.Info("Get Redis error", zap.String("key", key))
		return "", cmd.Err()
	}
	return cmd.Val(), nil
}
