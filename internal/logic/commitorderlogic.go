package logic

import (
	"camp/internal/config"
	"camp/internal/consts"
	"camp/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"sync"
	"time"

	"camp/internal/svc"
	"camp/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommitOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 提交订单
func NewCommitOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommitOrderLogic {
	return &CommitOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommitOrderLogic) CommitOrder(req *types.CommitOrderRequest) (resp *types.CommitOrderResponse, err error) {
	c := getConfig(l.svcCtx, req.Platform)
	resultList := make(chan *types.CommitRoomResult, len(req.CommitRoomList))
	wg := sync.WaitGroup{}
	for _, room := range req.CommitRoomList {
		wg.Add(1)
		go func(room *types.CommitRoom) {
			defer wg.Done()
			res := &types.CommitRoomResult{
				RoomId:   room.RoomId,
				RoomName: room.RoomName,
				Success:  true,
			}
			err := l.requestCommit(c, room)
			if err != nil {
				res.Success = false
				res.FailReason = err.Error()
				resultList <- res
				return
			}
			resultList <- res
		}(room)
	}
	wg.Wait()
	close(resultList)
	var resultListFinal []*types.CommitRoomResult
	for v := range resultList {
		resultListFinal = append(resultListFinal, v)
	}
	return &types.CommitOrderResponse{
		Result: resultListFinal,
	}, nil
}

func (l *CommitOrderLogic) requestCommit(config *config.CampConfig, room *types.CommitRoom) error {
	r := &types.CommitRequest{
		AppType:                 0,
		BookingType:             0,
		BoolPublish:             true,
		BoolUseContinueLiveInst: false,
		BoolUseIntegral:         false,
		Count:                   1,
		CustomerName:            config.CustomerName,
		CustomerPhone:           config.CustomerPhone,
		DiscountInstList:        []interface{}{},
		NtwNum:                  config.NtwNum,
		OpenID:                  config.OpenId,
		PromotionChannel:        3,
		Remark:                  "",
		RoomTypeID:              room.RoomId,
		StartDate:               room.RoomStartDate,
		EndDate:                 room.RoomEndDate,
	}
	header := map[string]string{
		"Referer":    consts.Referer,
		"User-Agent": consts.UserAgent,
		"token":      config.Token,
		"uid":        config.Uid,
	}
	randNum := rand.IntN(9)
	time.Sleep(time.Duration(randNum) * time.Second)
	data, err := utils.SendRequest(consts.CommitUrl, header, r)
	if err != nil {
		return err
	}
	var commitResp types.CommitResponse
	err = json.Unmarshal(data, &commitResp)
	if err != nil {
		return err
	}
	if commitResp.Code != "1" {
		return errors.New(commitResp.Msg)
	}
	return nil
}
