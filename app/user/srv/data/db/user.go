package db

import (
	"context"
	"go-5m3Micro/app/user/srv/data"
	metav1 "go-5m3Micro/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *users {
	return &users{db: db}
}

func (u users) List(ctx context.Context, req *metav1.ListMeta) (*data.UserDoList, error) {
	//TO DO SQL
	return nil, nil
}

// 接口验证
var _ data.UserStore = (*users)(nil)
