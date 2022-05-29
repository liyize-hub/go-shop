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

type ProductDao struct {
	*xorm.Engine
}

func NewProductDao(db *xorm.Engine) *ProductDao {
	//判断数据库连接是否存在
	if db == nil {
		db, err := datasource.NewMysqlConn()
		if err != nil {
			utils.Logger.Info("重新建立数据库连接失败", zap.Any("error", err))
		}
		return &ProductDao{db}
	}

	return &ProductDao{db}
}

func (d *ProductDao) AddProduct(product *models.Product) (err error) {
	count, err := d.Insert(product)
	if err != nil {
		utils.Logger.Error("插入失败", zap.Any("product", product))
		return
	}
	utils.SugarLogger.Infof("成功插入%d条数据,数据id为%d", count, product.ID)
	return
}

func (d *ProductDao) DeleteProductByID(productID int64) (bool, error) {
	count, err := d.ID(productID).UseBool().Update(&models.Product{Flag: 1})
	if err != nil {
		utils.Logger.Error("删除失败", zap.Int64("delete id", productID))
		return false, err
	}

	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (d *ProductDao) UpdateProductByID(productID int64, product *models.Product) (err error) {
	count, err := d.ID(productID).MustCols("flag", "num").Update(product)
	if err != nil {
		utils.Logger.Error("更新失败", zap.Int64("ProductID", productID), zap.Any("product", product))
		return
	}
	utils.SugarLogger.Infof("成功更新%d条数据,数据id为%d", count, productID)

	return
}

// 获取多个数据
func (d *ProductDao) GetProducts(product *models.Product) (*utils.ListAndCount, error) {
	products := make(map[int64]*models.Product)
	pro := []*models.Product{}
	sess := d.MustCols("flag").Limit(product.Size, (product.No-1)*product.Size).Asc("id")

	//时间范围
	if len(product.TimeRange) != 0 {
		first := time.Unix(product.TimeRange[0], 0).Format("2006-01-02 15:04:05")
		last := time.Unix(product.TimeRange[1], 0).Format("2006-01-02 15:04:05")
		sess = sess.Where("create_time between ? and ?", first, last)
	}

	err := sess.Find(products, product) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询失败", zap.Any("product", product))
		return nil, err
	}
	if len(products) == 0 {
		utils.Logger.Info("没有查到相关数据", zap.Any("product", product))
		return nil, errors.New("没有查到相关数据")
	}
	for _, v := range products {
		pro = append(pro, v)
	}

	//搜索总数
	count, _ := d.MustCols("flag").Count(product)

	return utils.Lists(pro, uint64(count)), nil
}

//获取单个数据
func (d *ProductDao) GetProduct(product *models.Product) (*models.Product, error) {
	exist, err := d.Get(product)
	if err != nil {
		utils.Logger.Error("查询失败", zap.Any("product", product))
		return nil, err
	}
	if !exist {
		utils.Logger.Info("没有查到相关数据", zap.Any("product", product))
		return nil, errors.New("没有查到相关数据")
	}

	return product, nil
}
