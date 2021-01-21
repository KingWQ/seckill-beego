package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"shop-seckill/controllers"
	"shop-seckill/services/mq"
	redisClient "shop-seckill/services/redis"
	"strconv"
	"time"
)


func main(){
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	defaultdb := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", defaultdb, 30, 30)
	mq.Consumer("", "go_seckill", callback)
}

func callback(s string){
	type Data struct{
		TaskName string `json:"task_name"`
		UserId   int    `json:"user_id"`
		SkuId    int    `json:"sku_id"`
		Address  string `json:"address"`
		Time     string `json:"time"`
	}
	var data Data
	err := json.Unmarshal([]byte(s), &data)

	if err == nil{

		redisConn := redisClient.PoolConnect()
		defer redisConn.Close()

		//校验库存
		stockKey := "go_stock:" + strconv.Itoa(data.SkuId)

		stock, _ := redis.Int(redisConn.Do("get", stockKey))
		if stock < 1{
					fmt.Println("库存不足")
					return
		}
		//skuInfo,err := models.GetSkuInfo(data.SkuId)
		//if err == nil{
		//	fmt.Println(skuInfo.Stock)
		//	if skuInfo.Stock < 1 {
		//		redisConn.Do("del", "go_stock:"+strconv.Itoa(data.SkuId))
		//		fmt.Println("库存不足")
		//		return
		//	}
		//
		//}else{
		//	fmt.Println("没有sku"+strconv.Itoa(data.SkuId))
		//	return
		//}


		no := controllers.RandString(21)
		now := time.Now().Format("2006-01-02 15:04:05")

		o := orm.NewOrm()
		o.Begin()

		res1,err1 := o.Raw("update product_skus set stock = stock - 1 where stock >= 1 and id = ?", data.SkuId).Exec()
		updateRow,_ := res1.RowsAffected()

		res2,err2 := o.Raw("INSERT INTO orders (`no`, `user_id`, `address`, `total_amount`,`created_at`) VALUES (?, ?, ?, ?,?)",
			no, data.UserId, data.Address, 1.0, now).Exec()
		orderId,err4 := res2.LastInsertId()

		_,err3 := o.Raw("INSERT INTO order_items (`price`, `product_id`, `product_sku_id`, `amount`,`order_id`) VALUES (?, ?, ?, ?,?)",
			1.0, 1, data.SkuId, 1, orderId).Exec()


		if(err1 != nil  || updateRow<1 || err2 != nil || err4 != nil || err3 !=nil ){
			fmt.Println("事务回滚")
			o.Rollback()
		}else{
			o.Commit()
			//redis库存更新
			stockKey := "go_stock:" + strconv.Itoa(data.SkuId)
			redisConn.Do("decr", stockKey)

			//下单成功用户存入缓存
			usrOrderKey := "go_user_order_" + strconv.Itoa(data.SkuId) + ":" + strconv.Itoa(data.UserId)
			redisConn.Do("set", usrOrderKey, orderId)
		}

	}
	fmt.Printf("msg is :%s\n", s)
}
type Orders struct{
	No string
	UserId int
	Address string
	TotalAmount float32
	CreatedAt string
}

type OrderItems struct{
	OrderId int64
	ProductId int
	ProductSkuId int
	Amount int
	Price float32
}