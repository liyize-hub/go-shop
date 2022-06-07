package services

import (
	"go-shop/dao"
	"go-shop/datasource"
	"go-shop/models"
	"go-shop/utils"
	"strconv"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

func GetAllCategory(ctx iris.Context) {
	db, err := datasource.NewMysqlConn()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	defer db.Close()

	categories := []*models.Category{}

	// 获取选中的商品
	id, _ := strconv.ParseInt(ctx.URLParam("id"), 10, 64)
	if id != 0 {
		page, _ := strconv.Atoi(ctx.URLParam("page"))
		perpage, _ := strconv.Atoi(ctx.URLParam("perpage"))
		product := &models.Product{}
		product.CategoryID = id
		product.No = page
		product.Size = perpage
		product.Flag = 0 //商品有效
		getProductsWithCategoryID(ctx, db, product)
		return
	}

	err = db.Asc("id").MustCols("flag").Find(&categories, &models.Category{Flag: 0})
	if err != nil {
		utils.Logger.Error("商品种类查询失败", zap.Any("err", err))
		utils.SendJSON(ctx, models.ErrorCode.ERROR, "商品种类查询失败", utils.Lists("", 0))
		return
	}

	if len(categories) == 0 {
		utils.Logger.Info("没有查到相关数据", zap.Any("category", categories))
		utils.SendJSON(ctx, models.ErrorCode.NotFound, "没有查到相关数据", utils.Lists("", 0))
		return
	}

	utils.SendJSON(ctx, models.ErrorCode.SUCCESS, "查询商品种类成功", utils.Lists(categories, 0))

}

func getProductsWithCategoryID(ctx iris.Context, db *xorm.Engine, product *models.Product) {
	products, count, err := dao.NewProductDao(db).GetProducts(product)
	if err != nil {
		utils.Logger.Error("商品查询失败", zap.Any("err", err))
		utils.SendJSON(ctx, models.ErrorCode.NotFound, "没有查到相关数据", utils.Lists("", 0))
		return
	}
	//连接redis数据库
	rdb, err := datasource.NewRedisConn()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer rdb.Close()

	//将查询到的用户数据进行转换成前端需要的内容
	var respList []interface{}
	for _, v := range products {
		resp := v.ProductToRespDesc(rdb)
		//rdb.HSet(ctx, v.Name+":"+strconv.FormatInt(v.ID, 10), )
		respList = append(respList, resp)
	}

	utils.SendJSON(ctx, models.ErrorCode.SUCCESS, "查询商品成功", utils.Lists(respList, uint64(count)))
}
