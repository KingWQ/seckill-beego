package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"shop-seckill/models"
	"shop-seckill/services/mq"
	redisClient "shop-seckill/services/redis"
	"strconv"
	"time"
)

// Operations about Users
type OrderController struct {
	beego.Controller
}


//@router /order/seckill [*]
func (this *OrderController) Seckill() {
	skuId, _ := this.GetInt("sku_id")
	address := this.GetString("address")

	if skuId == 0 {
		this.Data["json"] = ReturnError(4001, "sku id不能为空")
		this.ServeJSON()
		return
	}
	if address == "" {
		this.Data["json"] = ReturnError(4002, "收货地址不能为空")
		this.ServeJSON()
		return
	}

	redisConn := redisClient.PoolConnect()
	defer redisConn.Close()

	//1：在缓存中校验库存
	stockKey := "go_stock:" + strconv.Itoa(skuId)
	stock, err := redis.Int(redisConn.Do("get", stockKey))
	if err != nil {
		this.Data["json"] = ReturnError(4003, "该商品不存在")
		this.ServeJSON()
		return
	}
	if stock < 1 {
		this.Data["json"] = ReturnError(4005, "该商品库存不足")
		this.ServeJSON()
		return
	}

	//2：在缓存中检验秒杀是否开始
	xt := ExpireTime{}
	expire, err := redis.String(redisConn.Do("get", "go_expire:1"))
	if err != nil {
		this.Data["json"] = ReturnError(4006, "商品秒杀时间缓存不存在")
		this.ServeJSON()
		return

	} else {
		json.Unmarshal([]byte(expire), &xt)
		start := xt.Start
		end := xt.End

		local, _ := time.LoadLocation("Local")
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, local)
		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, local)
		now := time.Now()
		if startTime.After(now) {
			this.Data["json"] = ReturnError(4007, "秒杀还未开始")
			this.ServeJSON()
			return
		}
		if endTime.Before(now) {
			this.Data["json"] = ReturnError(4007, "秒杀还未开始")
			this.ServeJSON()
			return
		}

	}

	//3：在缓存中检验该用户是否秒杀过
	userId := rand.Intn(99) + 1
	usrOrderKey := "go_user_order_" + strconv.Itoa(skuId) + ":" + strconv.Itoa(userId)
	order, err := redis.String(redisConn.Do("get", usrOrderKey))
	if err == nil {
		this.Data["json"] = ReturnError(4008, "该用户已经下过单"+order)
		this.ServeJSON()
		return
	}

	//4：秒杀入队
	msg := &MsgData{}
	msg.TaskName = "seckill_order"
	msg.UserId = userId
	msg.SkuId = skuId
	msg.Address = address
	msg.Time = time.Now().Format("2006-01-02 15:04:05")
	msgStr, _ := json.Marshal(msg)

	//go func() {
	mq.Publish("", "go_seckill", string(msgStr))
	//}()

	this.Data["json"] = ReturnSuccess(2000, "秒杀中", nil, 0)
	this.ServeJSON()
	return
}

type MsgData struct {
	TaskName string `json:"task_name"`
	UserId   int    `json:"user_id"`
	SkuId    int    `json:"sku_id"`
	Address  string `json:"address"`
	Time     string `json:"time"`
}

type ExpireTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// @router /cache/set [*]
func (this *OrderController) Set() {
	redisConn := redisClient.PoolConnect()
	defer redisConn.Close()
	redisConn.Do("set", "go_stock:1", 10)

	xt := &ExpireTime{}
	xt.Start = "2021-01-15 00:00:00"
	xt.End = "2021-02-15 00:00:00"
	data, _ := json.Marshal(xt)
	redisConn.Do("set", "go_expire:1", string(data))


	this.Ctx.WriteString("设置缓存成功")
}

// @router /cache/get [*]
func (this *OrderController) Get() {
	redisConn := redisClient.PoolConnect()
	defer redisConn.Close()

	xt := ExpireTime{}
	expire, err := redis.String(redisConn.Do("get", "go_expire:1"))
	if err != nil {
		this.Ctx.WriteString("商品秒杀时间缓存不存在")
	} else {

		err := json.Unmarshal([]byte(expire), &xt)
		fmt.Println(err)
		fmt.Println(xt.Start)
		fmt.Println(xt.End)

	}

	stock, err := redis.String(redisConn.Do("get", "go_stock:1"))
	if err != nil {
		this.Ctx.WriteString("商品sku 库存缓存key 不存在")
	} else {
		this.Ctx.WriteString(stock)
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timeStr)

	skuInfo, err := models.GetSkuInfo(1)
	if err == nil {
		fmt.Println(skuInfo.Stock)
	} else {
		fmt.Println("没有sku")
		fmt.Println(err)
	}

	//res, err := models.UpdateStock(1)
	//if err == nil {
	//	fmt.Println("mysql row affected nums: ", res)
	//
	//} else {
	//	fmt.Println("更新库存失败")
	//	fmt.Println(err)
	//}
	//
	//id,err := models.SaveOrder("golang3567891qwqw2",2, "广州")
	//if err == nil{
	//	err := models.SaveItem(id,1)
	//	if err != nil{
	//		fmt.Println(err)
	//	}
	//}else{
	//	fmt.Println(err)
	//}


}
