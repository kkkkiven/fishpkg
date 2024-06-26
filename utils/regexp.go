package utils

import (
	//	"log"
	"regexp"
)

const (
	regUsername           = `^[a-zA-Z0-9_]{4,22}$`
	regPwd                = `^[\@A-Za-z0-9\!\#\$\%\^\&\*\.\~]{3,60}$`
	regNickname           = `^{0,40}$`
	regEmail              = `^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$`
	regPhone              = `^((\d3)|(\d{3}\-))?13[0-9]\d{8}|14[0-9]\d{8}|15[0-9]\d{8}|17[0-9]\d{8}|18[0-9]\d{8}`
	regUrl                = `^((https?|ftp|news|http):\/\/)?([a-z]([a-z0-9\-]*[\.。])+([a-z]{2}|aero|arpa|biz|com|coop|edu|gov|info|int|jobs|mil|museum|name|nato|net|org|pro|travel)|(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))(\/[a-z0-9_\-\.~]+)*(\/([a-z0-9_\-\.]*)(\?[a-z0-9+_\-\.%=&]*)?)?(#[a-z][a-z0-9_]*)?$`
	regGuid               = `[a-zA-Z0-9-_]{1,40}`
	regDescription        = `^{0,64}$`
	regOutTypeDescription = `^{0,20}$`
	regMac                = `^{0,40}$`
	regTradeNo            = `^[a-zA-Z0-9_-]{1,40}$`
	regAttach             = `^{0,127}$`

	// LT
	regID                  = `^[0-9]{0,11}$`
	regTitle               = `^[\s\S]{0,40}$`
	regIntro               = `^[\s\S]{0,120}$`
	regHash                = `^[\S]{0,160}$`
	reqAuthorName          = `^{0,60}$`
	reqAtType              = `^[1-7]{1}$`
	reqActionListAtType    = `^[1,3]{1}$`
	regContent             = `^[\s\S]{0,255}$`
	reqSrcType             = `^[1,2,3]{1}$`
	reqCommentActionAtType = `^[1,2,6]{1}$`
	reqUserActionAtType    = `^[1,2,3,4]{1}$`
	reqKeyword             = `^.{0,40}$`
	regReason              = `^[\s\S]{0,140}$`
	regUserSuggestTitle    = `^[\s\S]{0,120}$`
)

func CheckString(data string, pat string) bool {
	bFlag := false
	reg := regexp.MustCompile(pat)
	bFlag = reg.MatchString(data)
	return bFlag
}

// 检测用户名
func CheckUserName(username string) bool {
	return CheckString(username, regUsername)
}

// 昵称
func CheckNickname(nickname string) bool {
	return CheckString(nickname, regNickname)

}

// 密码
func CheckPwd(password string) bool {
	return CheckString(password, regPwd)
}

// 邮箱
func CheckEmail(email string) bool {
	return CheckString(email, regEmail)
}

// 手机号
func CheckPhone(phone string) bool {
	return CheckString(phone, regPhone)
}

// 网址
func CheckUrl(url string) bool {
	return CheckString(url, regUrl)
}

// guid
func CheckGuid(guid string) bool {
	return CheckString(guid, regGuid)
}

// 描述
func CheckDescription(description string) bool {
	return CheckString(description, regDescription)
}

// 平台交易类型描述
func CheckOutTypeDescription(description string) bool {
	return CheckString(description, regOutTypeDescription)
}

// mac 地址
func CheckMac(mac string) bool {
	return CheckString(mac, regMac)
}

// 订单号
func CheckTradeNo(tradeNo string) bool {
	return CheckString(tradeNo, regTradeNo)
}

// 附加参数
func CheckAttach(attach string) bool {
	return CheckString(attach, regAttach)
}

// 检测资源id
func CheckID(anchorId string) bool {
	return CheckString(anchorId, regID)
}

// 检测资源title
func CheckTitle(title string) bool {
	return CheckString(title, regTitle)
}

// 检测资源intro
func CheckIntro(intro string) bool {
	return CheckString(intro, regIntro)
}

// 检测资源Hash
func CheckHash(hash string) bool {
	return CheckString(hash, regHash)
}

// 检测资源AuthorName
func CheckAuthorName(authorName string) bool {
	return CheckString(authorName, reqAuthorName)
}

// 检测资源atType
func CheckAtType(atType string) bool {
	return CheckString(atType, reqAtType)
}

// 检测资源atType
func CheckActionListAtType(atType string) bool {
	return CheckString(atType, reqActionListAtType)
}

// 检测资源atType
func CheckContent(content string) bool {
	return CheckString(content, regContent)
}

// 检测资源atType
func CheckSrcType(srcType string) bool {
	return CheckString(srcType, reqSrcType)
}

// 检测资源atType
func CheckCommentActionAtType(atType string) bool {
	return CheckString(atType, reqCommentActionAtType)
}

// 检测资源atType
func CheckUserActionAtType(atType string) bool {
	return CheckString(atType, reqUserActionAtType)
}

// 检测资源keyword
func CheckKeyword(keyword string) bool {
	return CheckString(keyword, reqKeyword)
}

// 检测资源Reason
func CheckReason(reason string) bool {
	return CheckString(reason, regReason)
}

// 检测资源SuggestTitle
func CheckSuggestTitle(title string) bool {
	return CheckString(title, regUserSuggestTitle)
}
