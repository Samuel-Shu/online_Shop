package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"online_Shop_api/inventory_web/global"
	"online_Shop_api/inventory_web/proto"
)

func List(c *gin.Context) {
	goodsId, _ := c.Get("goods_id")
	detail, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{GoodsId: int32(goodsId.(uint))})
	if err != nil {
		zap.S().Errorw("获取库存详情失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, detail)

}

func Delete(c *gin.Context) {

}
func New(c *gin.Context) {

}

func Update(c *gin.Context) {

}
