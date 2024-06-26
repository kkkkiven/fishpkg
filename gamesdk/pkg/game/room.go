package game

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	sdk "git.yuetanggame.com/zfish/fishpkg/servicesdk/core"

	. "git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/errors"
	. "git.yuetanggame.com/zfish/fishpkg/gamesdk/pkg/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// reqGetRoomCardIDs 批量申请房卡号
type reqGetRoomCardIDs struct {
	RoomID   int64 `json:"roomid"`
	ServerID int64 `json:"serverid"`
	Num      int   `json:"num"`
}

// respGetRoomCardIDs 批量申请房卡号（响应）
type respGetRoomCardIDs struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	IDs  []string `json:"ids"`
}

// reqReleaseRoomCardIDs 批量释放房卡号
type reqReleaseRoomCardIDs struct {
	RoomID   int64    `json:"roomid"`
	ServerID int64    `json:"serverid"`
	IDs      []string `json:"ids"`
}

// respReleaseRoomCardIDs 批量释放房卡（响应）
type respReleaseRoomCardIDs struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type roomCardPool struct {
	sync.RWMutex

	serverID             uint32           // 当前服务IP
	roomID               int64            // 游戏房间ID
	currentSize          int32            // 当前房卡池大小
	initialSize          int32            // 初始房卡池大小
	scaleUpSizePerTime   int32            // 每次扩容申请的房卡数量
	scaleDownSizePerTime int32            // 每次缩容释放的房卡数量
	scaleUpThreshold     float64          // 触发扩容的阀值
	scaleDownThreshold   float64          // 触发缩容的阀值
	freeIdx              int              // 可分配房卡列表的起始数组下标
	shared               []string         // 房卡池
	used                 map[string]*card // 已使用房卡信息
	pendingIdx           int              // 待释放房卡列表的起始数组下标
	status               pstatus          // 房卡池状态: NOOP：无操作；NEW：申请中；FREE：释放中
	freeTimer            *time.Timer      // 释放计时器
}

// 房卡池状态，NOOP：无操作；NEW：向分配服申请房卡；FREE：向分配服释放房卡
type pstatus uint8

const (
	STAT_NOOP pstatus = 0 // 无操作
	STAT_NEW          = 1 // 申请中
	STAT_FREE         = 2 // 释放中
)

const (
	DEFAULT_INITIAL_SIZE             = 100 // 默认房卡池初值大小
	DEFAULT_SCALE_UP_SIZE_PER_TIME   = 50  // 默认每次房卡池扩容数量
	DEFAULT_SCALE_DOWN_SIZE_PER_TIME = 50  // 默认每次房卡池缩容数量
	DEFAULT_SCALE_UP_THRESHOLD       = 0.8 // 默认房卡池扩容阀值
	DEFAULT_SCALE_DOWN_THRESHOLD     = 0.4 // 默认房卡池缩容阀值
)

const RCPOOL_RELEASE_INTERVAL = 1 // 房卡池定时释放间隔，单位秒

// card 房卡
type card struct {
	idx int    // 房卡池底层数组下标
	Id  string // 房卡ID
}

var (
	once   sync.Once
	rcPool *roomCardPool
)

// StartRoomCardPool 开启游戏房卡池申请游戏房间号
// serverID: 当前服务ID，即整型的IP
// roomID: 当前服务房间ID
// initialSize: 初始房卡池大小
// scaleUpSizePerTime: 每次扩容申请的房间数量
// scaleDownSizePerTime: 每次缩容释放的房间数量
// scaleUpThreshold: 触发扩容的阀值，当使用率大于该值则触发扩容
// scaleDownThreshold: 触发缩容的阀值，当使用率低于该值则触发缩容。若扩容阀值不为0，则缩容阀值必须小于扩容阀值
func StartRoomCardPool(ctx context.Context, serverID uint32, roomID int64, initialSize, scaleUpSizePerTime, scaleDownSizePerTime int32, scaleUpThreshold, scaleDownThreshold float64) (err error) {
	once.Do(func() {
		if serverID == 0 {
			err = fmt.Errorf("bad server id")
			return
		}

		if scaleUpThreshold != 0 && scaleDownThreshold >= scaleUpThreshold {
			err = fmt.Errorf("bad parameter, the scale up threshold must greater than the scale down threshold")
			return
		}

		p := &roomCardPool{
			serverID:             serverID,
			roomID:               roomID,
			initialSize:          initialSize,
			scaleUpSizePerTime:   scaleUpSizePerTime,
			scaleDownSizePerTime: scaleDownSizePerTime,
			scaleUpThreshold:     scaleUpThreshold,
			scaleDownThreshold:   scaleDownThreshold,
			used:                 make(map[string]*card),
		}

		if initialSize <= 0 {
			p.initialSize = DEFAULT_INITIAL_SIZE
		}

		if scaleUpSizePerTime < 0 {
			p.scaleUpSizePerTime = DEFAULT_SCALE_UP_SIZE_PER_TIME
		}

		if scaleDownSizePerTime < 0 {
			p.scaleDownSizePerTime = DEFAULT_SCALE_DOWN_SIZE_PER_TIME
		}

		if scaleUpThreshold <= 0 || scaleUpThreshold > 1 {
			p.scaleUpThreshold = DEFAULT_SCALE_UP_THRESHOLD
		}

		if scaleDownThreshold < 0 || scaleDownThreshold >= 1 {
			p.scaleDownThreshold = DEFAULT_SCALE_DOWN_THRESHOLD
		}

		for {
			select {
			case <-ctx.Done():
				err = fmt.Errorf("called off")
				return
			default:
				if err = p.new(p.currentSize); err != nil {
					logs.Errorf("Obtain room card list failed, err:%s", err.Error())
					time.Sleep(2000 * time.Millisecond)
					continue
				}
			}

			break
		}

		p.freeTimer = time.AfterFunc(RCPOOL_RELEASE_INTERVAL*time.Second, loopFree)
		rcPool = p
	})

	return
}

// loopFree 循环检测释放
func loopFree() {
	rcPool.Lock()
	defer rcPool.Unlock()

	// 房卡池使用率小于等于缩容阀值并且当前房卡池数量大于初值，则进行异步缩容
	if float64(rcPool.freeIdx) <= float64(rcPool.currentSize)*rcPool.scaleDownThreshold && rcPool.currentSize > rcPool.initialSize {
		// logs.Infof("Pool down: used:%v,total:%v,usage:%v,threshold:%v", rcPool.freeIdx, rcPool.currentSize, float64(rcPool.freeIdx)/float64(rcPool.currentSize), rcPool.scaleDownThreshold)
		go rcPool.free(rcPool.currentSize)
	}
	rcPool.freeTimer = time.AfterFunc(RCPOOL_RELEASE_INTERVAL*time.Second, loopFree)
}

// GetRoomCardID 获取房卡
func GetRoomCardID() (string, error) {
	if rcPool == nil {
		return "", fmt.Errorf("please start room card pool first")
	}

	rcPool.Lock()
	defer rcPool.Unlock()

	// 房卡池枯竭，同步申请房卡
	if rcPool.freeIdx >= int(rcPool.currentSize) {
		rcPool.Unlock()
		err := rcPool.new(rcPool.currentSize)
		rcPool.Lock()
		if err != nil {
			return "", fmt.Errorf("run out of pool and new err: %s", err.Error())
		}
	}

	// 性能尖刺，处于回收阶段突然大量房卡申请，房卡池使用完后需要等待回收过程完成
	if rcPool.freeIdx >= rcPool.pendingIdx {
		return "", fmt.Errorf("waiting for free finished")
	}

	// 房卡池使用率达到扩容阀值，异步申请房卡列表
	if float64(rcPool.freeIdx) >= float64(rcPool.currentSize)*rcPool.scaleUpThreshold {
		// logs.Infof("Pool up: used:%v,total:%v,usage:%v,threshold:%v", rcPool.freeIdx, rcPool.currentSize, float64(rcPool.freeIdx)/float64(rcPool.currentSize), rcPool.scaleUpThreshold)
		go rcPool.new(rcPool.currentSize)
	}

	id := rcPool.shared[rcPool.freeIdx]

	rcPool.used[id] = &card{idx: rcPool.freeIdx}

	rcPool.freeIdx++

	return id, nil
}

// ReleaseRoomCardID 释放房卡
func ReleaseRoomCardID(roomCardId string) error {
	if rcPool == nil {
		return fmt.Errorf("please start room card pool first")
	}

	rcPool.Lock()
	defer rcPool.Unlock()

	// 房卡池未被使用
	if rcPool.freeIdx == 0 {
		return fmt.Errorf("multiple release")
	}

	e, ok := rcPool.used[roomCardId]
	if !ok {
		return fmt.Errorf("bad room id")
	}

	rcPool.freeIdx--

	if rcPool.freeIdx != 0 {
		// 由于房卡释放的无序性，须确保已分配房卡与未分配房卡分别处于连续的数组区域内；若只使用了一张房卡，则无须交换位置
		rcPool.shared[rcPool.freeIdx], rcPool.shared[e.idx] = rcPool.shared[e.idx], rcPool.shared[rcPool.freeIdx]
		rcPool.used[rcPool.shared[e.idx]].idx = e.idx
	}

	delete(rcPool.used, roomCardId)

	// 打乱房卡池顺序
	freeSize := rcPool.pendingIdx - rcPool.freeIdx
	if freeSize > 0 {
		idx := rcPool.freeIdx + rand.Intn(freeSize)
		rcPool.shared[rcPool.freeIdx], rcPool.shared[idx] = rcPool.shared[idx], rcPool.shared[rcPool.freeIdx]
	}

	// 房卡池使用率小于等于缩容阀值并且当前房卡池数量大于初值，则进行异步缩容
	if float64(rcPool.freeIdx) <= float64(rcPool.currentSize)*rcPool.scaleDownThreshold && rcPool.currentSize > rcPool.initialSize {
		// logs.Infof("Pool down: used:%v,total:%v,usage:%v,threshold:%v", rcPool.freeIdx, rcPool.currentSize, float64(rcPool.freeIdx)/float64(rcPool.currentSize), rcPool.scaleDownThreshold)
		go rcPool.free(rcPool.currentSize)
	}

	return nil
}

// new 申请房卡列表
func (p *roomCardPool) new(curSize int32) error {
	p.Lock()

	// 房卡池大小如果有变化，直接返回
	if p.currentSize != curSize {
		p.Unlock()
		return nil
	}

	// 房卡池状态为STAT_NEW则说明正在执行new调用返回错误
	if p.status == STAT_NEW {
		p.Unlock()
		return fmt.Errorf("multiple new")
	}
	p.status = STAT_NEW

	scaleUpSizePerTime := p.scaleUpSizePerTime
	// 首次分配
	if p.currentSize == 0 {
		scaleUpSizePerTime = p.initialSize
	}

	// 不进行扩容
	if scaleUpSizePerTime == 0 {
		p.Unlock()
		return nil
	}

	reqMsg := &reqGetRoomCardIDs{
		ServerID: int64(p.serverID),
		RoomID:   p.roomID,
		Num:      int(scaleUpSizePerTime),
	}

	p.Unlock()

	rspMsg := &respGetRoomCardIDs{}
	defer func() {
		p.Lock()
		p.shared = append(p.shared, rspMsg.IDs...)
		p.currentSize += int32(len(rspMsg.IDs))
		p.status = STAT_NOOP
		p.pendingIdx = int(p.currentSize)
		// logs.Infof("Pool: used:%v,total:%v,usage:%v", p.freeIdx, p.currentSize, float64(p.freeIdx)/float64(p.currentSize))
		p.Unlock()
	}()

	body, err := json.Marshal(reqMsg)
	if err != nil {
		return err
	}

	rsp, err := sdk.SendRequest(context.Background(), DISPATCHSERVER, 0, F_ID_ASSIGN_ROOM_CARD_IDS, body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rsp, rspMsg); err != nil {
		return err
	}

	if rspMsg.Code != int(ErrCode(SUCCESS)) {
		return fmt.Errorf("%s", rspMsg.Msg)
	}

	// logs.Infof("Assigned room card list : %v", rspMsg.IDs)
	return nil
}

// free 释放房卡列表
func (p *roomCardPool) free(curSize int32) error {
	p.Lock()

	// 房卡池大小如果有变化直接返回
	if p.currentSize != curSize {
		p.Unlock()
		return nil
	}

	// 若房卡池状态为STAT_FREE则说明正在执行free调用，返回错误
	if p.status == STAT_FREE {
		p.Unlock()
		return fmt.Errorf("multiple free")
	}
	p.status = STAT_FREE

	// 不进行缩容
	if p.scaleDownSizePerTime == 0 {
		p.Unlock()
		return nil
	}

	scaleDownSizePerTime := p.scaleDownSizePerTime

	// 房卡池大小最小为初值
	if p.currentSize-scaleDownSizePerTime < p.initialSize {
		scaleDownSizePerTime = p.currentSize - p.initialSize
	}

	for {
		// 如果缩容后的使用率大于扩容阀值，则减少缩容量
		if float64(int32(p.freeIdx)/(p.currentSize-scaleDownSizePerTime)) >= p.scaleUpThreshold {
			scaleDownSizePerTime /= 2
			continue
		}
		break
	}

	reqMsg := &reqReleaseRoomCardIDs{
		ServerID: int64(p.serverID),
		RoomID:   p.roomID,
		IDs:      p.shared[len(p.shared)-int(scaleDownSizePerTime):],
	}

	p.pendingIdx = len(p.shared) - int(scaleDownSizePerTime)
	p.Unlock()

	defer func() {
		p.Lock()
		p.shared = p.shared[:len(p.shared)-int(scaleDownSizePerTime)]
		p.currentSize = int32(len(p.shared[:]))
		p.status = STAT_NOOP
		p.pendingIdx = int(p.currentSize)
		// logs.Infof("Pool: used:%v,total:%v,usage:%v", p.freeIdx, p.currentSize, float64(p.freeIdx)/float64(p.currentSize))
		p.Unlock()
	}()

	body, err := json.Marshal(reqMsg)
	if err != nil {
		scaleDownSizePerTime = 0
		return err
	}

	rsp, err := sdk.SendRequest(context.Background(), DISPATCHSERVER, 0, F_ID_REVOKE_ROOM_CARD_IDS, body)
	if err != nil {
		scaleDownSizePerTime = 0
		return err
	}

	rspMsg := &respReleaseRoomCardIDs{}
	if err := json.Unmarshal(rsp, rspMsg); err != nil {
		scaleDownSizePerTime = 0
		return err
	}

	if rspMsg.Code != int(ErrCode(SUCCESS)) {
		scaleDownSizePerTime = 0
		return fmt.Errorf("%s", rspMsg.Msg)
	}

	// logs.Infof("Revoke room card list : %v", reqMsg.IDs)
	return nil
}
