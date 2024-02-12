// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.2
// - protoc             v4.25.0
// source: api/driver/driver.proto

package driver

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

const OperationDriverGetVerifyCode = "/api.driver.Driver/GetVerifyCode"

type DriverHTTPServer interface {
	GetVerifyCode(context.Context, *GetVerifyCodeRequest) (*GetVerifyCodeResponse, error)
}

func RegisterDriverHTTPServer(s *http.Server, srv DriverHTTPServer) {
	r := s.Route("/")
	r.POST("/driver/get-verify-code", _Driver_GetVerifyCode0_HTTP_Handler(srv))
}

func _Driver_GetVerifyCode0_HTTP_Handler(srv DriverHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetVerifyCodeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDriverGetVerifyCode)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetVerifyCode(ctx, req.(*GetVerifyCodeRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetVerifyCodeResponse)
		return ctx.Result(200, reply)
	}
}

type DriverHTTPClient interface {
	GetVerifyCode(ctx context.Context, req *GetVerifyCodeRequest, opts ...http.CallOption) (rsp *GetVerifyCodeResponse, err error)
}

type DriverHTTPClientImpl struct {
	cc *http.Client
}

func NewDriverHTTPClient(client *http.Client) DriverHTTPClient {
	return &DriverHTTPClientImpl{client}
}

func (c *DriverHTTPClientImpl) GetVerifyCode(ctx context.Context, in *GetVerifyCodeRequest, opts ...http.CallOption) (*GetVerifyCodeResponse, error) {
	var out GetVerifyCodeResponse
	pattern := "/driver/get-verify-code"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDriverGetVerifyCode))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}