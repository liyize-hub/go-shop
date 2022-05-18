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

func (d *AdminDao) AddAdmin(admin *models.Admin) (err error) {
	count, err := d.Insert(admin)
	if err != nil {
		utils.Logger.Error("插入失败", zap.Any("admin", admin))
		return
	}
	utils.SugarLogger.Infof("成功插入%d条数据,数据id为%d", count, admin.ID)
	return
}
