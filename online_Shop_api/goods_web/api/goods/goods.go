package goods

import (
	"context"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"online_Shop_api/goods_web/forms"
	"online_Shop_api/goods_web/global"
	"online_Shop_api/goods_web/proto"
	"strconv"
	"strings"
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转化为http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "内部错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func List(c *gin.Context) {
	request := &proto.GoodsFilterRequest{}

	priceMin := c.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := c.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := c.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}

	isNew := c.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := c.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := c.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := c.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	pagePerNums := c.DefaultQuery("pnum", "0")
	pagePerNumsInt, _ := strconv.Atoi(pagePerNums)
	request.PagePerNums = int32(pagePerNumsInt)

	keywords := c.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := c.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	//请求商品微服务【GoodsList】
	var err error
	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁，请稍后重试",
		})
		return
	}
	r, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginContext", c), request)
	if err != nil {
		zap.S().Errorw("[List] 查询【商品列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	e.Exit()

	c.JSON(http.StatusOK, gin.H{
		"total": r.Total,
		"data":  r.Data,
	})
}

func New(c *gin.Context) {
	GoodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&GoodsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.WithValue(context.Background(), "ginContext", c), &proto.CreateGoodsInfo{
		Name:            GoodsForm.Name,
		GoodsSn:         GoodsForm.GoodsSn,
		Stocks:          GoodsForm.Stocks,
		MarketPrice:     GoodsForm.MarketPrice,
		ShopPrice:       GoodsForm.ShopPrice,
		GoodsBrief:      GoodsForm.GoodsBrief,
		GoodsDesc:       GoodsForm.GoodsDesc,
		ShipFree:        *GoodsForm.ShipFree,
		Images:          GoodsForm.Images,
		DescImages:      GoodsForm.DescImages,
		GoodsFrontImage: GoodsForm.FrontImage,
		CategoryId:      uint32(GoodsForm.CategoryId),
		BrandId:         uint32(GoodsForm.Brand),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//todo 商品库存还未做
	c.JSON(http.StatusOK, rsp)
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", c), &proto.GoodInfoRequest{
		Id: int32(idInt),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
	}

	rsp := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"category": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}

	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteGoods(context.WithValue(context.Background(), "ginContext", c), &proto.DeleteGoodsInfo{
		Id: int32(i),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
	}

	c.Status(http.StatusOK)
	return
}

func Stocks(c *gin.Context) {
	id := c.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	//todo:商品库存接口
	return
}

func UpdateStatus(c *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := c.ShouldBindJSON(&goodsStatusForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.UpdateGoods(context.WithValue(context.Background(), "ginContext", c), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if _, err := global.GoodsSrvClient.UpdateGoods(context.WithValue(context.Background(), "ginContext", c), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      uint32(goodsForm.CategoryId),
		BrandId:         uint32(goodsForm.Brand),
	}); err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
