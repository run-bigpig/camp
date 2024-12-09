package logic

import (
	"camp/internal/config"
	"camp/internal/consts"
	"camp/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"sync"

	"camp/internal/svc"
	"camp/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoomListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取房间列表
func NewGetRoomListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoomListLogic {
	return &GetRoomListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoomListLogic) GetRoomList(req *types.GetRoomListRequst) (resp *types.GetRoomListResponse, err error) {
	c := getConfig(l.svcCtx, req.Platform)
	resultChannel := make(chan []*types.Room, len(c.ListId))
	list := make([]*types.Room, 0)
	wg := sync.WaitGroup{}
	for _, id := range c.ListId {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			roomList, err := l.requestList(c, req.RoomStartDate, req.RoomEndDate, id)
			if err != nil {
				l.Logger.Errorf("请求房间列表失败：%s", err.Error())
				return
			}
			resultChannel <- roomList
		}(id)
	}
	wg.Wait()
	close(resultChannel)
	for v := range resultChannel {
		list = append(list, v...)
	}
	return &types.GetRoomListResponse{
		RoomList: list,
	}, nil
}

func (l *GetRoomListLogic) requestList(config *config.CampConfig, startDate, endDate, id string) ([]*types.Room, error) {
	r := &types.ListRoomRequest{
		NtwNum:           config.NtwNum,
		OpenID:           config.OpenId,
		PromotionChannel: 3,
		RoomStartDate:    startDate,
		RoomEndDate:      endDate,
		RoomPageNum:      1,
		RoomPageSize:     50,
		ID:               id,
	}
	header := map[string]string{"uid": config.Uid, "token": config.Token, "User-Agent": consts.UserAgent, "Referer": consts.Referer}
	data, err := utils.SendRequest(consts.ListUrl, header, r)
	if err != nil {
		return nil, err
	}
	var listResp types.ListRoomResponse
	err = json.Unmarshal(data, &listResp)
	if err != nil {
		return nil, err
	}
	if listResp.Code != "1" {
		return nil, errors.New(listResp.Msg)
	}
	list := make([]*types.Room, 0, len(listResp.Data.List))
	for _, v := range listResp.Data.List {
		list = append(list, &types.Room{
			RoomId:   v.RoomTypeID,
			RoomName: v.Name,
		})
	}
	return list, nil
}
