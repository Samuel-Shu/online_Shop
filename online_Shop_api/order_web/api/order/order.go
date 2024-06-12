package order

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"online_Shop_api/order_web/api"
	"online_Shop_api/order_web/forms"
	"online_Shop_api/order_web/global"
	"online_Shop_api/order_web/proto"
	"online_Shop_api/user_web/middleware"
	"strconv"
)

func List(c *gin.Context)  {
	//订单的列表

	userId, _ := c.Get("userId")
	claims, _ := c.Get("claims")

	request := proto.OrderFilterRequest{}

	//如果是管理员用户，则返回所有的订单
	model := claims.(*middleware.MyClaims)
	if model.AuthorityId == 1 {
		request.UserId = userId.(int32)
	}

	pages := c.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	pagePerNums := c.DefaultQuery("pnum", "0")
	pagePerNumsInt, _ := strconv.Atoi(pagePerNums)
	request.PagePerNums = int32(pagePerNumsInt)


	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	orderList := make([]interface{}, 0)

	for _, item := range rsp.Data {
		tmpMap := map[string]interface{}{}

		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["add_time"] = item.AddTime

		orderList = append(orderList, tmpMap)
	}

	reMap["data"] = orderList

	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context)  {
	orderForm := forms.CreateOrderForm{}
	if err := c.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	userId, _ := c.Get("userId")

	rsp, err := global.OrderSrvClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  userId.(int32),
		Name:    orderForm.Name,
		Post:    orderForm.Post,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	//todo:返回支付宝的支付url
	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Detail(c *gin.Context)  {
	id := c.Param("id")
	userId, _ := c.Get("userId")
	claims, _ := c.Get("claims")

	i, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	//如果是管理员用户，则返回所有的订单
	request := proto.OrderRequest{
		Id: int32(i),
	}
	model := claims.(*middleware.MyClaims)
	if model.AuthorityId == 1 {
		request.UserId = userId.(int32)
	}

	rsp, err := global.OrderSrvClient.OrderDetail(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("查询订单详情失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := gin.H{}
	reMap["id"] = rsp.OrderInfo.Id
	reMap["status"] = rsp.OrderInfo.Status
	reMap["user"] = rsp.OrderInfo.UserId
	reMap["post"] = rsp.OrderInfo.Post
	reMap["total"] = rsp.OrderInfo.Total
	reMap["address"] = rsp.OrderInfo.Address
	reMap["name"] = rsp.OrderInfo.Name
	reMap["mobile"] = rsp.OrderInfo.Mobile
	reMap["pay_type"] = rsp.OrderInfo.PayType
	reMap["order_sn"] = rsp.OrderInfo.OrderSn

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		tmpMap := gin.H{
			"id": item.GoodsId,
			"name": item.GoodsName,
			"price": item.GoodsPrice,
			"image": item.GoodsImage,
			"nums": item.Nums,
		}

		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList

	c.JSON(http.StatusOK, reMap)
}
