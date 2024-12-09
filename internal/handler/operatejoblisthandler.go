package handler

import (
	"net/http"

	"camp/internal/logic"
	"camp/internal/svc"
	"camp/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取任务列表
func OperateJobListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OperateJobListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewOperateJobListLogic(r.Context(), svcCtx)
		resp, err := l.OperateJobList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
