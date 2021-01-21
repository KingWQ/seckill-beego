package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["shop-seckill/controllers:OrderController"] = append(beego.GlobalControllerRouter["shop-seckill/controllers:OrderController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/cache/get`,
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["shop-seckill/controllers:OrderController"] = append(beego.GlobalControllerRouter["shop-seckill/controllers:OrderController"],
        beego.ControllerComments{
            Method: "Set",
            Router: `/cache/set`,
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["shop-seckill/controllers:OrderController"] = append(beego.GlobalControllerRouter["shop-seckill/controllers:OrderController"],
        beego.ControllerComments{
            Method: "Seckill",
            Router: `/order/seckill`,
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
