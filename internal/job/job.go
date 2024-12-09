package job

import (
	"camp/internal/config"
	"camp/internal/consts"
	"camp/internal/notice"
	"camp/internal/types"
	"camp/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand/v2"
	"strings"
	"sync"
	"time"
)

var platformMap = map[int8]string{
	0: "柠檬树",
	1: "牧云",
}

type JobList struct {
	list  map[string]*types.Job
	close map[string]chan bool
	mu    sync.Mutex
}

type Job struct {
	logger  logx.Logger
	jobList *JobList
	conf    *config.Config
}

func NewJob(c *config.Config) *Job {
	return &Job{
		conf:    c,
		logger:  logx.WithContext(context.TODO()),
		jobList: &JobList{list: make(map[string]*types.Job), close: make(map[string]chan bool), mu: sync.Mutex{}},
	}
}

func (j *Job) AddJob(job *types.OperateJobRequest) error {
	j.logger.Info("add job")
	j.jobList.mu.Lock()
	defer j.jobList.mu.Unlock()
	jobId := uuid.New().String()
	j.jobList.list[jobId] = &types.Job{
		JobId: jobId,
	}
	j.jobList.close[jobId] = make(chan bool)
	go j.commit(job.Platform, job.Interval, job.Room, j.jobList.close[jobId])
	return nil
}

func (j *Job) DeleteJob(jobId string) error {
	j.logger.Info("delete job")
	j.jobList.mu.Lock()
	defer j.jobList.mu.Unlock()
	if _, ok := j.jobList.list[jobId]; !ok {
		return fmt.Errorf("jobId %s不存在", jobId)
	}
	j.jobList.close[jobId] <- true
	delete(j.jobList.list, jobId)
	delete(j.jobList.close, jobId)
	return nil
}

func (j *Job) commit(platform int8, interval int32, room *types.CommitRoom, close chan bool) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	var campConf *config.CampConfig
	if platform == 0 {
		campConf = j.conf.Nms
	} else {
		campConf = j.conf.My
	}
	for {
		select {
		case t := <-ticker.C:
			j.logger.Infof("%s %s开始提交订单", platformMap[platform], t.Format("2006-01-02 15:04:05"))
			err := j.requestCommit(campConf, room)
			if err != nil {
				if strings.Contains(err.Error(), "重新登录") {
					_ = notice.SendMsg(j.conf.WebHook.Url, j.conf.WebHook.Secret, fmt.Sprintf("%s %s %s", platformMap[platform], t.Format("2006-01-02 15:04:05"), err.Error()))
					return
				}
				j.logger.Errorf("%s %s订单提交失败：%s", platformMap[platform], t.Format("2006-01-02 15:04:05"), err.Error())
				continue
			}
			err = notice.SendMsg(j.conf.WebHook.Url, j.conf.WebHook.Secret, fmt.Sprintf("%s %s提交订单成功：%s", platformMap[platform], t.Format("2006-01-02 15:04:05"), room.RoomName))
			if err != nil {
				j.logger.Errorf("%s %s提交订单成功，但通知失败：%s", platformMap[platform], t.Format("2006-01-02 15:04:05"), err.Error())
				return
			}
			j.logger.Infof("%s %s提交订单成功：%s", platformMap[platform], t.Format("2006-01-02 15:04:05"), room.RoomName)
		case <-close:
			j.logger.Infof("%s %s任务结束", platformMap[platform], time.Now().Format("2006-01-02 15:04:05"))
			return
		}
	}
}

func (j *Job) requestCommit(config *config.CampConfig, room *types.CommitRoom) error {
	r := &types.CommitRequest{
		AppType:                 0,
		BookingType:             0,
		BoolPublish:             true,
		BoolUseContinueLiveInst: false,
		BoolUseIntegral:         false,
		Count:                   1,
		CustomerName:            config.CustomerName,
		CustomerPhone:           config.CustomerPhone,
		DiscountInstList:        []interface{}{},
		NtwNum:                  config.NtwNum,
		OpenID:                  config.OpenId,
		PromotionChannel:        3,
		Remark:                  "",
		RoomTypeID:              room.RoomId,
		StartDate:               room.RoomStartDate,
		EndDate:                 room.RoomEndDate,
	}
	header := map[string]string{
		"Referer":    consts.Referer,
		"User-Agent": consts.UserAgent,
		"token":      config.Token,
		"uid":        config.Uid,
	}
	randNum := rand.IntN(9)
	time.Sleep(time.Duration(randNum) * time.Second)
	data, err := utils.SendRequest(consts.CommitUrl, header, r)
	if err != nil {
		return err
	}
	var commitResp types.CommitResponse
	err = json.Unmarshal(data, &commitResp)
	if err != nil {
		return err
	}
	if commitResp.Code != "1" {
		return errors.New(commitResp.Msg)
	}
	return nil
}
