package models

import (
	"github.com/astaxie/beego/orm"
)

type ProductSku struct{
	Id int
	Title string
	Stock int
}

type SkuInfo struct{
	Id      int    `json:"id"`
	Stock    int `json:"stock"`
}

func init(){
	orm.RegisterModel(new(ProductSku))
}

func GetSkuInfo(skuId int) (SkuInfo, error){
	o := orm.NewOrm()
	var sku SkuInfo
	err := o.Raw("select id,stock from product_skus where id=? limit 1", skuId).QueryRow(&sku)
	return sku,err
}

func UpdateStock(skuId int) (int64, error){
	var (
		err error
	)
	o := orm.NewOrm()
	res,err := o.Raw("update product_skus set stock = stock - 1 where stock >= 1 and id = ?", skuId).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		return num,err
	}else{
		return 0,err
	}
}