//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

// IPC event to specified receiving destination

package target

// Type
//  0: Trigger the JS event of the specified target process
//  1: Trigger TgGoSub events for the specified target sub process
//  2: Trigger TgGoMain events for the specified target main process
type Type int8

const (
	TgJs     Type = iota //JS Event
	TgGoSub              //GO Event sub
	TgGoMain             //GO Event main
)

// ITarget
//
// ipc.NewTarget() *Target
type ITarget interface {
	BrowserId() int32 // Browser Window ID
	ChannelId() int64 // IPC channelID, frameId or GO IPC channelID
	TargetType() Type // Target type default 0: Trigger JS event
}

// Target Go IPC
//  receiving target of the event
type Target struct {
	browseId   int32
	channelId  int64
	targetType Type
}

// NewTarget Create a new Emit target
//	browserId: browser window ID
//	channelId: IPC channelID, frameId or GO IPC channelID
//	targetType: Optional parameter, target type default 0
//	  Type: TgJs:JS Event, TgGoSub:GO Sub Event, TgGoMain:GO Main Event
func NewTarget(browserId int32, channelId int64, targetType ...Type) ITarget {
	m := &Target{
		browseId:  browserId,
		channelId: channelId,
	}
	if len(targetType) > 0 {
		m.targetType = targetType[0]
	}
	return m
}

// NewTargetMain Create a new Emit target Main Process
//  targetType: TgGoMain
func NewTargetMain() ITarget {
	return &Target{
		targetType: TgGoMain,
	}
}

// TargetType
//  target type
//  0: Trigger JS event
//  1: Trigger Go Event
func (m *Target) TargetType() Type {
	return m.targetType
}

// BrowserId
//  return BrowserId
func (m *Target) BrowserId() int32 {
	return m.browseId
}

// ChannelId
//	return ChannelId
func (m *Target) ChannelId() int64 {
	return m.channelId
}
