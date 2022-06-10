package services

import (
	"errors"
	"go-shop/dao"
	"go-shop/models"
	"go-shop/utils"

	"github.com/go-redis/redis"
	"xorm.io/xorm"
)

type ProductService interface {
	InsertProduct(*models.Product) error                         //添加商品
	DeleteProduct(int64) error                                   //删除商品
	UpdateProduct(*models.Product) error                         //更新商品
	SelectAllProduct() (*utils.ListAndCount, error)              //查询所有商品
	SelectProducts(*models.Product) (*utils.ListAndCount, error) //查询多条商品数据
	SelectProductByID(int64) (*models.Product, error)            //查询一条商品数据
}

// 商品服务实现结构体
type productService struct {
	db  *xorm.Engine
	rdb *redis.Client
}

// 初始化函数
func NewProductService(db *xorm.Engine, redis *redis.Client) ProductService {
	return &productService{db: db, rdb: redis}
}

func (p *productService) InsertProduct(product *models.Product) (err error) {
	if product.Num == 0 || product.Price == 0 {
		return errors.New("插入的商品数量或者价格为0")
	}

	err = dao.NewProductDao(p.db).AddProduct(product)
	return err
}

func (p *productService) DeleteProduct(productID int64) (err error) {
	isOk, err := dao.NewProductDao(p.db).DeleteProductByID(productID)
	if err != nil {
		return
	}

	if isOk {
		utils.SugarLogger.Infof("删除商品成功,ID为:%d", productID)
	} else {
		utils.SugarLogger.Infof("删除商品失败,ID为:%d", productID)
	}

	return
}

func (p *productService) UpdateProduct(product *models.Product) (err error) {
	id := product.ID
	product.ID = 0 //清空主键

	//查询商品信息
	pro, err := dao.NewProductDao(p.db).GetProduct(&models.Product{ID: id})

	//请求添加秒杀活动
	if product.Status == 1 {
		product.Num -= product.ActivityNum
		//过滤第一次添加的情况
		if pro.ActivityNum != 0 {
			product.ActivityNum += pro.ActivityNum
		}
	}
	//请求删除秒杀活动
	if product.Status == 2 {
		product.Num += product.ActivityNum
		product.ActivityNum = 0
		product.Last = 0
		product.LowPrice = 0
	}

	err = dao.NewProductDao(p.db).UpdateProductByID(id, product)
	return
}

func (p *productService) SelectAllProduct() (*utils.ListAndCount, error) {
	products, count, err := dao.NewProductDao(p.db).GetProducts(&models.Product{})
	return utils.Lists(products, uint64(count)), err
}

func (p *productService) SelectProducts(product *models.Product) (*utils.ListAndCount, error) {
	products, count, err := dao.NewProductDao(p.db).GetProducts(product)

	return utils.Lists(products, uint64(count)), err
}

func (p *productService) SelectProductByID(productID int64) (*models.Product, error) {
	product, err := dao.NewProductDao(p.db).GetProduct(&models.Product{ID: productID})
	return product, err
}
