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

type AdminDao struct {
	*xorm.Engine
}

func NewAdminDao(db *xorm.Engine) *AdminDao {
	//判断数据库连接是否存在
	if db == nil {
		db, err := datasource.NewMysqlConn()
		if err != nil {
			utils.Logger.Info("商铺管理员重新建立数据库连接失败", zap.Any("error", err))
		}
		return &AdminDao{db}
	}

	return &AdminDao{db}
}

func (d *AdminDao) AddAdmin(admin *models.Admin) (err error) {
	count, err := d.Insert(admin)
	if err != nil {
		utils.Logger.Error("商铺管理员插入失败", zap.Any("admin", admin))
		return
	}
	utils.SugarLogger.Infof("商铺管理员成功插入%d条数据,数据id为%d", count, admin.ID)
	return
}

func (d *AdminDao) DeleteShopByID(AdminID int64) (bool, error) {
	count, err := d.ID(AdminID).Update(&models.Admin{Flag: 3})
	if err != nil {
		utils.Logger.Error("删除商铺管理员失败", zap.Int64("delete id", AdminID))
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
		utils.Logger.Error("更新商铺管理员失败", zap.Int64("AdminID", adminID), zap.Any("admin", admin))
		return
	}
	utils.SugarLogger.Infof("商铺管理员成功更新%d条数据,数据id为%d", count, adminID)

	return
}

func (d *AdminDao) GetAdmin(admin *models.Admin) (*models.Admin, error) {

	exist, err := d.Get(admin)
	if err != nil {
		utils.Logger.Error("查询管理员用户失败", zap.Any("admin", admin))
		return &models.Admin{}, err
	}
	if !exist {
		utils.Logger.Info("商铺管理员不存在", zap.Any("admin", admin))
		return &models.Admin{}, errors.New("商铺管理员不存在，请注册后使用！")
	}

	return admin, err
}

// 获取多个数据
func (d *AdminDao) GetShops(admin *models.Admin) (*utils.ListAndCount, error) {
	admins := []*models.Admin{}

	sess := d.MustCols("flag").Asc("id")

	//如果Page项不为空
	if admin.Size != 0 && admin.No != 0 {
		sess = sess.Limit(admin.Size, (admin.No-1)*admin.Size)
	}

	//时间范围
	if len(admin.TimeRange) != 0 {
		first := time.Unix(admin.TimeRange[0], 0).Format("2006-01-02 15:04:05")
		last := time.Unix(admin.TimeRange[1], 0).Format("2006-01-02 15:04:05")
		sess = sess.Where("create_time between ? and ?", first, last)
	}

	err := sess.Find(&admins, admin) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询商铺管理员失败", zap.Any("Admin", admin))
		return nil, err
	}
	if len(admins) == 0 {
		utils.Logger.Info("商铺管理员没有查到相关数据", zap.Any("Admin", admin))
		return nil, errors.New("没有查到商铺管理员相关数据")
	}

	//搜索总数
	count, _ := d.MustCols("flag").Count(admin)

	return utils.Lists(admins, uint64(count)), nil
}

//获取单个数据
func (d *AdminDao) GetShopByID(id int64) (*utils.ListAndCount, error) {
	var ad models.Admin
	exist, err := d.ID(id).Get(&ad)
	if err != nil {
		utils.Logger.Error("查询商铺管理员失败", zap.Any("admin", ad))
		return nil, err
	}
	if !exist {
		utils.Logger.Info("没有查到商铺管理员相关数据", zap.Any("admin", ad))
		return nil, errors.New("没有查到商铺管理员相关数据")
	}

	// 返回值必须为数组
	var admin []models.Admin
	admin = append(admin, ad)

	return utils.Lists(admin, 1), nil
}
