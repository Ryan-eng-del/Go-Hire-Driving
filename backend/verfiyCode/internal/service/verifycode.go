package service

import (
	"context"
	"math/rand"
	"strings"

	pb "verfiyCode/api/verifyCode"
)

type VerifyCodeService struct {
	pb.UnimplementedVerifyCodeServer
}

func NewVerifyCodeService() *VerifyCodeService {
	return &VerifyCodeService{}
}


func (s *VerifyCodeService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeReply, error) {
	return &pb.GetVerifyCodeReply{
	Code: RandCode(int(req.Length), req.Type),
	}, nil
}


type RandCodeConst string

const (
	DIGIT RandCodeConst = "0123456789"
	LETER RandCodeConst = "abcdfghigklmnopqsrtuvwxyz"
	MIXIN RandCodeConst = "0cdp1uvwx2abq34srtyz56klm7gh8figno9"
)

func RandCode(l int, t pb.TYPE) string {
	charCode := ""
	switch t {
	case pb.TYPE_DEFAULT:
		fallthrough
	case pb.TYPE_DIGIT:
		charCode = randCode(string(DIGIT), l, 4)
	case pb.TYPE_LETER:
		charCode = randCode(string(LETER), l, 5)
	case pb.TYPE_MIXED:
		charCode = randCode(string(MIXIN), l, 6)
	default:
		charCode = randCode(string(DIGIT), l, 4)
	}
	return charCode
}

func randCode(char string, l int, bitSize int) string{
	idxMask := 1 << bitSize - 1
	idxMax := 63 / bitSize
	sb := strings.Builder{}
	sb.Grow(l)

	for i, cache, remain := l - 1, rand.Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), idxMax
		}
		if idx := int(cache & int64(idxMask)); idx < len(char) {
			sb.WriteByte(char[idx])
			i--
		}
		cache >>= int64(bitSize)
		remain--
	}
	return sb.String()
}

