package host

import (
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

// 为了后期做资源解锁分两张表存储，<ip> ---> Host, IP, SLB, Redis, Mysql
type Host struct {
	// hash字段，在两个表插入数据时，当只需要插入一个表，另外一个表不需要插入，为了不插入另外一个表面做的hash校验
	ResourceHash string //         `json:"resource_hash"`
	DescribeHash string //        `json:"describe_hash"`
	*Resource
	*Describe
}

func NewDefaultHost() *Host {
	return &Host{
		Resource: &Resource{},
		Describe: &Describe{},
	}
}

func (h *Host) Validate() error {
	return validate.Struct(h)
}

type Vendor int

const (
	ALI_CLOUD Vendor = iota
	TX_CLOUD
	HW_CLOUD
)

// 主机元数据信息
type Resource struct {
	Id     string `json:"id" validate:"required"`     // 全局唯一Id
	Vendor Vendor `json:"vendor" validate:"required"` // 厂商
	Region string `json:"region" validate:"required"` // 地域
	Zone   string `json:"zone"`                       // 区域
	// 使用的13位时间戳
	// 为什么不用数据库Datetime，如果使用数据库的时间，数据库会默认加上时区
	// 后端使用时间戳，不加时区，都由前端加上时区进行展示
	CreateAt    int64             `json:"create_at" validate:"required"`  // 创建时间
	ExpireAt    int64             `json:"expire_at"`                      // 过期时间
	Category    string            `json:"category"`                       // 种类
	Type        string            `json:"type"`                           // 规格
	InstanceId  string            `json:"instance_id"`                    // 实例ID
	Name        string            `json:"name" validate:"required"`       // 名称
	Description string            `json:"description"`                    // 描述
	Status      string            `json:"status" validate:"required"`     // 服务商中的状态
	Tags        map[string]string `json:"tags"`                           // 标签
	UpdateAt    int64             `json:"update_at"`                      // 更新时间
	SyncAt      int64             `json:"sync_at"`                        // 同步时间
	SyncAccount string            `json:"sync_accout"`                    // 同步的账号
	PublicIP    string            `json:"public_ip"`                      // 公网IP
	PrivateIP   string            `json:"private_ip" validate:"required"` // 内网IP
	PayType     string            `json:"pay_type"`                       // 实例付费方式
}

// 主机具体信息
type Describe struct {
	CPU                     int    `json:"cpu" validate:"required"`    // 核数
	Memory                  int    `json:"memory" validate:"required"` // 内存
	GPUAmount               int    `json:"gpu_amount"`                 // GPU数量
	GPUSpec                 string `json:"gpu_spec"`                   // GPU类型
	OSType                  string `json:"os_type"`                    // 操作系统类型，分为Windows和Linux
	OSName                  string `json:"os_name"`                    // 操作系统名称
	SerialNumber            string `json:"serial_number"`              // 序列号
	ImageID                 string `json:"image_id"`                   // 镜像ID
	InternetMaxBandwidthOut int    `json:"internet_max_bandwidth_out"` // 公网出带宽最大值，单位为 Mbps
	InternetMaxBandwidthIn  int    `json:"internet_max_bandwidth_in"`  // 公网入带宽最大值，单位为 Mbps
	KeyPairName             string `json:"key_pair_name"`              // 秘钥对名称
	SecurityGroups          string `json:"security_groups"`            // 安全组  采用逗号分隔
}

// 查询主机列表信息 返回参数
type Set struct {
	Total int64
	Items []*Host
}

func NewSet() *Set {
	return &Set{
		// 初始化切片，否则会加对象加不进去的
		Items: []*Host{},
	}
}

func (s *Set) Add(item *Host) {
	s.Items = append(s.Items, item)
}
