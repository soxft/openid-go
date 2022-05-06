package mailutil

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"openid/config"
)

func SendByAliyun(mail Mail) error {
	AliyunConfig := config.C.Aliyun

	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", AliyunConfig.AccessKey, AliyunConfig.AccessSecret)
	if err != nil {
		return err
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dm.aliyuncs.com"
	request.Version = "2015-11-23"
	request.ApiName = "SingleSendMail"
	request.QueryParams["ToAddress"] = mail.ToAddress
	request.QueryParams["Subject"] = mail.Subject + " - " + config.C.Server.Title
	request.QueryParams["HtmlBody"] = mail.Content
	request.QueryParams["FromAlias"] = config.C.Server.Title
	request.QueryParams["AccountName"] = AliyunConfig.Email
	request.QueryParams["AddressType"] = "1"
	request.QueryParams["ReplyToAddress"] = "true"

	_, err = client.ProcessCommonRequest(request)
	return err
}