package job

import (
	"camp/internal/logic"
	"camp/internal/svc"
	"camp/internal/types"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

var platformMap = map[int8]string{
	0: "柠檬树",
	1: "牧云",
}

type Job struct {
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	logger      logx.Logger
	commitOrder *logic.CommitOrderLogic
}

func NewJob(ctx context.Context, svcCtx *svc.ServiceContext) *Job {
	return &Job{
		ctx:         ctx,
		svcCtx:      svcCtx,
		logger:      logx.WithContext(ctx),
		commitOrder: logic.NewCommitOrderLogic(ctx, svcCtx),
	}
}

func (l *Job) Run() {
	l.runMy()
	l.runNms()
}

func (l *Job) runNms() {
	if !l.svcCtx.Config.Nms.Job.Run {
		return
	}
	roomList := make([]*types.CommitRoom, 0)
	for _, roomId := range l.svcCtx.Config.Nms.Job.RoomId {
		for _, roomDate := range l.svcCtx.Config.Nms.Job.RoomDate {
			dates := strings.Split(roomDate, "#")
			roomList = append(roomList, &types.CommitRoom{
				RoomId:        roomId,
				RoomName:      roomId,
				RoomStartDate: dates[0],
				RoomEndDate:   dates[1],
			})
		}
	}
	go l.commit(0, roomList)
}

func (l *Job) runMy() {
	if !l.svcCtx.Config.My.Job.Run {
		return
	}
	roomList := make([]*types.CommitRoom, 0)
	for _, roomId := range l.svcCtx.Config.My.Job.RoomId {
		for _, roomDate := range l.svcCtx.Config.My.Job.RoomDate {
			dates := strings.Split(roomDate, "#")
			roomList = append(roomList, &types.CommitRoom{
				RoomId:        roomId,
				RoomName:      roomId,
				RoomStartDate: dates[0],
				RoomEndDate:   dates[1],
			})
		}
	}
	go l.commit(1, roomList)
}

func (l *Job) commit(platform int8, roomList []*types.CommitRoom) {
	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case t := <-ticker.C:
			l.logger.Infof("%s %s开始提交订单", platformMap[platform], t.Format("2006-01-02 15:04:05"))
			dataList, err := l.commitOrder.CommitOrder(&types.CommitOrderRequest{
				Platform:       platform,
				CommitRoomList: roomList,
			})
			if err != nil {
				l.logger.Errorf("%s %s提交订单失败：%s", platformMap[platform], t.Format("2006-01-02 15:04:05"), err.Error())
				continue
			}
			for _, v := range dataList.Result {
				if v.Success {
					fmt.Printf("%s %s提交订单成功：%s\n", platformMap[platform], t.Format("2006-01-02 15:04:05"), v.RoomName)
				}
			}
		}
	}
}
