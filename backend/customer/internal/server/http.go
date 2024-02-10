package server

import (
	"context"
	"customer/api/customer"
	v1 "customer/api/helloworld/v1"
	"customer/internal/biz"
	"customer/internal/conf"
	"customer/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	kratosJwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, customerService *service.CustomerService,logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			selector.Server(AllowCORS(),kratosJwt.Server(func(token *jwt.Token) (interface{}, error) {
				return []byte(biz.Secret), nil
			}), ValidateJWT(customerService)).Match(func (ctx context.Context, operation string) bool {
				noJWT := map[string]struct{}{
					"/api.customer.Customer/Login": {},
					"/api.customer.Customer/GetVerifyCode": {},
				}
				_, exits := noJWT[operation]
				if exits {
					return false
				} else {
					return true
				}
			}).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	customer.RegisterCustomerHTTPServer(srv, customerService)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
