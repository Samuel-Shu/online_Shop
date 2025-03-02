package pay

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
	"online_Shop_api/order_web/global"
	"online_Shop_api/order_web/proto"
)

func Notify(c *gin.Context) {
	//支付宝回调通知
	client, err := alipay.New(global.ServerConfig.AlipayInfo.AppID, global.ServerConfig.AlipayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝的url失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(global.ServerConfig.AlipayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝公钥失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	notice, err := client.GetTradeNotification(c.Request)
	if err != nil {
		fmt.Println("交易状态为：", notice.TradeStatus)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: notice.OutTradeNo,
		Status:  string(notice.TradeStatus),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.String(http.StatusOK, "success")
}
