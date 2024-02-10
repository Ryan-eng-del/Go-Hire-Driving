// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.2
// - protoc             v4.25.0
// source: api/customer/customer.proto

package customer

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationCustomerGetVerifyCode = "/api.customer.Customer/GetVerifyCode"
const OperationCustomerLogin = "/api.customer.Customer/Login"
const OperationCustomerLogout = "/api.customer.Customer/Logout"

type CustomerHTTPServer interface {
	GetVerifyCode(context.Context, *GetVerifyCodeRequst) (*GetVerifyCodeResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
}

func RegisterCustomerHTTPServer(s *http.Server, srv CustomerHTTPServer) {
	r := s.Route("/")
	r.GET("/customer/get-verify-code/{telephone}", _Customer_GetVerifyCode0_HTTP_Handler(srv))
	r.POST("/customer/login", _Customer_Login0_HTTP_Handler(srv))
	r.GET("/customer/logout", _Customer_Logout0_HTTP_Handler(srv))
}

func _Customer_GetVerifyCode0_HTTP_Handler(srv CustomerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetVerifyCodeRequst
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCustomerGetVerifyCode)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetVerifyCode(ctx, req.(*GetVerifyCodeRequst))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetVerifyCodeResponse)
		return ctx.Result(200, reply)
	}
}

func _Customer_Login0_HTTP_Handler(srv CustomerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCustomerLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*LoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginResponse)
		return ctx.Result(200, reply)
	}
}

func _Customer_Logout0_HTTP_Handler(srv CustomerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LogoutRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCustomerLogout)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Logout(ctx, req.(*LogoutRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LogoutResponse)
		return ctx.Result(200, reply)
	}
}

type CustomerHTTPClient interface {
	GetVerifyCode(ctx context.Context, req *GetVerifyCodeRequst, opts ...http.CallOption) (rsp *GetVerifyCodeResponse, err error)
	Login(ctx context.Context, req *LoginRequest, opts ...http.CallOption) (rsp *LoginResponse, err error)
	Logout(ctx context.Context, req *LogoutRequest, opts ...http.CallOption) (rsp *LogoutResponse, err error)
}

type CustomerHTTPClientImpl struct {
	cc *http.Client
}

func NewCustomerHTTPClient(client *http.Client) CustomerHTTPClient {
	return &CustomerHTTPClientImpl{client}
}

func (c *CustomerHTTPClientImpl) GetVerifyCode(ctx context.Context, in *GetVerifyCodeRequst, opts ...http.CallOption) (*GetVerifyCodeResponse, error) {
	var out GetVerifyCodeResponse
	pattern := "/customer/get-verify-code/{telephone}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCustomerGetVerifyCode))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CustomerHTTPClientImpl) Login(ctx context.Context, in *LoginRequest, opts ...http.CallOption) (*LoginResponse, error) {
	var out LoginResponse
	pattern := "/customer/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationCustomerLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CustomerHTTPClientImpl) Logout(ctx context.Context, in *LogoutRequest, opts ...http.CallOption) (*LogoutResponse, error) {
	var out LogoutResponse
	pattern := "/customer/logout"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCustomerLogout))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}