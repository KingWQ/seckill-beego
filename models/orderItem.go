package models

import (
	"github.com/astaxie/beego/orm"
)

type OrderItems struct{
	Id int
	OrderId int64
	ProductId int
	ProductSkuId int
	Amount int
	Price float32
}



func init(){
	orm.RegisterModel(new(OrderItems))
}

func SaveItem( orderId int64, skuId int) error{
	var (
		err error
		item OrderItems
	)
	o := orm.NewOrm()
	item.OrderId = orderId
	item.ProductId = 1
	item.ProductSkuId = skuId
	item.Amount = 1
	item.Price = 1.0
	_,err = o.Insert(&item)

	return err
}

