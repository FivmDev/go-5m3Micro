package data

import (
	"context"
	metav1 "go-5m3Micro/pkg/common/meta/v1"
)

type UserDo struct {
	Name string `json:"name"`
}

type UserDoList struct {
	TotalNumber int32     `json:"total_number"`
	Items       []*UserDo `json:"items,omitempty"`
}

type UserStore interface {
	List(ctx context.Context, req *metav1.ListMeta) (*UserDoList, error)
}
