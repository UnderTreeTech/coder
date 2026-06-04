package service

import (
	"user/internal/dao"
)

type Config struct {
	MaxIdPoolNum int
	// 是否开启debug调试验签接口
	DebugMode bool
}

type Service struct {
	cfg *Config
	dao dao.Dao
}

func New(d dao.Dao, cfg *Config) *Service {

	svc := &Service{
		dao: d,
		cfg: cfg,
	}

	return svc
}

func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) Ping() {

}
