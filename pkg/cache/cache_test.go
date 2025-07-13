package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestCache() ICache {
	return &Cache{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func TestCache_SetAndGet(t *testing.T) {
	t.Parallel()

	cache := createTestCache()

	cache.Set("test_key", "test_value", time.Minute)
	value, found := cache.Get("test_key")

	assert.True(t, found)
	assert.Equal(t, "test_value", value)
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()

	cache := createTestCache()

	cache.Set("delete_key", "delete_value", time.Minute)
	cache.Delete("delete_key")

	_, found := cache.Get("delete_key")
	assert.False(t, found)
}

func TestCache_Flush(t *testing.T) {
	t.Parallel()

	cache := createTestCache()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)

	cache.Flush()

	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")

	assert.False(t, found1)
	assert.False(t, found2)
}

func TestCache_GetNonExistentKey(t *testing.T) {
	t.Parallel()

	cache := createTestCache()

	value, found := cache.Get("non_existent_key")

	assert.False(t, found)
	assert.Nil(t, value)
}

func TestCache_Expiration(t *testing.T) {
	t.Parallel()

	cache := createTestCache()

	cache.Set("expire_key", "expire_value", 1*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	_, found := cache.Get("expire_key")
	assert.False(t, found)
}

func TestCache_Singleton(t *testing.T) {

	mu := &sync.Mutex{}
	mu.Lock()
	originalInstance := instance
	instance = nil
	once = sync.Once{}
	mu.Unlock()

	defer func() {
		mu.Lock()
		instance = originalInstance
		once = sync.Once{}
		mu.Unlock()
	}()

	cache1 := NewCache(5*time.Minute, 10*time.Minute)
	cache2 := NewCache(10*time.Minute, 20*time.Minute)

	assert.Equal(t, cache1, cache2)
}

func TestGetInstance(t *testing.T) {

	mu := &sync.Mutex{}
	mu.Lock()
	originalInstance := instance
	instance = nil
	once = sync.Once{}
	mu.Unlock()

	defer func() {
		mu.Lock()
		instance = originalInstance
		once = sync.Once{}
		mu.Unlock()
	}()

	cache := GetInstance()
	require.NotNil(t, cache)

	cache2 := GetInstance()
	assert.Equal(t, cache, cache2)
}

func TestCache_SetWithDifferentTypes(t *testing.T) {
	t.Parallel()

	cache := createTestCache()

	cache.Set("string_key", "string_value", time.Minute)
	cache.Set("int_key", 42, time.Minute)
	cache.Set("bool_key", true, time.Minute)
	cache.Set("slice_key", []string{"a", "b", "c"}, time.Minute)

	strValue, found := cache.Get("string_key")
	assert.True(t, found)
	assert.Equal(t, "string_value", strValue)

	intValue, found := cache.Get("int_key")
	assert.True(t, found)
	assert.Equal(t, 42, intValue)

	boolValue, found := cache.Get("bool_key")
	assert.True(t, found)
	assert.Equal(t, true, boolValue)

	sliceValue, found := cache.Get("slice_key")
	assert.True(t, found)
	assert.Equal(t, []string{"a", "b", "c"}, sliceValue)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	cache := createTestCache()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", index)
			value := fmt.Sprintf("value_%d", index)
			cache.Set(key, value, time.Minute)
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", index)
			cache.Get(key)
		}(i)
	}

	wg.Wait()

	value, found := cache.Get("key_0")
	if found {
		assert.Equal(t, "value_0", value)
	}
}
