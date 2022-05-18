package services

import (
	"go-shop/dao"
	"go-shop/models"
	"go-shop/utils"

	"xorm.io/xorm"
)

/**
 * 管理员服务
 * 标准的开发模式将每个实体的提供的功能以接口标准的形式定义,供控制层进行调用。
 */
type AdminService interface {
	GetByAdminNameAndPassword(username, password string) (*models.Admin, error) //通过用户名+密码 获取管理员实体
	AddAdmin(*models.Admin) error                                               //添加管理员
	//获取管理员总数
	//GetAdminCount() (int64, error)
}

// 管理员服务实现结构体
type adminService struct {
	db *xorm.Engine
}

// 初始化函数
func NewAdminService(db *xorm.Engine) AdminService {
	return &adminService{db: db}
}

// 通过用户名和密码查询管理员
func (ac *adminService) GetByAdminNameAndPassword(username, password string) (*models.Admin, error) {
	admin, err := dao.NewAdminDao(ac.db).GetAdmin(&models.Admin{Name: username})
	if err != nil {
		return &models.Admin{}, err
	}

	err = utils.ValidatePassword(password, admin.Pwd) //验证密码
	if err != nil {
		return &models.Admin{}, err
	}

	return admin, nil
}

// 注册管理员
func (ac *adminService) AddAdmin(admin *models.Admin) (err error) {
	pwdByte, err := utils.GeneratePassword(admin.Pwd) //加密密码
	if err != nil {
		return
	}
	admin.Pwd = string(pwdByte)
	return dao.NewAdminDao(ac.db).AddAdmin(admin)
}
