package services

import (
	"go-shop/dao"
	"go-shop/models"
	"go-shop/utils"

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
	db *xorm.Engine
}

// 初始化函数
func NewProductService(db *xorm.Engine) ProductService {
	return &productService{db: db}
}

func (p *productService) InsertProduct(product *models.Product) (err error) {
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
	err = dao.NewProductDao(p.db).UpdateProductByID(id, product)
	return
}

func (p *productService) SelectAllProduct() (*utils.ListAndCount, error) {
	productListandCount, err := dao.NewProductDao(p.db).GetProducts(&models.Product{})
	return productListandCount, err
}

func (p *productService) SelectProducts(product *models.Product) (*utils.ListAndCount, error) {
	productListandCount, err := dao.NewProductDao(p.db).GetProducts(product)
	return productListandCount, err
}

func (p *productService) SelectProductByID(productID int64) (*models.Product, error) {
	product, err := dao.NewProductDao(p.db).GetProduct(&models.Product{ID: productID})
	return product, err
}
