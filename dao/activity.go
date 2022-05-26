package dao

import (
	"errors"
	"go-shop/datasource"
	"go-shop/models"
	"go-shop/utils"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

type ActivityDao struct {
	*xorm.Engine
}

func NewActivityDao(db *xorm.Engine) *ActivityDao {
	//判断数据库连接是否存在
	if db == nil {
		db, err := datasource.NewMysqlConn()
		if err != nil {
			utils.Logger.Info("重新建立数据库连接失败", zap.Any("error", err))
		}
		return &ActivityDao{db}
	}

	return &ActivityDao{db}
}

func (a *ActivityDao) AddActivity(activity *models.Activity) (err error) {
	if activity.Num == 0 || activity.Price == 0 || activity.Last == 0 {
		return errors.New("插入的秒杀活动的商品的数量或者价格为0或者持续时间为0")
	}
	count, err := a.Insert(activity)
	if err != nil {
		utils.Logger.Error("插入秒杀活动失败", zap.Any("Activity", activity))
		return
	}
	utils.SugarLogger.Infof("成功插入%d条数据,数据id为%d", count, activity.ID)
	return
}

func (a *ActivityDao) DeleteActivityByID(activityID int64) (bool, error) {
	count, err := a.ID(activityID).UseBool().Update(&models.Activity{Flag: 1})
	if err != nil {
		utils.Logger.Error("删除秒杀活动失败", zap.Int64("delete id", activityID))
		return false, err
	}

	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (a *ActivityDao) UpdateActivityByID(ActivityID int64, Activity *models.Activity) (err error) {
	count, err := a.ID(ActivityID).MustCols("flag", "num").Update(Activity)
	if err != nil {
		utils.Logger.Error("更新秒杀活动失败", zap.Int64("ActivityID", ActivityID), zap.Any("Activity", Activity))
		return
	}
	utils.SugarLogger.Infof("成功更新秒杀活动%d条数据,数据id为%d", count, ActivityID)

	return
}

// 获取多个数据
func (a *ActivityDao) GetActivitys(Activity *models.Activity) (*utils.ListAndCount, error) {
	Activitys := make(map[int64]*models.Activity)
	pro := []*models.Activity{}
	err := a.MustCols("flag").Limit(Activity.Size, (Activity.No-1)*Activity.Size).Asc("id").Find(Activitys, Activity) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询秒杀活动失败", zap.Any("Activity", Activity))
		return nil, err
	}
	if len(Activitys) == 0 {
		utils.Logger.Info("没有查到相关数据", zap.Any("Activity", Activity))
		return nil, errors.New("没有查到相关数据")
	}
	for _, v := range Activitys {
		pro = append(pro, v)
	}

	//搜索总数
	count, _ := a.MustCols("flag").Count(Activity)

	return utils.Lists(pro, uint64(count)), nil
}

// 获取单个数据
func (a *ActivityDao) GetActivity(activity *models.Activity) (*models.Activity, error) {
	exist, err := a.Get(activity)
	if err != nil {
		utils.Logger.Error("查询秒杀活动失败", zap.Any("Activity", activity))
		return nil, err
	}
	if !exist {
		utils.Logger.Info("没有查到相关数据", zap.Any("Activity", activity))
		return nil, errors.New("没有查到相关数据")
	}

	return activity, nil
}
