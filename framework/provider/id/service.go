package id

import "github.com/rs/xid"

type LadeIDService struct {
}

func NewLadeIDService(params ...interface{}) (interface{}, error) {
	return &LadeIDService{}, nil
}

func (s *LadeIDService) NewID() string {
	return xid.New().String()
}
