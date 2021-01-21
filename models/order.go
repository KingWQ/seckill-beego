package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Orders struct{
	Id int
	No string
	UserId int
	Address string
	TotalAmount float32
	CreatedAt string
}



func init(){
	orm.RegisterModel(new(Orders))
}

func SaveOrder( no string, userId int, address string) (int64, error){
	var (
		err error
		order Orders
	)
	o := orm.NewOrm()
	order.No = no
	order.Address = address
	order.UserId = userId
	order.TotalAmount= 1.0
	order.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	id,err := o.Insert(&order)

	return id,err
}

