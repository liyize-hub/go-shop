package dao

import (
	"errors"
	"go-shop/datasource"
	"go-shop/models"
	"go-shop/utils"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

type UserDao struct {
	*xorm.Engine
}

func NewUserDao(db *xorm.Engine) *UserDao {
	//判断数据库连接是否存在
	if db == nil {
		db, err := datasource.NewMysqlConn()
		if err != nil {
			utils.Logger.Info("重新建立数据库连接失败", zap.Any("error", err))
		}
		return &UserDao{db}
	}

	return &UserDao{db}
}

func (u *UserDao) GetUser(user *models.User) (*models.User, error) {

	exist, err := u.Get(user)
	if err != nil {
		utils.Logger.Error("查询管理员用户失败", zap.Any("user", user))
		return &models.User{}, err
	}
	if !exist {
		utils.Logger.Info("此用户不存在", zap.Any("user", user))
		return &models.User{}, errors.New("此用户不存在，请注册后使用！")
	}

	return user, nil
}

func (u *UserDao) AddUser(user *models.User) (err error) {
	count, err := u.Insert(user)
	if err != nil {
		utils.Logger.Error("插入失败", zap.Any("user", user))
		return
	}
	utils.SugarLogger.Infof("成功插入%d条数据,数据id为%d", count, user.ID)
	return
}
