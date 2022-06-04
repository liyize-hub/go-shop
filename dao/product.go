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
			utils.Logger.Info("商品数据重新建立数据库连接失败", zap.Any("error", err))
		}
		return &ProductDao{db}
	}

	return &ProductDao{db}
}

func (d *ProductDao) AddProduct(product *models.Product) (err error) {
	count, err := d.Insert(product)
	if err != nil {
		utils.Logger.Error("插入商品失败", zap.Any("product", product))
		return
	}
	utils.SugarLogger.Infof("商品成功插入%d条数据,数据id为%d", count, product.ID)
	return
}

func (d *ProductDao) DeleteProductByID(productID int64) (bool, error) {
	count, err := d.ID(productID).Update(&models.Product{Flag: 1})
	if err != nil {
		utils.Logger.Error("删除商品失败", zap.Int64("delete id", productID))
		return false, err
	}

	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (d *ProductDao) UpdateProductByID(productID int64, product *models.Product) (err error) {
	sess := d.ID(productID).MustCols("flag", "num", "activity_num")

	//删除秒杀活动操作
	if product.Status == 2 {
		sess = sess.MustCols("status","low_price","last")
		product.Status = 0
	}

	count, err := sess.Update(product)
	if err != nil {
		utils.Logger.Error("更新商品失败", zap.Int64("ProductID", productID), zap.Any("product", product))
		return
	}
	utils.SugarLogger.Infof("商品成功更新%d条数据,数据id为%d", count, productID)

	return
}

// 获取多个数据
func (d *ProductDao) GetProducts(product *models.Product) (*utils.ListAndCount, error) {
	products := []*models.Product{}

	sess := d.MustCols("flag").Asc("id")

	//如果Page项不为空
	if product.Size != 0 && product.No != 0 {
		sess = sess.Limit(product.Size, (product.No-1)*product.Size)
	}

	//时间范围
	if len(product.TimeRange) != 0 {
		first := time.Unix(product.TimeRange[0], 0).Format("2006-01-02 15:04:05")
		last := time.Unix(product.TimeRange[1], 0).Format("2006-01-02 15:04:05")
		sess = sess.Where("create_time between ? and ?", first, last)
	}

	err := sess.Find(&products, product) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询商品失败", zap.Any("product", product))
		return nil, err
	}
	if len(products) == 0 {
		utils.Logger.Info("商品没有查到相关数据", zap.Any("product", product))
		return nil, errors.New("商品没有查到相关数据")
	}

	//搜索总数
	count, _ := d.MustCols("flag").Count(product)

	return utils.Lists(products, uint64(count)), nil
}

//获取单个数据
func (d *ProductDao) GetProduct(product *models.Product) (*models.Product, error) {
	exist, err := d.MustCols("flag").Get(product)
	if err != nil {
		utils.Logger.Error("查询商品失败", zap.Any("product", product))
		return nil, err
	}
	if !exist {
		utils.Logger.Info("商品没有查到相关数据", zap.Any("product", product))
		return nil, errors.New("商品没有查到相关数据")
	}

	return product, nil
}
