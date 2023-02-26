package http

import (
	"net/http"
	"strconv"

	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"github.com/jacknotes/restful-api-demo/apps/host"
)

func (h *handler) CreateHost(w http.ResponseWriter, r *http.Request) {
	// 用于接收前端传递过来的数据，为了做数据格式转换
	payload := &struct {
		*host.Resource
		*host.Describe
	}{
		Resource: &host.Resource{},
		Describe: &host.Describe{},
	}

	// 传入req初始化的结构体给request.GetDataFromRequest方法进行反序列化（JSON -> Struct）
	if err := request.GetDataFromRequest(r, payload); err != nil {
		response.Failed(w, err)
		return
	}

	req := host.NewDefaultHost()
	req.Resource = payload.Resource
	req.Describe = payload.Describe

	// 1. Context一定要传，如果用户中断了请求，你的后端逻辑也会跟着中断
	// 2. req: 通过http协议传递进来
	// 3. h.host通过Config()中h.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)已经变成grpc对象Service
	ins, err := h.host.CreateHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) QueryHost(w http.ResponseWriter, r *http.Request) {
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
		PageSize:   int64(pageSize),
		PageNumber: int64(pageNumber),
		Keywords:   qs.Get("keywords"),
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

// 查询主机详情
// httprouter.Params保存着我们的路径查询参数
func (h *handler) DescribeHost(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := &host.DescribeHostRequest{
		// 从路径获取:id的值
		Id: ctx.PS.ByName("id"),
	}

	// 调用业务层逻辑函数，查询主机列表
	host, err := h.host.DescribeHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	// 传递给success的是一个切片对象
	// success 会把这个切片对象序列化为一个JSON
	// 补充返回的数据
	// 封装这个函数的原因，是为了标准化输出格式
	response.Success(w, host)

}

// 更新单个主机
func (h *handler) UpdateHost(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := host.NewPutUpdateHostRequest()

	// 传入req初始化的结构体给request.GetDataFromRequest方法进行反序列化（JSON -> Struct）
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	// 查询的是用户传入"/hosts/:id"中的id的值，并且将原Id赋值给req.Id
	req.Resource.Id = ctx.PS.ByName("id")
	req.Describe.ResourceId = req.Resource.Id

	host, err := h.host.UpdateHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	// 传递给success的是一个切片对象
	// success 会把这个切片对象序列化为一个JSON
	// 补充返回的数据
	// 封装这个函数的原因，是为了标准化输出格式
	response.Success(w, host)
}

func (h *handler) PatchHost(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := host.NewPatchUpdateHostRequest()

	// 传入req初始化的结构体给request.GetDataFromRequest方法进行反序列化（JSON -> Struct）
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.Resource.Id = ctx.PS.ByName("id")
	req.Describe.ResourceId = req.Resource.Id

	host, err := h.host.UpdateHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	// 传递给success的是一个切片对象
	// success 会把这个切片对象序列化为一个JSON
	// 补充返回的数据
	// 封装这个函数的原因，是为了标准化输出格式
	response.Success(w, host)
}

// 删除主机
func (h *handler) DeleteHost(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := &host.DeleteHostRequest{
		Id: ctx.PS.ByName("id"),
	}

	host, err := h.host.DeleteHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, host)
}
