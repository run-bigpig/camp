package logic

import (
	"camp/internal/config"
	"camp/internal/svc"
)

func getConfig(svcCtx *svc.ServiceContext, platform int8) *config.CampConfig {
	switch platform {
	case 0:
		return svcCtx.Config.Nms
	case 1:
		return svcCtx.Config.My
	default:
		return svcCtx.Config.Nms
	}
}
