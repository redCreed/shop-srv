package load

import "time"

type SheddingStat struct {
	name  string
	total int64
	pass  int64
	drop  int64
}

//用于保持当前请求的的通过和拒绝的数量，并上报到prometheus
type snapshot struct {
	total int64
	pass  int64
	drop  int64
}

//根据服务名称
func NewSheddingStat(name string) *SheddingStat {
	st := &SheddingStat{
		name:  name,
		total: 0,
		pass:  0,
		drop:  0,
	}

	return st
}

func (s *SheddingStat)run()  {
	//每一分钟记录日志
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range  ticker.C {

	}
}