package handler

import (
	"net/http"

	"camp/internal/logic"
	"camp/internal/svc"
	"camp/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 操作任务
func OperateJobHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OperateJobRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewOperateJobLogic(r.Context(), svcCtx)
		resp, err := l.OperateJob(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
