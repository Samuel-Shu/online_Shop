package handler

import (
	"context"
	_ "crypto/md5"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"online_Shop/user_srv/global"
	"online_Shop/user_srv/model"
	"online_Shop/user_srv/proto"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Model2Resp(user model.User) proto.UserInfoResponse {
	UserInfoResp := proto.UserInfoResponse{
		Password: user.Password,
		Email:    user.Email,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     uint32(user.Role),
	}
	// 在grpc的message字段中有默认值，不能直接赋值nil，易出错
	if user.Birthday != nil {
		UserInfoResp.Birthday = uint64(user.Birthday.Unix())
	}

	return UserInfoResp
}

// Paginate 分页查询
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (u *UserServer) GetUserList(ctx context.Context, in *proto.PageInfo) (*proto.UserListResponse, error) {
	//获取用戶信息
	var users []model.User
	tx := global.DB.Find(&users)

	if tx.Error != nil {
		return nil, tx.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(tx.RowsAffected)
	//1、使用作用域函数Scopes复用分页逻辑
	global.DB.Scopes(Paginate(int(in.Pn), int(in.PSize))).Find(&users)

	for _, user := range users {
		userInfoRes := Model2Resp(user)
		rsp.Data = append(rsp.Data, &userInfoRes)
	}
	return rsp, nil
}

func (u *UserServer) GetUserByEmail(ctx context.Context, in *proto.EmailRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	tx := global.DB.Where("email = ?", in.Email).First(&user)
	if tx.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	resp := Model2Resp(user)
	return &resp, nil
}

func (u *UserServer) GetUserById(ctx context.Context, in *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	tx := global.DB.Where("id = ?", in.Id).First(&user)
	if tx.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	resp := Model2Resp(user)
	return &resp, nil
}

func (u *UserServer) CreateUser(ctx context.Context, in *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	//新建用户
	//1、先查询用户是否已经存在：若存在，则返回注册错误状态码；若不存在，则继续注册
	var user model.User
	tx := global.DB.Where("email = ?", in.Email).First(&user)
	if tx.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已经存在")
	}
	user.Email = in.Email
	user.NickName = in.NickName

	//密码加密，出于安全性考虑，这里使用md5加盐加密处理
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(in.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	tx = global.DB.Create(&user)
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, tx.Error.Error())
	}

	userInfoRes := Model2Resp(user)
	return &userInfoRes, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, in *proto.UpdateUserInfo) (*proto.MyEmpty, error) {
	//修改用户信息
	//1、只有用户已经存在才能修改
	var user model.User
	tx := global.DB.First(&user, in.Id)
	if tx.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthday := time.Unix(int64(in.Birthday), 0)
	user.NickName = in.NickName
	user.Gender = in.Gender
	user.Birthday = &birthday

	tx = global.DB.Save(user)
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, tx.Error.Error())
	}

	return nil, nil
}

func (u *UserServer) CheckPassword(ctx context.Context, in *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	//校验密码
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(in.EncryptedPassword, "$")
	check := password.Verify(in.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}
