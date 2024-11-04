package logic

import (
	"camp/internal/config"
	"camp/internal/svc"
)

const (
	UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.53(0x1800352c) NetType/WIFI Language/zh_CN"
	Referer   = "https://servicewechat.com/wxae20944101069a31/1/page-frame.html"
	ListUrl   = "https://dingdandao.com/v2/micro/web/micro/goods/public/room/page"
	CommitUrl = "https://dingdandao.com/v2/micro/web/micro/order/V2/commitAccomOrder"
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
