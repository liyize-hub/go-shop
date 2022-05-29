package services

import (
	"errors"
	"go-shop/dao"
	"go-shop/models"
	"go-shop/utils"

	"xorm.io/xorm"
)

type ActivityService interface {
	InsertActivity(*models.Activity) error                         //添加秒杀活动
	DeleteActivity(int64) error                                    //删除秒杀活动
	UpdateActivity(*models.Activity) error                         //更新秒杀活动
	SelectAllActivity() (*utils.ListAndCount, error)               //查询所有秒杀活动
	SelectActivitys(*models.Activity) (*utils.ListAndCount, error) //查询多条秒杀活动数据
	SelectActivityByID(int64) (*models.Activity, error)            //查询一条秒杀活动数据
}

// 秒杀活动服务实现结构体
type activityService struct {
	db *xorm.Engine
}

// 初始化函数
func NewActivityService(db *xorm.Engine) ActivityService {
	return &activityService{db: db}
}

func (a *activityService) InsertActivity(activity *models.Activity) (err error) {
	// 获取商品结构体
	product, err := dao.NewProductDao(a.db).GetProduct(&models.Product{ID: activity.ProductID})
	if err != nil {
		return
	}

	// 计算商品剩余数量
	num := product.Num - activity.Num
	if num < 0 {
		return errors.New("参与秒杀活动商品数量超过已上架商品数量")
	}

	// 添加秒杀活动
	err = dao.NewActivityDao(a.db).AddActivity(activity)
	if err != nil {
		return
	}

	// 更新商品数量
	err = dao.NewProductDao(a.db).UpdateProductByID(activity.ProductID, &models.Product{Num: num})
	return
}

func (a *activityService) DeleteActivity(activityID int64) (err error) {
	isOk, err := dao.NewActivityDao(a.db).DeleteActivityByID(activityID)
	if err != nil {
		return
	}

	if isOk {
		utils.SugarLogger.Infof("删除秒杀活动成功,ID为:%d", activityID)
	} else {
		utils.SugarLogger.Infof("删除秒杀活动失败,ID为:%d", activityID)
	}

	// 获取秒杀活动结构体
	activity, err := dao.NewActivityDao(a.db).GetActivity(&models.Activity{ID: activityID})
	if err != nil {
		return
	}

	// 秒杀的商品卖光了
	if activity.Num <= 0 {
		return
	}

	// 获取商品数量，秒杀的商品有剩余
	product, err := dao.NewProductDao(a.db).GetProduct(&models.Product{ID: activity.ProductID})
	if err != nil {
		return
	}

	// 计算商品剩余数量
	num := product.Num + activity.Num
	if num <= 0 {
		return errors.New("未知错误")
	}

	//恢复商品数量
	err = dao.NewProductDao(a.db).UpdateProductByID(activity.ProductID, &models.Product{Num: num})
	return
}

func (a *activityService) UpdateActivity(activity *models.Activity) (err error) {
	id := activity.ID
	activity.ID = 0 //清空主键
	err = dao.NewActivityDao(a.db).UpdateActivityByID(id, activity)
	return
}

func (a *activityService) SelectAllActivity() (*utils.ListAndCount, error) {
	activityListandCount, err := dao.NewActivityDao(a.db).GetActivitys(&models.Activity{})
	return activityListandCount, err
}

func (a *activityService) SelectActivitys(activity *models.Activity) (*utils.ListAndCount, error) {
	activityListandCount, err := dao.NewActivityDao(a.db).GetActivitys(activity)

	return activityListandCount, err
}

func (a *activityService) SelectActivityByID(activityID int64) (*models.Activity, error) {
	activity, err := dao.NewActivityDao(a.db).GetActivity(&models.Activity{ID: activityID})
	return activity, err
}
