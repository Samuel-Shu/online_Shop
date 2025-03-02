package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/userop_srv/global"
	"online_Shop/userop_srv/model"
	"online_Shop/userop_srv/proto"
)

func (*UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var address []model.Address
	var rsp proto.AddressListResponse
	var addressResponse []*proto.AddressResponse

	if result := global.DB.Where(&model.Address{User: req.UserId}).Find(&address); result.RowsAffected != 0 {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, address := range address {
		addressResponse = append(addressResponse, &proto.AddressResponse{
			Id:           int32(address.ID),
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}
	rsp.Data = addressResponse

	return &rsp, nil
}

func (*UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	var address model.Address

	address.User = req.UserId
	address.Province = req.Province
	address.City = req.City
	address.Address = req.Address
	address.SignerName = req.SignerName
	address.SignerMobile = req.SignerMobile

	global.DB.Save(&address)

	return &proto.AddressResponse{Id: int32(address.ID)}, nil
}

func (*UserOpServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*proto.EmptyWithAddress, error) {
	if r := global.DB.Where("id = ? and user = ?", req.Id, req.UserId).Delete(&model.Address{}); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}
	return &proto.EmptyWithAddress{}, nil
}

func (*UserOpServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.EmptyWithAddress, error) {
	var address model.Address

	if r := global.DB.Where("id = ? and user = ?", req.Id, req.UserId).First(&address); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	if address.Province != "" {
		address.Province = req.Province
	}

	if address.City != "" {
		address.City = req.City
	}

	if address.District != "" {
		address.District = req.District
	}

	if address.Address != "" {
		address.Address = req.Address
	}

	if address.SignerName != "" {
		address.SignerName = req.SignerName
	}

	if address.SignerMobile != "" {
		address.SignerMobile = req.SignerMobile
	}

	global.DB.Save(&address)

	return &proto.EmptyWithAddress{}, nil
}
