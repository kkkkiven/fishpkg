package jsonmodel

type RespPay struct {
	Code int32       ` json:"code,omitempty"`
	Msg  string      ` json:"msg,omitempty"`
	Data interface{} ` json:"data,omitempty" `
}

type CallBackPay struct {
	OrderId       string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
}
