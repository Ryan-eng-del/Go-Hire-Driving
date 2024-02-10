package server

import (
	"context"
	"customer/internal/service"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

const (

	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized(reason, "JWT token is missing")
	ErrMissingKeyFunc         = errors.Unauthorized(reason, "keyFunc is missing")
	ErrTokenInvalid           = errors.Unauthorized(reason, "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized(reason, "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized(reason, "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized(reason, "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized(reason, "Wrong context for middleware")
	ErrNeedTokenProvider      = errors.Unauthorized(reason, "Token provider is missing")
	ErrSignToken              = errors.Unauthorized(reason, "Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.Unauthorized(reason, "Can not get key while signing token")
)

// Server is a server auth middleware. Check the token and extract the info from token.
func ValidateJWT(customerService *service.CustomerService) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			claims, ok := jwt.FromContext(ctx)
			if !ok {
				return nil, ErrMissingJwtToken
			}
			claimsMap := claims.(jwtv4.MapClaims)
			userId := claimsMap["jti"]	
			token, err := customerService.CustomData.GetToken(userId)
			if token == "" || err != nil {
				return nil, ErrMissingJwtToken
			}

			if header, ok := transport.FromServerContext(ctx); ok {
				auths := strings.SplitN(header.RequestHeader().Get(authorizationKey), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
					return nil, ErrMissingJwtToken
				}
				jwtToken := auths[1]
				if jwtToken != token {
					return nil, ErrUnSupportSigningMethod
				}
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

func AllowCORS() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				ht := tr.(http.Transporter)
				replyHeader := ht.ReplyHeader()
				replyHeader.Set("Access-Control-Allow-Origin", "*")
				replyHeader.Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,HEAD,OPTIONS,PATCH,PUT,DELETE")
				replyHeader.Set("Access-Control-Allow-Headers","Content-Type, X-Requested-With, Authorization, Access-control-Allow-Credentials, User-Agent")
				replyHeader.Set("Access-control-Allow-Credentials", "true")
			}
			return handler(ctx, req)
		}
	}
}
