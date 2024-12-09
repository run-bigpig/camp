package logic

import (
	"camp/internal/config"
	"context"
	"errors"

	"camp/internal/svc"
	"camp/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新token
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest) (resp *types.RefreshTokenResponse, err error) {
	switch req.Platform {
	case 0:
		config.RefreshToken(l.svcCtx.Config, req.Token, "")
		return &types.RefreshTokenResponse{
			Token: l.svcCtx.Config.Nms.Token,
		}, nil
	case 1:
		config.RefreshToken(l.svcCtx.Config, "", req.Token)
		return &types.RefreshTokenResponse{
			Token: l.svcCtx.Config.My.Token,
		}, nil
	default:
		return nil, errors.New("平台参数错误")
	}
}
