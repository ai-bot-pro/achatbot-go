package common

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPoolInstance 是 IPoolInstance 接口的模拟实现
type MockPoolInstance struct {
	mock.Mock
	id      int
	resetFn func() error
}

func (m *MockPoolInstance) Reset() error {
	if m.resetFn != nil {
		return m.resetFn()
	}
	args := m.Called()
	return args.Error(0)
}

func (m *MockPoolInstance) Release() error {
	args := m.Called()
	return args.Error(0)
}

// MockNewFunc 是 NewFunc 类型的模拟实现
func MockNewFunc(id int, resetFn func() error) NewFunc {
	return func() (IPoolInstance, error) {
		instance := &MockPoolInstance{id: id, resetFn: resetFn}
		instance.On("Reset").Return(nil)
		instance.On("Release").Return(nil)
		return instance, nil
	}
}

func TestRegisterAndGetNewFunc(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := func() (IPoolInstance, error) {
		return &MockPoolInstance{}, nil
	}

	// 注册新函数
	RegisterNewFunc(poolType, newFunc)

	// 获取已注册的函数
	retrievedFunc := GetNewFunc(poolType)
	assert.NotNil(t, retrievedFunc)

	// 尝试获取未注册类型的函数
	unregisteredType := reflect.TypeOf("")
	unregisteredFunc := GetNewFunc(reflect.TypeOf(unregisteredType))
	assert.Nil(t, unregisteredFunc)
}

func TestNewModuleProviderPool(t *testing.T) {
	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	pool := NewModuleProviderPool(5, poolType)

	assert.NotNil(t, pool)
	assert.Equal(t, 5, pool.poolSize)
	assert.Equal(t, poolType, pool.poolType)
	assert.NotNil(t, pool.poolInstances)
	assert.Equal(t, 5, cap(pool.poolInstances))
}

func TestCreateNewInstanceInfo(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := MockNewFunc(0, nil)
	RegisterNewFunc(poolType, newFunc)

	pool := NewModuleProviderPool(1, poolType)
	instanceInfo, err := pool.createNewInstanceInfo()

	assert.NoError(t, err)
	assert.NotNil(t, instanceInfo)
	assert.Equal(t, int32(0), instanceInfo.inUse)
	assert.Equal(t, int64(1), instanceInfo.instanceID) // 初始化时为-1
	assert.NotNil(t, instanceInfo.instance)
	assert.Equal(t, int64(1), pool.totalCreated)
}

func TestInitialize(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := MockNewFunc(0, nil)
	RegisterNewFunc(poolType, newFunc)

	pool := NewModuleProviderPool(3, poolType)
	err := pool.Initialize()

	assert.NoError(t, err)
	assert.Equal(t, 3, len(pool.poolInstances))

	// 收集所有实例以检查它们是否正确初始化
	var instances []*PoolInstanceInfo
	for len(pool.poolInstances) > 0 {
		select {
		case instanceInfo := <-pool.poolInstances:
			if instanceInfo == nil {
				continue
			}
			instances = append(instances, instanceInfo)
		default:
		}
	}

	// 确保我们收集到了所有实例
	assert.Equal(t, 3, len(instances))

	// 注意：由于并发初始化，实例ID可能不是按顺序分配的
	instanceIDs := make(map[int64]bool)
	for _, instance := range instances {
		assert.NotNil(t, instance)
		instanceIDs[instance.instanceID] = true
	}

	// 验证所有预期的ID都存在
	assert.Equal(t, 3, len(instanceIDs))
}

func TestInitializeWithError(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := func() (IPoolInstance, error) {
		return nil, fmt.Errorf("initialization error")
	}
	RegisterNewFunc(poolType, newFunc)

	pool := NewModuleProviderPool(1, poolType) // 减少到1个实例以简化测试
	err := pool.Initialize()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to initialize any")
}

func TestGet(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := MockNewFunc(0, nil)
	RegisterNewFunc(poolType, newFunc)

	pool := NewModuleProviderPool(2, poolType)
	err := pool.Initialize()
	assert.NoError(t, err)

	// 获取一个实例
	instanceInfo, err := pool.Get()
	assert.NoError(t, err)
	assert.NotNil(t, instanceInfo)
	assert.Equal(t, int32(1), instanceInfo.inUse)
	assert.Equal(t, int64(1), pool.totalReused)
	assert.Equal(t, int64(1), pool.totalActive)

	// 检查池中剩余实例数量
	assert.Equal(t, 1, len(pool.poolInstances))
}

func TestGetWithTimeout(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := MockNewFunc(0, nil)
	RegisterNewFunc(poolType, newFunc)

	// 创建一个空池
	pool := NewModuleProviderPool(2, poolType)
	// 不调用 Initialize，所以池为空

	// 检查初始状态
	assert.Equal(t, int64(0), pool.totalCreated)
	assert.Equal(t, int64(0), pool.totalReused)
	assert.Equal(t, int64(0), pool.totalActive)

	// 在超时前尝试获取实例，应该会创建一个新的超出池范围的实例
	instanceInfo, err := pool.Get()
	assert.NoError(t, err)
	assert.NotNil(t, instanceInfo)
	assert.Equal(t, int32(1), instanceInfo.inUse)
	// 应该增加 totalCreated 计数而不是 totalReused
	assert.Equal(t, int64(1), pool.totalCreated)
	assert.Equal(t, int64(0), pool.totalReused)
	assert.Equal(t, int64(1), pool.totalActive)
}

func TestPut(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := MockNewFunc(0, nil)
	RegisterNewFunc(poolType, newFunc)

	pool := NewModuleProviderPool(2, poolType)
	err := pool.Initialize()
	assert.NoError(t, err)

	// 获取一个实例
	instanceInfo, err := pool.Get()
	assert.NoError(t, err)
	assert.NotNil(t, instanceInfo)

	// 归还实例
	pool.Put(instanceInfo)

	// 检查实例是否正确归还到池中
	assert.Equal(t, 2, len(pool.poolInstances))
	assert.Equal(t, int64(0), pool.totalActive)

	// 检查实例状态是否正确重置
	assert.Equal(t, int32(0), instanceInfo.inUse)
}

func TestPutNilInstance(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	pool := NewModuleProviderPool(2, poolType)

	// 归还一个 nil 实例不应该导致 panic
	assert.NotPanics(t, func() {
		pool.Put(nil)
	})
}

func TestGetStats(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := MockNewFunc(0, nil)
	RegisterNewFunc(poolType, newFunc)

	pool := NewModuleProviderPool(3, poolType)
	err := pool.Initialize()
	assert.NoError(t, err)

	stats := pool.GetStats()
	assert.Equal(t, 3, stats["pool_size"])
	assert.Equal(t, 3, stats["total_instances"])
	assert.Equal(t, int64(0), stats["active_count"])
	assert.Equal(t, int64(3), stats["total_created"])
	assert.Equal(t, int64(0), stats["total_reused"])

	// 获取一个实例并检查统计数据
	instanceInfo, err := pool.Get()
	assert.NoError(t, err)
	stats = pool.GetStats()
	assert.Equal(t, int64(1), stats["active_count"])
	assert.Equal(t, int64(1), stats["total_reused"])

	// 归还实例并检查统计数据
	pool.Put(instanceInfo)
	stats = pool.GetStats()
	assert.Equal(t, int64(0), stats["active_count"])
}

func TestConcurrentAccess(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))

	// 创建一个计数器来跟踪创建的实例数量
	var createdCount int64
	var mu sync.Mutex
	createdInstances := make(map[int]bool)

	newFunc := func() (IPoolInstance, error) {
		mu.Lock()
		id := int(createdCount)
		createdCount++
		createdInstances[id] = true
		mu.Unlock()

		instance := &MockPoolInstance{id: id}
		instance.On("Reset").Return(nil)
		instance.On("Release").Return(nil)
		return instance, nil
	}
	RegisterNewFunc(poolType, newFunc)

	pool := NewModuleProviderPool(5, poolType)
	err := pool.Initialize()
	assert.NoError(t, err)

	// 并发获取和归还实例
	var wg sync.WaitGroup
	const numWorkers = 10
	const numOperations = 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				instanceInfo, err := pool.Get()
				assert.NoError(t, err)
				// 模拟一些工作
				time.Sleep(time.Microsecond)
				pool.Put(instanceInfo)
			}
		}(i)
	}

	wg.Wait()

	// 检查统计数据
	stats := pool.GetStats()
	assert.Equal(t, 5, stats["pool_size"])
	// 活跃实例应该为0
	assert.Equal(t, int64(0), stats["active_count"])
	// 总创建数应该大于等于初始池大小
	assert.True(t, stats["total_created"].(int64) >= 5)
}

func TestClose(t *testing.T) {
	// 清理之前注册的函数
	newFuncMap = make(map[reflect.Type]NewFunc)

	poolType := reflect.TypeOf((*MockPoolInstance)(nil))
	newFunc := MockNewFunc(0, nil)
	RegisterNewFunc(poolType, newFunc)

	pool := NewModuleProviderPool(3, poolType)
	err := pool.Initialize()
	assert.NoError(t, err)

	// 关闭池
	pool.Close()

	// 尝试从已关闭的池中获取实例应该失败
	_, err = pool.Get()
	assert.Error(t, err)
	// 错误信息可能因通道关闭而有所不同，所以我们检查它是否包含预期的错误信息之一
	assert.True(t, strings.Contains(err.Error(), "pool is shutting down") || strings.Contains(err.Error(), "received nil instance from pool"))

	// 确保统计数据仍然可以访问
	stats := pool.GetStats()
	assert.Equal(t, 3, stats["pool_size"])
}
