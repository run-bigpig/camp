package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Nms     *CampConfig
	My      *CampConfig
	WebHook *WebHook
}

type CampConfig struct {
	Uid           string
	OpenId        string
	Token         string
	ListId        []string
	NtwNum        string
	CustomerPhone string
	CustomerName  string
}

type WebHook struct {
	Url    string
	Secret string
}

func RefreshToken(c *Config, NmsToken, MyToken string) {
	if NmsToken != "" {
		c.Nms.Token = NmsToken
	}
	if MyToken != "" {
		c.My.Token = MyToken
	}
}
