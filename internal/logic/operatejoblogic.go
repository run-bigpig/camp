package logic

import (
	"context"
	"errors"

	"camp/internal/svc"
	"camp/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperateJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 操作任务
func NewOperateJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperateJobLogic {
	return &OperateJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperateJobLogic) OperateJob(req *types.OperateJobRequest) (resp *types.Empty, err error) {
	switch req.Operate {
	case "add":
		err = l.svcCtx.Job.AddJob(req)
		return nil, err
	case "delete":
		err = l.svcCtx.Job.DeleteJob(req.JobId)
		return nil, err
	}
	return nil, errors.New("无效操作")
}
