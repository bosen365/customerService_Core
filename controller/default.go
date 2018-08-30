// 健康检查

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type DefaultController struct {
}

func InitHealth() *DefaultController {
	return &DefaultController{}
}

// @Summary 健康检查
// @Description 应用程序健康检查接口
// @Tags Default
// @Accept json
// @Produce json
// @Success 200 {string} json ""
// @Router /v1/health [get]
func (c *DefaultController) Health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format("2006-01-02 15:04:05")})
}

// @Summary 系统初始化
// @Description 在线客服系统进行初始化
// @Tags Default
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/init [get]
func (c *DefaultController) Init(context *gin.Context) {
	// 获取访问客服信息
	roomKf, err := handle.AuthToken2Model(context)
	ReturnErrInfo(context, err)
	kfDb := &model.Kf{Id: roomKf.KfId}
	if err := kfDb.Get(); err != nil {
		ReturnErrInfo(context, err.Error())
	}

	// 获取聊天历史记录
	messageDb := model.MessageLinkCustomer{Message: model.Message{KfId: kfDb.Id}}
	messages, err := messageDb.GetKfHistoryMsg()
	ReturnErrInfo(context, err)

	var initOnlineCustomers []InitOnlineCustomer
	var mapCus = map[int]*InitOnlineCustomer{}

	// 组织数据
	for _, singeMsg := range messages {
		if _, ok := mapCus[singeMsg.KfId]; ok {
			mapCus[singeMsg.KfId].CustomerMessages = append(mapCus[singeMsg.KfId].CustomerMessages, InitMessage{
				Id:                singeMsg.Id,
				MessageType:       singeMsg.MsgType,
				MessageContent:    singeMsg.Msg,
				MessageOperCode:   singeMsg.OperCode,
				MessageAck:        singeMsg.KfAck,
				MessageCteateTime: singeMsg.CreateTime,
			})
		} else {
			mapCus[singeMsg.KfId] = &InitOnlineCustomer{
				RoomToken:          singeMsg.CustomerToken,
				CustomerNickName:   singeMsg.CustomerNickName,
				CustomerHeadImgUrl: singeMsg.CustomerHeadImgUrl,
				CustomerMessages: []InitMessage{
					{
						Id:                singeMsg.Id,
						MessageType:       singeMsg.MsgType,
						MessageContent:    singeMsg.Msg,
						MessageOperCode:   singeMsg.OperCode,
						MessageAck:        singeMsg.KfAck,
						MessageCteateTime: singeMsg.CreateTime,
					},
				},
			}
		}
	}
	for _, v := range mapCus {
		initOnlineCustomers = append(initOnlineCustomers, *v)
	}

	// 组织排队数据
	var waitQueues []WaitQueueResponse
	if waitQueueRooms, err := logic.GetWaitQueue(); err != nil {
		ReturnErrInfo(context, err)
	} else {
		for _, value := range waitQueueRooms {
			waitQueues = append(waitQueues, WaitQueueResponse{
				CustomerId:         value.CustomerId,
				CustomerNickName:   value.CustomerNickName,
				CustomerHeadImgUrl: value.CustomerHeadImgUrl,
				//Messages:           value.CustomerMsgs,
				PreviousKf: WaitQueuePreviousKf{},
			})
		}
	}

	context.JSON(http.StatusOK, InitResponse{
		Mine: InitMine{
			Id:         kfDb.Id,
			UserName:   kfDb.NickName,
			HeadImgUrl: kfDb.HeadImgUrl,
			Status:     string(common.KF_ONLINE),
		},
		InitOnlineCustomer: initOnlineCustomers,
		WaitQueueResponse:  waitQueues,
	})
}

// 异常返回
func ReturnErrInfo(context *gin.Context, err interface{}) {
	if err != nil {
		log.Printf("发生异常：%#v", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "接口调用异常，或联系管理员",
		})
		panic(err)
	}
}

// 成功返回
func ReturnSuccessInfo(context *gin.Context) {
	context.JSON(http.StatusInternalServerError, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

type InitResponse struct {
	Mine               InitMine             `json:"mine"`
	InitOnlineCustomer []InitOnlineCustomer `json:"init_online_customer"`
	WaitQueueResponse  []WaitQueueResponse  `json:"wait_queue_response"`
}
type InitMine struct {
	Id         int    `json:"id"`
	UserName   string `json:"user_name"`
	HeadImgUrl string `json:"head_img_url"`
	Status     string `json:"status"`
}
type InitOnlineCustomer struct {
	RoomToken          string        `json:"room_token"` // 会话的Token，实际上就是用户的OpenId
	CustomerNickName   string        `json:"customer_nick_name"`
	CustomerHeadImgUrl string        `json:"customer_head_img_url"`
	CustomerMessages   []InitMessage `json:"customer_messages"`
}
type InitMessage struct {
	Id                int       `json:"id"`
	MessageType       string    `json:"message_type"`
	MessageContent    string    `json:"message_content"`
	MessageOperCode   int       `json:"message_oper_code"`
	MessageAck        bool      `json:"message_ack"`
	MessageCteateTime time.Time `json:"message_cteate_time"`
}
