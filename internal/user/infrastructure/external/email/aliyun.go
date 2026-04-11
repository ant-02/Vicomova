package email

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"vicomova/pkg/config"
)

type EmailService interface {
	SendVerificationCode(ctx context.Context, toEmail, code string) error
}

type AliyunEmailService struct {
	accessKey    string
	accessSecret string
	accountName  string
	region       string
}

func NewAliyunEmailService(cfg *config.AliyunEmailConfig) (*AliyunEmailService, error) {
	return &AliyunEmailService{
		accessKey:    cfg.AccessKey,
		accessSecret: cfg.AccessSecret,
		accountName:  cfg.AccountName,
		region:       cfg.Region,
	}, nil
}

type aliyunResponse struct {
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
}

func (s *AliyunEmailService) SendVerificationCode(ctx context.Context, toEmail, code string) error {
	apiURL := "https://dm.aliyuncs.com/"

	params := url.Values{
		"AccessKeyId":      {s.accessKey},
		"AccountName":      {s.accountName},
		"Action":           {"SingleSendMail"},
		"AddressType":      {"1"},
		"Format":           {"JSON"},
		"HtmlBody":         {fmt.Sprintf("您的验证码是: <b>%s</b>，10分钟内有效。", code)},
		"RegionId":         {s.region},
		"SignatureMethod":  {"HMAC-SHA1"},
		"SignatureNonce":   {fmt.Sprintf("%d", time.Now().UnixNano())},
		"SignatureVersion": {"1.0"},
		"Subject":          {"Vicomova 邮箱验证码"},
		"Timestamp":        {time.Now().UTC().Format("2006-01-02T15:04:05Z")},
		"ToAddress":        {toEmail},
		"Version":          {"2015-11-23"},
	}

	sig := s.signature(params)
	params.Set("Signature", sig)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		var ar aliyunResponse
		if json.Unmarshal(body, &ar) == nil && ar.Code != "" {
			return fmt.Errorf("aliyun error: %s - %s", ar.Code, ar.Message)
		}
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// signature 计算 HMAC-SHA1 签名
// 阿里云签名规则：stringToSign = "POST&%2F&" + percentEncode(canonicalizedQueryString)
// canonicalizedQueryString = sortedPercentEncodedKeyValuePairs (key 和 value 分别编码，用 = 连接，pairs 之间用 & 连接)
func (s *AliyunEmailService) signature(params url.Values) string {
	// 1. 按字母顺序排序 key（排除 Signature）
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "Signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 2. 构建 canonicalized query string
	// url.Values 已对 key 和 value 做了 percent-encode，直接用
	var canonical strings.Builder
	for i, k := range keys {
		if i > 0 {
			canonical.WriteString("&")
		}
		canonical.WriteString(k)
		canonical.WriteString("=")
		canonical.WriteString(params.Get(k))
	}

	// 3. stringToSign = "POST&%2F&" + percentEncode(canonicalString)
	stringToSign := "POST&%2F&" + percentEncode(canonical.String())

	// 4. HMAC-SHA1
	mac := hmac.New(sha1.New, []byte(s.accessSecret+"&"))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// percentEncode RFC 3986 编码
func percentEncode(s string) string {
	encoded := url.QueryEscape(s)
	// url.QueryEscape 将空格编码为 +，但阿里云要求 + 编码为 %2B
	encoded = strings.ReplaceAll(encoded, "+", "%2B")
	return encoded
}