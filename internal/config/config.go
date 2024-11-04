package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Nms *CampConfig
	My  *CampConfig
}

type CampConfig struct {
	Uid           string
	OpenId        string
	Token         string
	ListId        []string
	NtwNum        string
	CustomerPhone string
	CustomerName  string
	Job           *Job
}

type Job struct {
	Run      bool
	RoomId   []string
	RoomDate []string
}
