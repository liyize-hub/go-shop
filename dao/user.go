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
		db, err := datasource.NewMysqlConn()
		if err != nil {
			utils.Logger.Info("用户重新建立数据库连接失败", zap.Any("error", err))
		}
		return &UserDao{db}
	}

	return &UserDao{db}
}

func (d *UserDao) AddUser(User *models.User) (err error) {
	count, err := d.Insert(User)
	if err != nil {
		utils.Logger.Error("用户插入失败", zap.Any("User", User))
		return
	}
	utils.SugarLogger.Infof("用户成功插入%d条数据,数据id为%d", count, User.ID)
	return
}

func (d *UserDao) DeleteUserByID(UserID int64) (bool, error) {
	count, err := d.ID(UserID).Update(&models.User{Flag: 3})
	if err != nil {
		utils.Logger.Error("删除用户失败", zap.Int64("delete id", UserID))
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

func (d *UserDao) GetUser(User *models.User) (*models.User, error) {

	exist, err := d.Get(User)
	if err != nil {
		utils.Logger.Error("查询用户失败", zap.Any("User", User))
		return &models.User{}, err
	}
	if !exist {
		utils.Logger.Info("用户不存在", zap.Any("User", User))
		return &models.User{}, errors.New("用户不存在，请注册后使用！")
	}

	return User, err
}

// 获取多个数据
func (d *UserDao) GetUsers(User *models.User) (*utils.ListAndCount, error) {
	Users := []*models.User{}

	sess := d.MustCols("flag").Asc("id")

	//如果Page项不为空
	if User.Size != 0 && User.No != 0 {
		sess = sess.Limit(User.Size, (User.No-1)*User.Size)
	}

	//时间范围
	if len(User.TimeRange) != 0 {
		first := time.Unix(User.TimeRange[0], 0).Format("2006-01-02 15:04:05")
		last := time.Unix(User.TimeRange[1], 0).Format("2006-01-02 15:04:05")
		sess = sess.Where("create_time between ? and ?", first, last)
	}

	err := sess.Find(&Users, User) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询用户失败", zap.Any("User", User))
		return nil, err
	}
	if len(Users) == 0 {
		utils.Logger.Info("用户没有查到相关数据", zap.Any("User", User))
		return nil, errors.New("没有查到用户相关数据")
	}

	//搜索总数
	count, _ := d.MustCols("flag").Count(User)

	return utils.Lists(Users, uint64(count)), nil
}

//获取单个数据
func (d *UserDao) GetUserByID(id int64) (*utils.ListAndCount, error) {
	var ad models.User
	exist, err := d.ID(id).Get(&ad)
	if err != nil {
		utils.Logger.Error("查询用户失败", zap.Any("User", ad))
		return nil, err
	}
	if !exist {
		utils.Logger.Info("没有查到用户相关数据", zap.Any("User", ad))
		return nil, errors.New("没有查到用户相关数据")
	}

	// 返回值必须为数组
	var User []models.User
	User = append(User, ad)

	return utils.Lists(User, 1), nil
}
