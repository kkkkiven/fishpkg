package game

var	assign Assigner

// Assigner 桌号分配接口定义
type Assigner interface {

	// Assign 分配桌号和座次
	Assign(ext ... interface{})(deskID int32, seatID int32, err error)

	// Recycle 回收座次
	Recycle(deskID int32, seatID int32)
}

func SetAssign(a Assigner)  {
	assign = a
}
