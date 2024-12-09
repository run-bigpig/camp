package notice

import (
	"camp/internal/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type FeiShuMsg struct {
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
	MsgType   string `json:"msg_type"`
	Content   struct {
		Text string `json:"text"`
	} `json:"content"`
}

type FeiShuResponse struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

func SendMsg(url, secret string, msg string) error {
	t := time.Now().Unix()
	sign, err := genSign(secret, t)
	if err != nil {
		return err
	}
	req := &FeiShuMsg{
		Timestamp: t,
		Sign:      sign,
		MsgType:   "text",
	}
	req.Content.Text = msg
	resp, err := utils.SendRequest(url, nil, req)
	if err != nil {
		return err
	}
	var res FeiShuResponse
	err = json.Unmarshal(resp, &res)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf("send feishu msg failed, code:%d, msg:%s", res.Code, res.Msg)
	}
	return nil
}

func genSign(secret string, timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
