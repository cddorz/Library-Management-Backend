package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/smartwalle/alipay/v3"
	"log"
	"net/http"
)

func SignCheck(context *gin.Context) {
	Cfg, err := ini.Load("app.ini")
	if err != nil {
		log.Fatal("Fail to Load config: ", err)
	}

	ali, _ := Cfg.GetSection("alipay")
	appid, _ := ali.GetKey("appid")
	private, _ := ali.GetKey("private")
	alipublic, _ := ali.GetKey("alipublic")

	client, _ := alipay.New(appid.String(), private.String(), false)
	client.LoadAppPublicCertFromFile("appCertPublicKey.crt")
	client.LoadAliPayRootCertFromFile("alipayRootCert.crt")        // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("alipayCertPublicKey.crt") // 加载支付宝公钥证书
	client.LoadAliPayPublicKey(alipublic.String())
	req := context.Request
	req.ParseForm()
	ok, err := client.VerifySign(req.Form)
	fmt.Println(ok, err)
	if !ok {
		context.JSON(http.StatusOK, gin.H{"SignCheck": false})
	} else {
		context.JSON(http.StatusOK, gin.H{"SignCheck": true})
	}

}

// AliPayHandlerMobile 手机网页支付
// 传参示例http://127.0.0.1/pay?subject=fine&outtradeno=12340&totalamount=10
func AliPayHandlerMobile(context *gin.Context) {
	Cfg, err := ini.Load("app.ini")
	if err != nil {
		log.Fatal("Fail to Load config: ", err)
	}

	ali, _ := Cfg.GetSection("alipay")
	appid, _ := ali.GetKey("appid")
	private, _ := ali.GetKey("private")

	client, _ := alipay.New(appid.String(), private.String(), false)
	client.LoadAppPublicCertFromFile("appCertPublicKey.crt")
	client.LoadAliPayRootCertFromFile("alipayRootCert.crt")        // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("alipayCertPublicKey.crt") // 加载支付宝公钥证书
	var p = alipay.TradeWapPay{}
	p.NotifyURL = "_"
	p.ReturnURL = "http://127.0.0.1/pay/signcheck"
	p.QuitURL = "http://127.0.0.1/pay/signcheck"

	p.Subject = context.Query("subject")         // 订单标题
	p.OutTradeNo = context.Query("outtradeno")   // 商户订单号，64个字符以内、可包含字母、数字、下划线；需保证在商户端不重复
	p.TotalAmount = context.Query("totalamount") // 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	p.ProductCode = "QUICK_WAP_WAY"              // 销售产品码，与支付宝签约的产品码名称。 参考官方文档, App 支付时默认值为 QUICK_MSECURITY_PAY
	/*
		1、app支付product_code：QUICK_MSECURITY_PAY；
		2、手机网站支付product_code：QUICK_WAP_WAY；
		3、电脑网站支付product_code：FAST_INSTANT_TRADE_PAY；
		4、统一收单交易支付接口product_code：FACE_TO_FACE_PAYMENT；
		5、周期扣款签约product_code：CYCLE_PAY_AUTH；
	*/

	url, err := client.TradeWapPay(p)
	if err != nil {
		fmt.Println("pay client.TradeWapPay error:", err)
		return
	}

	binary, _ := url.MarshalBinary()
	fmt.Println(string(binary))
	data := make(map[string]interface{})
	data["url"] = string(binary)
	context.JSON(http.StatusOK, data)

}

// AliPayHandlerPC 电脑端网页支付
// 传参示例http://127.0.0.1/pay?subject=fine&outtradeno=12340&totalamount=10
func AliPayHandlerPC(context *gin.Context) {
	Cfg, err := ini.Load("app.ini")
	if err != nil {
		log.Fatal("Fail to Load config: ", err)
	}

	ali, _ := Cfg.GetSection("alipay")
	appid, _ := ali.GetKey("appid")
	private, _ := ali.GetKey("private")

	client, _ := alipay.New(appid.String(), private.String(), false)
	client.LoadAppPublicCertFromFile("appCertPublicKey.crt")
	client.LoadAliPayRootCertFromFile("alipayRootCert.crt")        // 加载支付宝根证书
	client.LoadAliPayPublicCertFromFile("alipayCertPublicKey.crt") // 加载支付宝公钥证书
	var p = alipay.TradePagePay{}
	p.NotifyURL = "_"
	p.ReturnURL = "http://127.0.0.1/pay/signcheck"

	p.Subject = context.Query("subject")         // 订单标题
	p.OutTradeNo = context.Query("outtradeno")   // 商户订单号，64个字符以内、可包含字母、数字、下划线；需保证在商户端不重复
	p.TotalAmount = context.Query("totalamount") // 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"     // 销售产品码，与支付宝签约的产品码名称。 参考官方文档, App 支付时默认值为 QUICK_MSECURITY_PAY
	/*
		1、app支付product_code：QUICK_MSECURITY_PAY；
		2、手机网站支付product_code：QUICK_WAP_WAY；
		3、电脑网站支付product_code：FAST_INSTANT_TRADE_PAY；
		4、统一收单交易支付接口product_code：FACE_TO_FACE_PAYMENT；
		5、周期扣款签约product_code：CYCLE_PAY_AUTH；
	*/

	url, err := client.TradePagePay(p)
	if err != nil {
		fmt.Println("pay client.TradePagePay error:", err)
		return
	}

	binary, _ := url.MarshalBinary()
	fmt.Println(string(binary))
	data := make(map[string]interface{})
	data["url"] = string(binary)
	context.JSON(http.StatusOK, data)
}
