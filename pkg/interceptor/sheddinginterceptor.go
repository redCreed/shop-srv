package interceptor

import (
	"context"
	"github.com/zeromicro/go-zero/core/load"
	"google.golang.org/grpc"
	"sync"
)
const serviceType = "rpc"
var (
	sheddingStat *load.SheddingStat
	lock  sync.Mutex
)

func SheddingInterceptor(shedder load.Shedder ) grpc.UnaryServerInterceptor   {
	ensureSheddingStat()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		sheddingStat.IncrementTotal()
		var promise load.Promise
		promise, err = shedder.Allow()
		if err != nil {
			sheddingStat.IncrementDrop()
			return
		}

		defer func() {

		}()


		return handler(ctx, req)
	}
}


func ensureSheddingStat() {
	lock.Lock()
	if sheddingStat == nil {
		sheddingStat = load.NewSheddingStat(serviceType)
	}
	lock.Unlock()
}
