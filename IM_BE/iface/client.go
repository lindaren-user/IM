package iface

import "IM_BE/dto"

// 解除循环导入
type Client interface {
	GetId() uint64
	GetMessage(dto dto.MessageDto)
}
