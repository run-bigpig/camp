package logic

import (
	"context"

	"camp/internal/svc"
	"camp/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperateJobListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取任务列表
func NewOperateJobListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperateJobListLogic {
	return &OperateJobListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperateJobListLogic) OperateJobList(req *types.OperateJobListRequest) (resp *types.OperateJobListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
