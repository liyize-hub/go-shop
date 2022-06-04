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

type ActivityDao struct {
	*xorm.Engine
}

func NewActivityDao(db *xorm.Engine) *ActivityDao {
	//判断数据库连接是否存在
	if db == nil {
		db, err := datasource.NewMysqlConn()
		if err != nil {
			utils.Logger.Info("秒杀活动重新建立数据库连接失败", zap.Any("error", err))
		}
		return &ActivityDao{db}
	}

	return &ActivityDao{db}
}

func (a *ActivityDao) AddActivity(activity *models.Activity) (err error) {
	if activity.Num == 0 || activity.Price == 0 || activity.Last == 0 || activity.ProductID == 0 {
		return errors.New("插入的秒杀活动的秒杀活动的数量或者价格为0或者持续时间为0或者ID为0")
	}
	count, err := a.Insert(activity)
	if err != nil {
		utils.Logger.Error("插入秒杀活动失败", zap.Any("Activity", activity))
		return
	}
	utils.SugarLogger.Infof("秒杀活动成功插入%d条数据,数据id为%d", count, activity.ID)
	return
}

func (a *ActivityDao) DeleteActivityByID(activityID int64) (bool, error) {
	count, err := a.ID(activityID).Update(&models.Activity{Flag: 1})
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
func (d *ActivityDao) GetActivitys(activity *models.Activity) (*utils.ListAndCount, error) {
	activitys := []*models.Activity{}

	sess := d.MustCols("flag").Asc("id")

	//如果Page项不为空
	if activity.Size != 0 && activity.No != 0 {
		sess = sess.Limit(activity.Size, (activity.No-1)*activity.Size)
	}

	//时间范围
	if len(activity.TimeRange) != 0 {
		first := time.Unix(activity.TimeRange[0], 0).Format("2006-01-02 15:04:05")
		last := time.Unix(activity.TimeRange[1], 0).Format("2006-01-02 15:04:05")
		sess = sess.Where("create_time between ? and ?", first, last)
	}

	err := sess.Find(&activitys, activity) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询秒杀活动失败", zap.Any("Activity", activity))
		return nil, err
	}
	if len(activitys) == 0 {
		utils.Logger.Info("秒杀活动没有查到相关数据", zap.Any("Activity", activity))
		return nil, errors.New("秒杀活动没有查到相关数据")
	}

	//搜索总数
	count, _ := d.MustCols("flag").Count(activity)

	return utils.Lists(activitys, uint64(count)), nil
}

//获取单个数据
func (d *ActivityDao) GetActivity(activity *models.Activity) (*models.Activity, error) {
	exist, err := d.MustCols("flag").Get(activity)
	if err != nil {
		utils.Logger.Error("查询秒杀活动失败", zap.Any("Activity", activity))
		return nil, err
	}
	if !exist {
		utils.Logger.Info("秒杀活动没有查到相关数据", zap.Any("Activity", activity))
		return nil, errors.New("秒杀活动没有查到相关数据")
	}

	return activity, nil
}
