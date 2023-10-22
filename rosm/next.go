package rosm

import ()

var (
	nextMessList = map[int64]chan *CTX{}
	//nextEmoticonList = map[string]chan *CTX{}
)

// 获取本房间全体的下一句话
func (ctx *CTX) GetNextAllMess() (chan *CTX, func()) {
	next := make(chan *CTX, 1)
	id := ctx.Being.RoomID + ctx.Being.RoomID2
	nextMessList[id] = next
	return next, func() {
		close(next)
		delete(nextMessList, id)
	}
}

// 获取本房间该用户的下一句话
func (ctx *CTX) GetNextUserMess() (chan *CTX, func()) {
	next := make(chan *CTX, 1)
	id := ctx.Being.RoomID + ctx.Being.RoomID2 + ctx.Being.User.ID
	nextMessList[id] = next
	return next, func() {
		close(next)
		delete(nextMessList, id)
	}
}

func (ctx *CTX) sendNext() (block bool) {
	if len(nextMessList) == 0 {
		return false
	}
	//先匹配个人
	if c, ok := nextMessList[ctx.Being.RoomID+ctx.Being.RoomID2+ctx.Being.User.ID]; ok {
		c <- ctx
		return true
	}
	if c, ok := nextMessList[ctx.Being.RoomID+ctx.Being.RoomID2]; ok {
		c <- ctx
		return true
	}
	return false
}
