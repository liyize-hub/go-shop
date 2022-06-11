package dao

import (
	"errors"
	"go-shop/datasource"
	"go-shop/models"
	"go-shop/utils"
	"time"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

type UserDao struct {
	*xorm.Engine
}

func NewUserDao(db *xorm.Engine) *UserDao {
	//判断数据库连接是否存在
	if db == nil {
		db = datasource.DB
	}
	return &UserDao{db}
}

func (d *UserDao) AddUser(user *models.User) (err error) {
	count, err := d.Insert(user)
	if err != nil {
		utils.Logger.Error("用户插入失败", zap.Any("User", user))
		return
	}
	utils.SugarLogger.Infof("用户成功插入%d条数据,数据id为%d", count, user.ID)
	return
}

func (d *UserDao) DeleteUserByID(userID int64) (bool, error) {
	count, err := d.ID(userID).Update(&models.User{Flag: 3})
	if err != nil {
		utils.Logger.Error("删除用户失败", zap.Int64("delete id", userID))
		return false, err
	}

	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (d *UserDao) UpdateUserByID(UserID int64, User *models.User) (err error) {
	count, err := d.ID(UserID).MustCols("flag").Update(User)
	if err != nil {
		utils.Logger.Error("更新用户失败", zap.Int64("UserID", UserID), zap.Any("User", User))
		return
	}
	utils.SugarLogger.Infof("用户成功更新%d条数据,数据id为%d", count, UserID)

	return
}

func (d *UserDao) GetUser(user *models.User) (*models.User, bool, error) {

	exist, err := d.Get(user)
	if err != nil {
		utils.Logger.Error("查询用户失败", zap.Any("User", user))
		return &models.User{}, false, err
	}
	if !exist {
		return &models.User{}, false, nil
	}

	return user, true, err
}

// 获取多个数据
func (d *UserDao) GetUsers(user *models.User) (*utils.ListAndCount, error) {
	users := []*models.User{}

	sess := d.MustCols("flag").Asc("id")

	//如果Page项不为空
	if user.Size != 0 && user.No != 0 {
		sess = sess.Limit(user.Size, (user.No-1)*user.Size)
	}

	//时间范围
	if len(user.TimeRange) != 0 {
		first := time.Unix(user.TimeRange[0], 0).Format("2006-01-02 15:04:05")
		last := time.Unix(user.TimeRange[1], 0).Format("2006-01-02 15:04:05")
		sess = sess.Where("create_time between ? and ?", first, last)
	}

	err := sess.Find(&users, user) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询用户失败", zap.Any("User", user))
		return nil, err
	}
	if len(users) == 0 {
		utils.Logger.Info("用户没有查到相关数据", zap.Any("User", user))
		return nil, errors.New("没有查到用户相关数据")
	}

	//搜索总数
	count, _ := d.MustCols("flag").Count(user)

	return utils.Lists(users, uint64(count)), nil
}
