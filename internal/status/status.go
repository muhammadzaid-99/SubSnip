// internal/status/status.go
package status

import "sync"

var (
	mu    sync.RWMutex
	store map[string]string
)

func Init() {
	if store == nil {
		store = make(map[string]string)
	}
}

func Set(taskID, value string) {
	mu.Lock()
	defer mu.Unlock()
	store[taskID] = value
}

func Get(taskID string) string {
	mu.RLock()
	defer mu.RUnlock()
	return store[taskID]
}
