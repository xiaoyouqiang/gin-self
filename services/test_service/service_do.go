package test_service

import (
	"gin-self/model/mysql/test1_model"
)

func (s *service) Create(name string) (int32, error) {

	data := test1_model.Test1{
		Name: name,
	}

	query := test1_model.NewQueryBuilder()

	id, _ := query.Create(data)

	return id, nil
}
