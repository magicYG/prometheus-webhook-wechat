package notifier

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"prometheus-webhook-wechat/models"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func SendNotification(notifyTargets []models.Target, notifyContent string, logger log.Logger, callID string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 取消证书验证
	}
	httpClient := &http.Client{Transport: tr}

	notification := &models.WechatNotification{
		MsgType: "markdown",
		Markdown: &models.WechatNotificationMarkdown{
			Content: notifyContent,
		},
	}
	reqBody, err := json.Marshal(notification)
	if err != nil {

		level.Error(logger).Log("TraceID", callID, "Msg", "Encoding request body error", "Error", err)
		return
	}

	for _, v := range notifyTargets {
		level.Debug(logger).Log("reqbody", string(reqBody))
		httpReq, err := http.NewRequest("POST", v.URL, bytes.NewReader(reqBody))
		if err != nil {
			level.Error(logger).Log("TraceID", callID, "Msg", "Building request body error", "URL", v.URL, "Error", err)
			continue
		}
		httpReq.Header.Set("Content-Type", "application/json")
		rsp, err := httpClient.Do(httpReq)
		if err != nil {
			level.Error(logger).Log("TraceID", callID, "Msg", "Sending request to error", "URL", v.URL, "Error", err)
			continue
		}
		if rsp.StatusCode != 200 {
			level.Error(logger).Log("TraceID", callID, "Msg", "Call wechat API failed", "URL", v.URL, "ResponseCode", rsp.StatusCode)
			continue
		}
		var robotRsp models.WechatNotificationResponse
		rspDec := json.NewDecoder(rsp.Body)
		if err := rspDec.Decode(&robotRsp); err != nil {
			level.Error(logger).Log("TraceID", callID, "Msg", "Error decoding response from wechat", "Error", err)
			continue
		}
		level.Info(logger).Log("TraceID", callID, "WechatResponseCode", robotRsp.ErrCode, "WechatResponseMsg", robotRsp.ErrMsg, "URL", v.URL)
	}
}
