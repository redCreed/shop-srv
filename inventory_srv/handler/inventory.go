package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/dtm/dtmcli"
	"github.com/dtm-labs/dtm/dtmgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/model"
	"mxshop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

//设置库存
func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存， 如果我要更新库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

//库存详情数据
func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}



//func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
//	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
//	//数据库基本的一个应用场景：数据库事务
//	//并发情况之下 可能会出现超卖 1
//	client := goredislib.NewClient(&goredislib.Options{
//		Addr: "192.168.18.100:6379",
//	})
//	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
//	rs := redsync.New(pool)
//	tx := global.DB.Begin()
//	for _, goodInfo := range req.GoodsInfo {
//		var inv model.Inventory
//		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
//		//mysql悲观锁
//		//if result := global.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
//		//	tx.Rollback() //回滚之前的操作
//		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
//		//}
//
//		//分布式锁
//		if err := mutex.Lock(); err != nil {
//			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
//		}
//		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
//			tx.Rollback() //回滚之前的操作
//			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
//		}
//		//判断库存是否存在
//		if inv.Stocks < goodInfo.Num {
//			tx.Rollback() //回滚之前的操作
//			return nil, status.Errorf(codes.InvalidArgument, "库存不足")
//		}
//		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
//		inv.Stocks -= goodInfo.Num
//		//更新操作
//		tx.Save(inv)
//
//		if ok, err := mutex.Unlock(); !ok || err != nil {
//			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
//		}
//	}
//
//	tx.Commit()
//	return &emptypb.Empty{}, nil
//}
// 分布式需要考虑
//func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
//	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
//	//数据库基本的一个应用场景：数据库事务
//	//并发情况之下 可能会出现超卖 1
//	client := goredislib.NewClient(&goredislib.Options{
//		Addr: "192.168.0.104:6379",
//	})
//	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
//	rs := redsync.New(pool)
//	tx := global.DB.Begin()
//	for _, goodInfo := range req.GoodsInfo {
//		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
//		//分布式锁
//		if err := mutex.Lock(); err != nil {
//			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
//		}
//		var inv model.Inventory
//		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
//			tx.Rollback() //回滚之前的操作
//			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
//		}
//		//增加， 会出现数据不一致的问题 - 锁，分布式锁
//		inv.Stocks += goodInfo.Num
//		tx.Save(inv)
//		if ok, err := mutex.Unlock(); !ok || err != nil {
//			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
//		}
//	}
//
//	tx.Commit()
//	return &emptypb.Empty{}, nil
//}

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	empty := &emptypb.Empty{}
	barrer,err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return empty, status.Error(codes.Internal, err.Error())
	}
	db, err := global.DB.DB()
	if err != nil {
		return empty, status.Error(codes.Internal, err.Error())
	}
	if err := barrer.CallWithDB(db, func(tx *sql.Tx) error {
		stocksSql := fmt.Sprintf(`select stocks from inventory where id =1`)
		row, err := tx.Query(stocksSql)
		if err != nil {
			return err
		}
		stocks := 0
		for row.Next() {
			row.Scan(&stocks)
		}
		if stocks <2 {
			//需要回滚
			return  status.Error(codes.Aborted,  dtmcli.ResultFailure)
		}


		sql := fmt.Sprintf("update inventory set stocks=stocks-1 where id=1")
		_, err = tx.Exec(sql)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure
		return empty, status.Error(codes.Internal, err.Error())
	}

	return empty, nil
}

func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	empty := &emptypb.Empty{}
	barrer,err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return empty, status.Error(codes.Internal, err.Error())
	}
	db, err := global.DB.DB()
	if err != nil {
		return empty, status.Error(codes.Internal, err.Error())
	}
	if err := barrer.CallWithDB(db, func(tx *sql.Tx) error {
		sql := fmt.Sprintf("update inventory set stocks=2 where id=1")
		_, err := tx.Exec(sql)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure
		return empty, status.Error(codes.Internal, err.Error())
	}

	return empty, nil
}