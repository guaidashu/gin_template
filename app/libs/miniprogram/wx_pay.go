/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 24/09/2021
 * @Desc: 微信支付 (目前只支持小程序)
 */

package miniprogram

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gin_template/app/libs/random"
	"sort"
	"strings"
	"time"
)

type (
	WxPay interface {
		// CreateOrder 创建订单
		CreateOrder(params CreateOrderParams) (err error)
		// GenSign 生成签名
		GenSign(key string, m map[string]string) (string, error)
		// Pay 小程序微信支付
		Pay(charge *Charge) (map[string]string, error)
	}

	defaultWxPay struct {
		createOrderUrl string
		appId          string       // 开发者appId
		mchId          string       // 商户号ID
		key            string       // 密钥
		privateKey     []byte       // 私钥文件内容
		publicKey      []byte       // 公钥文件内容
		httpsClient    *HTTPSClient // 双向证书链接
	}

	CreateOrderParams struct {
	}

	// Charge 支付参数
	Charge struct {
		AppId       string  `json:"-"`
		TradeNum    string  `json:"tradeNum,omitempty"`
		Origin      string  `json:"origin,omitempty"`
		UserId      string  `json:"userId,omitempty"`
		PayMethod   int64   `json:"payMethod,omitempty"`
		MoneyFee    float64 `json:"MoneyFee,omitempty"`
		CallbackURL string  `json:"callbackURL,omitempty"`
		ReturnURL   string  `json:"returnURL,omitempty"`
		ShowURL     string  `json:"showURL,omitempty"`
		Describe    string  `json:"describe,omitempty"`
		OpenId      string  `json:"openid,omitempty"`
		CheckName   bool    `json:"check_name,omitempty"`
		ReUserName  string  `json:"re_user_name,omitempty"`
		// 阿里提现
		AliAccount     string `json:"ali_account"`
		AliAccountType string `json:"ali_account_type"`
	}

	// WechatBaseResult 基本信息
	WechatBaseResult struct {
		ReturnCode string `xml:"return_code" json:"return_code,omitempty"`
		ReturnMsg  string `xml:"return_msg" json:"return_msg,omitempty"`
	}

	// WeChatResult 微信支付返回
	WeChatReResult struct {
		PrepayID string `xml:"prepay_id" json:"prepay_id,omitempty"`
		CodeURL  string `xml:"code_url" json:"code_url,omitempty"`
	}

	// WechatReturnData 返回通用数据
	WechatReturnData struct {
		AppID      string `xml:"appid,omitempty" json:"appid,omitempty"`
		MchID      string `xml:"mch_id,omitempty" json:"mch_id,omitempty"`
		MchAppid   string `xml:"mch_appid,omitempty" json:"mch_appid,omitempty"`
		DeviceInfo string `xml:"device_info,omitempty" json:"device_info,omitempty"`
		NonceStr   string `xml:"nonce_str,omitempty" json:"nonce_str,omitempty"`
		Sign       string `xml:"sign,omitempty" json:"sign,omitempty"`
		ResultCode string `xml:"result_code,omitempty" json:"result_code,omitempty"`
		ErrCode    string `xml:"err_code,omitempty" json:"err_code,omitempty"`
		ErrCodeDes string `xml:"err_code_des,omitempty" json:"err_code_des,omitempty"`
	}

	// WechatResultData 结果通用数据
	WechatResultData struct {
		OpenID         string `xml:"openid,omitempty" json:"openid,omitempty"`
		IsSubscribe    string `xml:"is_subscribe,omitempty" json:"is_subscribe,omitempty"`
		TradeType      string `xml:"trade_type,omitempty" json:"trade_type,omitempty"`
		BankType       string `xml:"bank_type,omitempty" json:"bank_type,omitempty"`
		FeeType        string `xml:"fee_type,omitempty" json:"fee_type,omitempty"`
		TotalFee       string `xml:"total_fee,omitempty" json:"total_fee,omitempty"`
		CashFeeType    string `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty"`
		CashFee        string `xml:"cash_fee,omitempty" json:"cash_fee,omitempty"`
		TransactionID  string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty"`
		OutTradeNO     string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty"`
		Attach         string `xml:"attach,omitempty" json:"attach,omitempty"`
		TimeEnd        string `xml:"time_end,omitempty" json:"time_end,omitempty"`
		PartnerTradeNo string `xml:"partner_trade_no,omitempty" json:"partner_trade_no,omitempty"`
		PaymentNo      string `xml:"payment_no,omitempty" json:"payment_no,omitempty"`
		PaymentTime    string `xml:"payment_time,omitempty" json:"payment_time,omitempty"`
		DetailId       string `xml:"detail_id,omitempty" json:"detail_id,omitempty"`
	}

	WeChatQueryResult struct {
		WechatBaseResult
		WeChatReResult
		WechatReturnData
		WechatResultData
		TradeState     string `xml:"trade_state" json:"trade_state,omitempty"`
		TradeStateDesc string `xml:"trade_state_desc" json:"trade_state_desc,omitempty"`
	}
)

func NewWxPay() WxPay {
	return &defaultWxPay{
		createOrderUrl: "https://api.mch.weixin.qq.com/pay/unifiedorder",
	}
}

func (s *defaultWxPay) CreateOrder(params CreateOrderParams) (err error) {
	return
}

func (s *defaultWxPay) GenSign(key string, m map[string]string) (string, error) {
	var signData []string
	for k, v := range m {
		if v != "" && k != "sign" && k != "key" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key

	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", errors.New("WechatGenSign md5.Write: " + err.Error())
	}
	signByte := c.Sum(nil)

	return strings.ToUpper(fmt.Sprintf("%x", signByte)), nil
}

// Pay 支付
func (s *defaultWxPay) Pay(charge *Charge) (map[string]string, error) {
	var m = make(map[string]string)
	appId := s.appId
	if charge.AppId != "" {
		appId = charge.AppId
	}
	m["appid"] = appId
	m["mch_id"] = s.mchId
	m["nonce_str"] = random.GenRandomTimestampStr() + charge.UserId
	m["body"] = TruncatedText(charge.Describe, 32)
	m["out_trade_no"] = charge.TradeNum
	m["total_fee"] = WechatMoneyFeeToString(charge.MoneyFee)
	m["spbill_create_ip"] = LocalIP()
	m["notify_url"] = charge.CallbackURL
	m["trade_type"] = "JSAPI"
	m["openid"] = charge.OpenId
	m["sign_type"] = "MD5"

	sign, err := s.GenSign(s.key, m)
	if err != nil {
		return map[string]string{}, err
	}
	m["sign"] = sign

	// 转出xml结构
	xmlRe, err := s.PostWechat("https://api.mch.weixin.qq.com/pay/unifiedorder", m, nil)
	if err != nil {
		return map[string]string{}, err
	}

	var c = make(map[string]string)
	c["appId"] = appId
	c["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	c["nonceStr"] = random.GenRandomTimestampStr() + charge.UserId
	c["package"] = fmt.Sprintf("prepay_id=%s", xmlRe.PrepayID)
	c["signType"] = "MD5"
	sign2, err := s.GenSign(s.key, c)
	if err != nil {
		return map[string]string{}, errors.New("WechatWeb: " + err.Error())
	}
	c["paySign"] = sign2
	delete(c, "appId")
	return c, nil
}

// 对微信下订单或者查订单
func (s *defaultWxPay) PostWechat(url string, data map[string]string, h *HTTPSClient) (WeChatQueryResult, error) {
	var xmlRe WeChatQueryResult
	buf := bytes.NewBufferString("")

	for k, v := range data {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())

	hc := new(HTTPSClient)
	if h != nil {
		hc = h
	} else {
		hc = HTTPSC
	}

	re, err := hc.PostData(url, "text/xml:charset=UTF-8", xmlStr)
	if err != nil {
		return xmlRe, errors.New("HTTPSC.PostData: " + err.Error())
	}

	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		return xmlRe, errors.New("xml.Unmarshal: " + err.Error())
	}

	if xmlRe.ReturnCode != "SUCCESS" {
		// 通信失败
		return xmlRe, errors.New("xmlRe.ReturnMsg: " + xmlRe.ReturnMsg)
	}

	if xmlRe.ResultCode != "SUCCESS" {
		// 业务结果失败
		return xmlRe, errors.New("xmlRe.ErrCodeDes: " + xmlRe.ErrCodeDes)
	}
	return xmlRe, nil
}

// 关闭订单
func (s *defaultWxPay) CloseOrder(outTradeNo string) (WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = s.appId
	m["mch_id"] = s.mchId
	m["nonce_str"] = random.GenRandomTimestampStr()
	m["out_trade_no"] = outTradeNo
	m["sign_type"] = "MD5"

	sign, err := s.GenSign(s.key, m)
	if err != nil {
		return WeChatQueryResult{}, err
	}
	m["sign"] = sign

	// 转出xml结构
	result, err := s.PostWechat("https://api.mch.weixin.qq.com/pay/closeorder", m, nil)
	if err != nil {
		return WeChatQueryResult{}, err
	}

	return result, err
}

// 支付到用户的微信账号
func (s *defaultWxPay) PayToClient(charge *Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["mch_appid"] = charge.AppId
	m["mchid"] = s.mchId
	m["nonce_str"] = random.GenRandomTimestampStr()
	m["partner_trade_no"] = charge.TradeNum
	m["openid"] = charge.OpenId
	m["amount"] = WechatMoneyFeeToString(charge.MoneyFee)
	m["spbill_create_ip"] = LocalIP()
	m["desc"] = TruncatedText(charge.Describe, 32)

	// 是否验证用户名称
	if charge.CheckName {
		m["check_name"] = "FORCE_CHECK"
		m["re_user_name"] = charge.ReUserName
	} else {
		m["check_name"] = "NO_CHECK"
	}

	sign, err := s.GenSign(s.key, m)
	if err != nil {
		return map[string]string{}, err
	}
	m["sign"] = sign

	// 转出xml结构
	result, err := s.PostWechat("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", m, s.httpsClient)
	if err != nil {
		return map[string]string{}, err
	}

	return s.struct2Map(result)
}

// QueryOrder 查询订单
func (s *defaultWxPay) QueryOrder(tradeNum string) (WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = s.appId
	m["mch_id"] = s.mchId
	m["out_trade_no"] = tradeNum
	m["nonce_str"] = random.GenRandomTimestampStr()

	sign, err := s.GenSign(s.key, m)
	if err != nil {
		return WeChatQueryResult{}, err
	}
	m["sign"] = sign

	return s.PostWechat("https://api.mch.weixin.qq.com/pay/orderquery", m, nil)
}

func (s *defaultWxPay) struct2Map(obj interface{}) (map[string]string, error) {

	j2 := make(map[string]string)

	j1, err := json.Marshal(obj)
	if err != nil {
		return j2, err
	}

	err2 := json.Unmarshal(j1, &j2)
	return j2, err2
}
