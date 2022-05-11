package sms

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"gohub/pkg/logger"
)

// Tencent 实现 sms.Driver interface
type Tencent struct{}

// Send 实现 sms.Driver interface 的 Send 方法
func (s *Tencent) Send(phone string, message Message, config map[string]string) bool {
	logger.DebugJSON("短信[腾讯云]", "配置信息", config)

	credential := common.NewCredential(
		config["access_key_id"],
		config["access_key_secret"],
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := tencentsms.NewClient(credential, "ap-guangzhou", cpf)

	request := tencentsms.NewSendSmsRequest()
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.SignName = common.StringPtr(config["sign_name"])
	request.TemplateId = common.StringPtr(config["template_code"])
	// 为了和阿里云保持兼容，使用腾讯云发送短信时传入的 sms.Message{} 中的 Data 格式还是
	// Data: map[string]string{"code": code}
	request.TemplateParamSet = common.StringPtrs([]string{message.Data["code"]})
	request.SmsSdkAppId = common.StringPtr(config["sdk_app_id"])

	response, err := client.SendSms(request)
	logger.DebugJSON("短信[腾讯云]", "请求内容", request)
	logger.DebugJSON("短信[腾讯云]", "接口响应", response)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logger.ErrorString("短信[腾讯云]", "调用接口错误", err.Error())
		return false
	}
	if err != nil {
		logger.ErrorString("短信[腾讯云]", "服务商返回错误", err.Error())
		return false
	}

	// 当前只有一条
	statusSet := response.Response.SendStatusSet
	code := *statusSet[0].Code
	if code == "Ok" {
		logger.DebugString("短信[腾讯云]", "发信成功", response.ToJsonString())
		return true
	} else {
		logger.ErrorString("短信[腾讯云]", "发信失败", response.ToJsonString())
		return false
	}

}
