package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/micro/go-micro/client"

	"github.com/gin-gonic/gin"

	message "github.com/liuyuexclusive/future.srv.basic/proto/message"
	role "github.com/liuyuexclusive/future.srv.basic/proto/role"
	user "github.com/liuyuexclusive/future.srv.basic/proto/user"

	"github.com/liuyuexclusive/utils/webutil"
)

type LoginInput struct {
	UserName string `json:"username"` //用户名
	Password string `json:"password"` //密码
}

type LoginOutput struct {
	Token string `json:"token"` //令牌
}

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		res, err := user.NewUserService("go.micro.srv.basic", client.DefaultClient).Validate(context.TODO(), &user.ValidateRequest{Token: token})
		if err != nil {
			c.JSON(401, err.Error())
			c.Abort()
		}
		c.Set("username", res.Name)
	}
}

// Auth godoc
// @Summary
// @Description
// @Tags 获取token
// @Accept  json
// @Produce  json
// @Param account body handler.LoginInput true "input"
// @Success 200 {object} handler.LoginOutput "output"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /basic/login [post]
func Login(c *gin.Context) {
	var model LoginInput
	if ok := webutil.ReadBody(c, &model); ok {
		res, err := user.NewUserService("go.micro.srv.basic", client.DefaultClient).Auth(context.TODO(), &user.AuthRequest{Id: model.UserName, Key: model.Password})
		if err != nil {
			webutil.Bad(c, err)
		} else {
			webutil.OK(c, LoginOutput{Token: res.Token})
		}
	}
}

func Logout(c *gin.Context) {
	webutil.OK(c, "")
}

type RoleAddOrUpdateInput struct {
	ID   int64  `json:"id"`   //角色ID
	Name string `json:"name"` //角色名称
}

// RoleAddOrUpdate godoc
// @Summary
// @Description
// @Tags 新增or更新角色
// @Accept  json
// @Produce  json
// @Param Authorization header string true "授权码"
// @Param account body handler.RoleAddOrUpdateInput true "输入参数"
// @Success 200 {string} string "answer"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /basic/role/add [post]
func RoleAddOrUpdate(c *gin.Context) {

	var model RoleAddOrUpdateInput
	if ok := webutil.ReadBody(c, &model); ok {
		_, err := role.NewRoleService("go.micro.srv.basic", client.DefaultClient).AddOrUpdate(context.TODO(), &role.RoleAddOrUpdateRequest{Id: model.ID, Name: model.Name})
		if err != nil {
			webutil.Bad(c, err)
		} else {
			webutil.OK(c, "操作成功")
		}
	}
}

type CurrentUserOutput struct {
	Access []string `json:"access"`
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
}

func CurrentUser(c *gin.Context) {
	res, err := user.NewUserService("go.micro.srv.basic", client.DefaultClient).Get(context.TODO(), &user.GetRequest{Name: c.GetString("username")})
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	webutil.OK(c, CurrentUserOutput{Avatar: res.Avatar, Name: res.Name, Access: []string{"super_admin"}})
}

// func Test(c *gin.Context) {
// 	micro.NewPublisher("go.micro.srv.basic1", client.DefaultClient).Publish(context.TODO(), &basic.TestMessage{Name: "jiaojiao"})
// }

type Message struct {
	Title      string `json:"title"`
	CreateTime string `json:"create_time"`
	MsgID      int64  `json:"msg_id"`
}

type MessageOutput struct {
	Unread []Message `json:"unread"`
	Readed []Message `json:"readed"`
	Trash  []Message `json:"trash"`
}

func MessageInit(c *gin.Context) {
	res, err := message.NewMessageService("go.micro.srv.basic", client.DefaultClient).Init(context.TODO(), &message.InitRequest{To: c.GetString("username")})
	if err == nil {
		webutil.Bad(c, err)
	}
	var funcMap = func(s []*message.InitResponse_Message) []Message {
		news := make([]Message, 0)
		for _, v := range s {
			news = append(news, Message{Title: v.Title, MsgID: v.Id})
		}
		return news
	}
	webutil.OK(c, MessageOutput{Unread: funcMap(res.Unread), Readed: funcMap(res.Readed), Trash: funcMap(res.Trash)})
}

type MessageContentResult struct {
}

func MessageContent(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("msg_id"))
	if err != nil {
		webutil.Bad(c, err)
	}
	res, err := message.NewMessageService("go.micro.srv.basic", client.DefaultClient).Get(context.TODO(), &message.GetRequest{Id: int64(id)})
	if err != nil {
		webutil.Bad(c, err)
	}
	c.JSON(200, fmt.Sprintf(`%s`, res.Content))
}

type HasReadInput struct {
	MsgID int64 `json:"msg_id"`
}

func HasRead(c *gin.Context) {
	var input HasReadInput
	if webutil.ReadBody(c, &input) {
		_, err := message.NewMessageService("go.micro.srv.basic", client.DefaultClient).ChangeStatus(context.TODO(), &message.ChangeStatusRequest{Id: input.MsgID, Status: message.ChangeStatusRequest_Readed})
		if err != nil {
			webutil.Bad(c, err)
		}
		c.JSON(200, "")
	}
}

func RemoveReaded(c *gin.Context) {
	var input HasReadInput
	if webutil.ReadBody(c, &input) {
		_, err := message.NewMessageService("go.micro.srv.basic", client.DefaultClient).ChangeStatus(context.TODO(), &message.ChangeStatusRequest{Id: input.MsgID, Status: message.ChangeStatusRequest_Trash})
		if err != nil {
			webutil.Bad(c, err)
		}
		c.JSON(200, "")
	}
}

func Restore(c *gin.Context) {
	var input HasReadInput
	if webutil.ReadBody(c, &input) {
		_, err := message.NewMessageService("go.micro.srv.basic", client.DefaultClient).ChangeStatus(context.TODO(), &message.ChangeStatusRequest{Id: input.MsgID, Status: message.ChangeStatusRequest_Readed})
		if err != nil {
			webutil.Bad(c, err)
		}
		c.JSON(200, "")
	}
}

func MessageCount(c *gin.Context) {
	res, err := message.NewMessageService("go.micro.srv.basic", client.DefaultClient).Init(context.TODO(), &message.InitRequest{To: c.GetString("username")})
	if err != nil {
		webutil.Bad(c, err)
	} else {
		webutil.OK(c, len(res.Unread))
	}
}

func AddErrorLog(c *gin.Context) {
	var data map[string]interface{}
	if ok := webutil.ReadBody(c, &data); ok {
		fmt.Println(data)
		webutil.OK(c, ok)
	}
}
