package http

import (
	"net/http"
	"strconv"

	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"github.com/jacknotes/restful-api-demo/app/host"
	"github.com/julienschmidt/httprouter"
)

func (h *handler) CreateHost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := request.ReadBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	h.log.Debugf("receive body: %s", string(body))
	response.Success(w, "ok")
}

func (h *handler) QueryHost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		pageSize   = 20
		pageNumber = 1
	)

	// 获取请求URL中的所有参数
	qs := r.URL.Query()

	psStr := qs.Get("page_size")
	if psStr != "" {
		pageSize, _ = strconv.Atoi(psStr)
	}

	pnStr := qs.Get("page_number")
	if pnStr != "" {
		pageNumber, _ = strconv.Atoi(pnStr)
	}

	// 构建查询主机请求body值
	req := &host.QueryHostRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}

	// 调用业务层逻辑函数，查询主机列表
	set, err := h.host.QueryHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	// 传递给success的是一个切片对象
	// success 会把这个切片对象序列化为一个JSON
	// 补充返回的数据
	// 封装这个函数的原因，是为了标准化输出格式
	response.Success(w, set)

}
