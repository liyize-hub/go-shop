package dao

import (
	"errors"
	"go-shop/datasource"
	"go-shop/models"
	"go-shop/utils"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

type AdminDao struct {
	*xorm.Engine
}

func NewAdminDao(db *xorm.Engine) *AdminDao {
	//判断数据库连接是否存在
	if db == nil {
		db, err := datasource.NewMysqlConn()
		if err != nil {
			utils.Logger.Info("重新建立数据库连接失败", zap.Any("error", err))
		}
		return &AdminDao{db}
	}

	return &AdminDao{db}
}

func (d *AdminDao) GetAdmin(admin *models.Admin) (*models.Admin, error) {

	exist, err := d.Get(admin)
	if err != nil {
		utils.Logger.Error("查询管理员用户失败", zap.Any("admin", admin))
		return &models.Admin{}, err
	}
	if !exist {
		utils.Logger.Info("此用户不存在", zap.Any("admin", admin))
		return &models.Admin{}, errors.New("此用户不存在，请注册后使用！")
	}

	return admin, err
}

func (d *AdminDao) GetShops(admin *models.Admin) (*utils.ListAndCount, error) {
	admins := make(map[int64]*models.Admin)
	ad := []*models.Admin{}
	err := d.MustCols("flag").Limit(admin.Size, (admin.No-1)*admin.Size).Asc("id").Find(admins, admin) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询失败", zap.Any("Admin", admin))
		return nil, err
	}
	
	for _, v := range admins {
		if v.ID != 1 { //过滤超级管理员id=1
			ad = append(ad, v)
		}
	}

	if len(ad) == 0 {
		utils.Logger.Info("没有查到相关数据", zap.Any("Admin", admin))
		return nil, errors.New("没有查到相关数据")
	}

	//搜索总数
	count, _ := d.MustCols("flag").Count(admin)

	return utils.Lists(ad, uint64(count)), nil
}

func (d *AdminDao) AddAdmin(admin *models.Admin) (err error) {
	count, err := d.Insert(admin)
	if err != nil {
		utils.Logger.Error("插入失败", zap.Any("admin", admin))
		return
	}
	utils.SugarLogger.Infof("成功插入%d条数据,数据id为%d", count, admin.ID)
	return
}

func (d *AdminDao) DeleteShopByID(AdminID int64) (bool, error) {
	count, err := d.ID(AdminID).Update(&models.Admin{Flag: 3})
	if err != nil {
		utils.Logger.Error("删除失败", zap.Int64("delete id", AdminID))
		return false, err
	}

	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (d *AdminDao) UpdateShopByID(adminID int64, admin *models.Admin) (err error) {
	count, err := d.ID(adminID).MustCols("flag").Update(admin)
	if err != nil {
		utils.Logger.Error("更新失败", zap.Int64("AdminID", adminID), zap.Any("admin", admin))
		return
	}
	utils.SugarLogger.Infof("成功更新%d条数据,数据id为%d", count, adminID)

	return
}
