package address

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"online_Shop_api/userop_web/api"
	"online_Shop_api/userop_web/forms"
	"online_Shop_api/userop_web/global"
	"online_Shop_api/userop_web/middleware"
	"online_Shop_api/userop_web/proto"
	"strconv"
)

func List(c *gin.Context) {
	request := &proto.AddressRequest{}

	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.MyClaims)

	if currentUser.AuthorityId != 2 {
		userId, _ := c.Get("userID")
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.AddressSrvClient.GetAddressList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取地址列表失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["province"] = value.Province
		reMap["city"] = value.City
		reMap["district"] = value.District
		reMap["address"] = value.Address
		reMap["signer_name"] = value.SignerName
		reMap["signer_mobile"] = value.SignerMobile

		result = append(result, reMap)
	}
	reMap["data"] = result

	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := c.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	rsp, err := global.AddressSrvClient.CreateAddress(context.Background(), &proto.AddressRequest{
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
	if err != nil {
		zap.S().Errorw("新建地址失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id

	c.JSON(http.StatusOK, request)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.AddressSrvClient.DeleteAddress(context.Background(), &proto.AddressRequest{Id: int32(i)})
	if err != nil {
		zap.S().Errorw("删除地址失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func Update(c *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := c.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.AddressSrvClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           int32(i),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		SignerMobile: addressForm.SignerMobile,
		SignerName:   addressForm.SignerName,
		Address:      addressForm.Address,
	})
	if err != nil {
		zap.S().Errorw("更新地址失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
