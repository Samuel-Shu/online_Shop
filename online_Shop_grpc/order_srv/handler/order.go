package handler

import (
	"context"
	_ "crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/order_srv/global"
	"online_Shop/order_srv/model"
	"online_Shop/order_srv/proto"
	"time"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func GenerateOrderSn(userId int32) string {
	/*
		订单号生成规则
			年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(uint64(time.Now().UnixNano()))
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

func (o *OrderServer) CartItemList(c context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	//获取用户的购物车列表
	var shopCarts []model.ShoppingCart

	result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)

	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.CartItemListResponse{
		Total: int32(result.RowsAffected),
	}

	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      int32(shopCart.ID),
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}

	return rsp, nil
}

func (o *OrderServer) CreateCartItem(c context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//将商品添加到购物车
	//1、购物车中原本没有这件商品，--直接新建一条记录即可
	//2、这个商品之前添加到了购物车 --合并记录（主要是对购买个数的修改）
	var shopCart model.ShoppingCart
	if r := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); r.RowsAffected == 1 {
		//更新操作，如果记录已经存在，就合并购物车记录
		shopCart.Nums += req.Nums
	} else {
		//插入操作，直接添加一条新的记录即可
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: int32(shopCart.ID)}, nil
}

func (o *OrderServer) UpdateCartItem(c context.Context, req *proto.CartItemRequest) (*proto.MyEmptyWithOrder, error) {
	//更新购物车记录，更新数量和选中状态

	var shopCart model.ShoppingCart

	if r := global.DB.First(&shopCart, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	shopCart.Checked = req.Checked

	global.DB.Save(&shopCart)

	return &proto.MyEmptyWithOrder{}, nil
}

func (o *OrderServer) DeleteCartItem(c context.Context, req *proto.CartItemRequest) (*proto.MyEmptyWithOrder, error) {

	if r := global.DB.Delete(&model.ShoppingCart{}, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	return &proto.MyEmptyWithOrder{}, nil
}

func (o *OrderServer) OrderList(c context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp *proto.OrderListResponse
	var total int64

	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)
	//分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&orders)

	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      int32(order.ID),
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SignerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return rsp, nil
}

func (o *OrderServer) OrderDetail(c context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp *proto.OrderInfoDetailResponse

	//这个订单id是否是当前用户的订单，如果在web传输过来一个id的订单，web层应该先查询一下订单id是否是当前用户的订单
	//如果是后台管理系统，那么只传递order的id，如果是电商系统，还需要传递一个用户id
	if r := global.DB.Where("id = ? AND user_id = ?", req.Id, req.UserId).First(&order); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	rsp.OrderInfo.Id = int32(order.ID)
	rsp.OrderInfo.UserId = order.User
	rsp.OrderInfo.OrderSn = order.OrderSn
	rsp.OrderInfo.PayType = order.PayType
	rsp.OrderInfo.Status = order.Status
	rsp.OrderInfo.Post = order.Post
	rsp.OrderInfo.Total = order.OrderMount
	rsp.OrderInfo.Address = order.Address
	rsp.OrderInfo.Name = order.SignerName
	rsp.OrderInfo.Mobile = order.SignerMobile

	var orderGoods []model.OrderGoods
	global.DB.Where(&model.OrderGoods{Order: int32(order.ID)}).Find(&orderGoods)
	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			Nums:       orderGood.Nums,
		})
	}
	return rsp, nil
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	parentSpan := opentracing.SpanFromContext(o.Ctx)

	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopCart", opentracing.ChildOf(parentSpan.Context()))

	if r := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); r.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
	}
	shopCartSpan.Finish()
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	//跨服务调用 -- 商品微服务
	queryGoodsSpan := opentracing.GlobalTracer().StartSpan("query_goods", opentracing.ChildOf(parentSpan.Context()))
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsInfo{
		Id: goodsIds,
	})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	queryGoodsSpan.Finish()

	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	//跨微服务调用 -- 库存微服务
	shopInvSpan := opentracing.GlobalTracer().StartSpan("query_inv", opentracing.ChildOf(parentSpan.Context()))
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: goodsInvInfo,
	}); err != nil {
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}
	shopInvSpan.Finish()

	//生成订单表
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount
	saveOrderSpan := opentracing.GlobalTracer().StartSpan("save_order", opentracing.ChildOf(parentSpan.Context()))
	if r := tx.Save(&orderInfo); r.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "订单创建失败"
		return primitive.CommitMessageState
	}
	saveOrderSpan.Finish()

	o.OrderAmount = orderAmount
	o.ID = int32(orderInfo.ID)
	for _, orderGood := range orderGoods {
		orderGood.Order = int32(orderInfo.ID)
	}

	//批量插入orderGoods
	saveOrderGoodsSpan := opentracing.GlobalTracer().StartSpan("save_order_goods", opentracing.ChildOf(parentSpan.Context()))

	if r := tx.CreateInBatches(orderGoods, 100); r.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
	}
	saveOrderGoodsSpan.Finish()

	//删除购物车的记录
	deleteShopCartSpan := opentracing.GlobalTracer().StartSpan("delete_shop_cart", opentracing.ChildOf(parentSpan.Context()))
	if r := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(model.ShoppingCart{}); r.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}
	deleteShopCartSpan.Finish()

	tx.Commit()

	return primitive.RollbackMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	if r := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); r.RowsAffected == 0 {
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}

func (o *OrderServer) CreateOrder(c context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单
			1、从购物车中获取到选中的商品
			2、商品的价格自己查询 - 访问商品服务（跨微服务）
			3、库存的扣减 - 访问库存服务（跨微服务）
			4、订单的基本信息表 - 订单的商品信息表
			5、从购物车中删除已购买的记录
	*/
	orderListener := OrderListener{}
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{}, producer.WithNameServer([]string{"192.168.220.128:9876"}))

	if err != nil {
		zap.S().Errorf("生成producer失败：%s\n", err.Error())
		return nil, err
	}

	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败：%s\n", err.Error())
		return nil, err
	}

	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SignerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}

	jsonString, _ := json.Marshal(order)

	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		fmt.Printf("发送失败：%s\n", err)
		return nil, status.Errorf(codes.Internal, "发送消息失败")
	}

	if res.State == primitive.CommitMessageState {
		return nil, status.Errorf(codes.Internal, "新建订单失败")
	}

	return &proto.OrderInfoResponse{Id: int32(orderListener.ID), OrderSn: order.OrderSn, Total: order.OrderMount}, nil
}

func (o *OrderServer) UpdateOrderStatus(c context.Context, req *proto.OrderStatus) (*proto.MyEmptyWithOrder, error) {
	if r := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	return &proto.MyEmptyWithOrder{}, nil
}
