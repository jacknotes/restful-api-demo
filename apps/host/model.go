package host

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
)

var (
	validate = validator.New()
)

func NewDefaultHost() *Host {
	return &Host{
		Resource: &Resource{
			CreateAt: time.Second.Milliseconds(),
		},
		Describe: &Describe{},
	}
}

func (h *Host) Validate() error {
	return validate.Struct(h)
}

func (h *Host) Patch(res *Resource, desc *Describe) error {
	if res != nil {
		// 合并，patch方法常用，将res的值覆盖h.Resource的值，其余值保留
		err := mergo.MergeWithOverwrite(h.Resource, res)
		if err != nil {
			return err
		}
	}

	if desc != nil {
		err := mergo.MergeWithOverwrite(h.Describe, desc)
		if err != nil {
			return err
		}
	}
	h.Resource.UpdateAt = time.Now().UnixMilli()
	return nil
}

// go 1.17允许获取毫秒
func (h *Host) Update(res *Resource, desc *Describe) {
	h.Resource = res
	h.Describe = desc
	h.Resource.UpdateAt = time.Now().UnixMilli()
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
