package utils

import (
	"IM_BE/dto"
	"IM_BE/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"time"
)

// MessageChanger 将 MessageReqDto 转化为 MessageRespDto
func MessageChanger(messageBytes []byte) (*dto.MessageRespDto, error) {
	messageReq := &dto.MessageReqDto{} // 注意，前端传过来的类型
	if err := json.Unmarshal(messageBytes, messageReq); err != nil {
		GetLogger().Error("解析消息体失败", zap.Error(err))
		return nil, err
	}
	GetLogger().Debug("获取消息成功", zap.String("messageBytes", string(messageBytes)))

	messageResp := &dto.MessageRespDto{}
	if err := copier.Copy(messageResp, messageReq); err != nil {
		GetLogger().Error("复制消息体失败", zap.Error(err))
		return nil, err
	}

	switch messageReq.ContentType {
	case model.AudioMessage, model.ImageMessage, model.FileMessage, model.TextMessage, model.SystemNotice:
		switch messageReq.ContentType {
		case model.TextMessage, model.SystemNotice:
			if err := json.Unmarshal(messageReq.Content, &messageResp.Content); err != nil {
				GetLogger().Error("解析 content 失败", zap.Error(err))
				return nil, err
			}
		case model.ImageMessage:
			fileBytes := messageReq.Content
			fileName := fmt.Sprintf("upload/image/%d.jpg", time.Now().UnixNano())

			fileUrl, err := UpdateFile(fileBytes, fileName, "image/jpeg")
			if err != nil {
				GetLogger().Error("文件上传失败", zap.Error(err))
				return nil, err
			}

			messageResp.Content = fileUrl
		case model.AudioMessage:
			//TODO
		}

	default:
		GetLogger().Error("消息格式错误")
		return nil, errors.New("消息格式错误")
	}

	GetLogger().Debug("消息转化成功")
	return messageResp, nil
}
