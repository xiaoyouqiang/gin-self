package test_service

import (
	"context"
	"sync"

	"gin-self/model/mysql/test1_model"
)

//service 结构体 必须实现 Service 接口
var _ Service = (*service)(nil)
var instance *service
var once sync.Once

type service struct {
	ctx context.Context
}

//Service 服务模板 需要实现其中的方法
type Service interface {
	i()

	Create(name string) (int32, error)
}

//service 结构体 必须在本包中实现 Service 接口
func (s *service) i() {}

//GetInstance 获取该服务实例
func GetInstance(ctx context.Context) Service {
	once.Do(func() {
		instance = &service{
			ctx: ctx,
		}
	})
	return instance
}

func (s *service) Create(name string) (int32, error) {

	data := test1_model.Test1{
		Name: name,
	}

	query := test1_model.NewQueryBuilder()

	id, _ := query.Create(data)

	return id, nil
}
