package jsonmodel

type Award struct {
	DataId    int32 `json:"data_id"`
	DataValue int64 `json:"data_value"`
}

type Attach struct {
	Award []Award     `json:"award"`
	Data  interface{} `json:"data"`
}

type MailMsg struct {
	Userid      []int64 `json:"userid"`
	Title       string  `json:"title"`
	Body        string  `json:"body"`
	Attach      Attach  `json:"attach"`
	Type        int     `json:"type"`
	Sender      string  `json:"sender"`
	RXKeyNumber string  `json:"key_number"`
}

type RespMail struct {
	Code int32       ` json:"code,omitempty"`
	Msg  string      ` json:"msg,omitempty"`
	Data interface{} ` json:"data,omitempty" `
}

type MailAward struct {
	Userid int64 `json:"userid"`
	Mailid int64 `json:"mailid"`
}

type RspMailAward struct {
	Code    int32           `json:"code"`
	Msg     string          `json:"msg,omitempty"`
	MailIds []int64         `json:"mailids,omitempty"`
	Data    map[int32]int64 `json:"data,omitempty"`
}

type MailStatus struct {
	Userid   int64   `json:"userid"`
	MailsIds []int64 `json:"mailid"`
}

type RspMailStatus struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

type MailBatch struct {
	Userid   int64   `json:"userid"`
	MailsIds []int64 `json:"mailids"`
}

type RspMailBatchRead struct {
	Code    int32   `json:"code"`
	Msg     string  `json:"msg,omitempty"`
	MailIds []int64 `json:"mailids,omitempty"`
}

type RspMailBatchAward struct {
	Code    int32           `json:"code"`
	Msg     string          `json:"msg,omitempty"`
	MailIds []int64         `json:"mailids,omitempty"`
	Data    map[int32]int64 `json:"data,omitempty"`
	Content []string        `json:"content,omitempty"`
}

type RspMailBatchDel struct {
	Code    int32   `json:"code"`
	Msg     string  `json:"msg,omitempty"`
	MailIds []int64 `json:"mailids,omitempty"`
}

type RspMailAllDetail struct {
	Code    int32                    `json:"code"`
	Msg     string                   `json:"msg,omitempty"`
	MailIds []int64                  `json:"mailids,omitempty"` // 邮件ID
	Data    map[int64]MailDetailData `json:"data,omitempty"`    // 发货道具
}

type MailDetailData struct {
	MailID  int64           `json:"mail_id"`
	Rewards map[int32]int64 `json:"rewards"`
}

type MailEvaluation struct {
	Userid  int64 `json:"userid"`
	MailsId int64 `json:"mailid"`
	Pleased int   `json:"pleased"`
}
