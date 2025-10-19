package common

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/weedge/pipeline-go/pkg/logger"
)

type NewFunc func() (IPoolInstance, error)

var newFuncMap = make(map[reflect.Type]NewFunc)

func RegisterNewFunc(poolType reflect.Type, newFunc NewFunc) {
	newFuncMap[poolType] = newFunc
}
func GetNewFunc(poolType reflect.Type) NewFunc {
	if newFunc, ok := newFuncMap[poolType]; ok {
		return newFunc
	}
	return nil
}

type PoolInstanceInfo struct {
	instanceID int64
	inUse      int32
	lastUsed   int64
	instance   IPoolInstance
}

// ModuleProviderPool module provider 资源池
type ModuleProviderPool struct {
	poolInstances chan *PoolInstanceInfo
	poolSize      int
	poolType      reflect.Type
	newFunc       NewFunc

	// stats info
	totalCreated int64 // init create or beyond create +1
	totalReused  int64 // get from poolInstances channel to reuse +1
	totalActive  int64 // get from poolInstances channel become active +1 and put to poolInstances channel become disactive -1

	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewModuleProviderPool 创建新的资源池
func NewModuleProviderPool(poolSize int, poolType reflect.Type) *ModuleProviderPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &ModuleProviderPool{
		poolInstances: make(chan *PoolInstanceInfo, poolSize),
		poolSize:      poolSize,
		poolType:      poolType,
		newFunc:       GetNewFunc(poolType),
		ctx:           ctx,
		cancel:        cancel,
	}

	return pool
}

// createNewInstance 创建新的实例
func (p *ModuleProviderPool) createNewInstanceInfo() (*PoolInstanceInfo, error) {
	instance, err := p.newFunc()
	if err != nil {
		return nil, err
	}

	instanceInfo := &PoolInstanceInfo{
		instance:   instance,
		lastUsed:   time.Now().UnixNano(),
		inUse:      0,
		instanceID: atomic.AddInt64(&p.totalCreated, 1),
	}

	logger.Infof("Created New InstanceInfo")
	return instanceInfo, nil
}

// Initialize 并行初始化池
func (p *ModuleProviderPool) Initialize() error {
	logger.Infof("Initializing pool with %d instances...", p.poolSize)
	if p.newFunc == nil {
		return fmt.Errorf("no NewFunc registered for pool type %s", p.poolType)
	}

	var initWg sync.WaitGroup
	errorChan := make(chan error, p.poolSize)
	for i := 0; i < p.poolSize; i++ {
		initWg.Add(1)
		go func(instanceID int) {
			defer initWg.Done()

			instanceInfo, err := p.createNewInstanceInfo()
			if err != nil {
				errorChan <- fmt.Errorf("createNewInstanceInfo err: %s", err.Error())
				return
			}

			select {
			case p.poolInstances <- instanceInfo:
				logger.Infof("%s Instance#%d Initialized", p.poolType, instanceID)
			default:
				instanceInfo.instance.Release()
				errorChan <- fmt.Errorf("pool queue full, %s instance#%d release", p.poolType, instanceID)
			}
		}(i)
	}

	initWg.Wait()
	close(errorChan)

	var initErrors []error
	for err := range errorChan {
		if err != nil {
			initErrors = append(initErrors, err)
			logger.Warnf("Initialization warning: %v", err)
		}
	}

	successCount := p.poolSize - len(initErrors)
	logger.Infof("pool initialized with %d/%d %s instances", successCount, p.poolSize, p.poolType)

	if len(initErrors) > 0 && successCount == 0 {
		return fmt.Errorf("failed to initialize any %s instances", p.poolType)
	}

	return nil
}

// Get 获取实例
func (p *ModuleProviderPool) Get() (*PoolInstanceInfo, error) {
	logger.Infof("Attempting to get %s instance from pool (available: %d)", p.poolType, len(p.poolInstances))

	select {
	case instanceInfo := <-p.poolInstances:
		if instanceInfo == nil {
			return nil, fmt.Errorf("received nil instance from pool")
		}
		logger.Infof("Got %s instanceInfoInfo#%d from pool", p.poolType, instanceInfo.instanceID)
		if atomic.CompareAndSwapInt32(&instanceInfo.inUse, 0, 1) {
			instanceInfo.lastUsed = time.Now().UnixNano()
			atomic.AddInt64(&p.totalReused, 1)
			atomic.AddInt64(&p.totalActive, 1)
			logger.Infof("%s instance#%d marked as in-use (active: %d)", p.poolType, instanceInfo.instanceID, atomic.LoadInt64(&p.totalActive))
			return instanceInfo, nil
		}
		logger.Warnf("%s instance#%d already in use, returning to pool", p.poolType, instanceInfo.instanceID)
		select {
		case p.poolInstances <- instanceInfo:
		default:
		}
		return p.Get() // 递归重试
	case <-time.After(100 * time.Millisecond):
		logger.Warnf("pool timeout, creating new beyond instance")
		instanceInfo, err := p.createNewInstanceInfo()
		if err != nil {
			return nil, fmt.Errorf("createNewInstanceInfo err: %s", err.Error())
		}
		instanceInfo.inUse = 1
		return instanceInfo, nil
	case <-p.ctx.Done():
		return nil, fmt.Errorf("pool is shutting down")
	}
}

// Put 归还实例
func (p *ModuleProviderPool) Put(instanceInfo *PoolInstanceInfo) {
	if instanceInfo == nil {
		logger.Warnf("Attempted to put nil instance")
		return
	}

	logger.Infof("Returning instance#%d to pool(%s)", instanceInfo.instanceID, p.poolType)

	if atomic.CompareAndSwapInt32(&instanceInfo.inUse, 1, 0) {
		instanceInfo.lastUsed = time.Now().UnixNano()

		// instance reset
		if err := instanceInfo.instance.Reset(); err != nil {
			logger.Warnf("Failed to reset instance#%d to pool(%s) err: %v", instanceInfo.instanceID, p.poolType, err)
		}

		select {
		case p.poolInstances <- instanceInfo:
			atomic.AddInt64(&p.totalActive, -1)
			logger.Infof("instance#%d marked as available (active: %d) to pool(%s)", instanceInfo.instanceID, atomic.LoadInt64(&p.totalActive), p.poolType)
		default:
			logger.Warnf("pool queue full, release instance#%d to pool(%s)", instanceInfo.instanceID, p.poolType)
			instanceInfo.instance.Release()
			instanceInfo = nil // gc to free, maybe use sync.pool
		}
	} else {
		logger.Warnf("instance#%d was not in use, cannot return", instanceInfo.instanceID)
	}
}

// GetStats 获取统计信息
func (p *ModuleProviderPool) GetStats() map[string]any {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]any{
		"pool_size":       p.poolSize,
		"total_instances": len(p.poolInstances),
		"active_count":    atomic.LoadInt64(&p.totalActive),
		"total_created":   atomic.LoadInt64(&p.totalCreated),
		"total_reused":    atomic.LoadInt64(&p.totalReused),
	}
}

// Shutdown 关闭池
func (p *ModuleProviderPool) Close() {
	logger.Infof("Pool Closing...")

	p.cancel()

	p.mu.Lock()
	defer p.mu.Unlock()

	// Drain the channel and release all instances
	for {
		select {
		case instanceInfo := <-p.poolInstances:
			if instanceInfo != nil && instanceInfo.instance != nil {
				instanceInfo.instance.Release()
			}
		default:
			// Channel is empty
			close(p.poolInstances)
			logger.Infof("Pool Closed")
			return
		}
	}
}
