package tache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync/atomic"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/xhofe/gsync"
)

// Manager is the manager of all tasks
type Manager[T Task] struct {
	tasks           gsync.MapOf[string, T]
	queue           gsync.QueueOf[T]
	workers         *WorkerPool[T]
	opts            *Options
	debouncePersist func()
	running         atomic.Bool

	idGenerator func() string
	logger      *slog.Logger
}

// NewManager create a new manager
func NewManager[T Task](opts ...Option) *Manager[T] {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	m := &Manager[T]{
		workers: NewWorkerPool[T](options.Works),
		opts:    options,
		idGenerator: func() string {
			return nanoid.Must()
		},
		logger: options.Logger,
	}
	m.running.Store(options.Running)
	if m.opts.PersistPath != "" || (m.opts.PersistReadFunction != nil && m.opts.PersistWriteFunction != nil) {
		m.debouncePersist = func() {
			_ = m.persist()
		}
		if m.opts.PersistDebounce != nil {
			m.debouncePersist = newDebounce(func() {
				_ = m.persist()
			}, *m.opts.PersistDebounce)
		}
		err := m.recover()
		if err != nil {
			m.logger.Error("recover error", "error", err)
		}
	} else {
		m.debouncePersist = func() {}
	}
	return m
}

// Add a task to manager
func (m *Manager[T]) Add(task T) {
	ctx, cancel := context.WithCancel(context.Background())
	task.SetCtx(ctx)
	task.SetCancelFunc(cancel)
	task.SetPersist(m.debouncePersist)
	if task.GetID() == "" {
		task.SetID(m.idGenerator())
	}
	if _, maxRetry := task.GetRetry(); maxRetry == 0 {
		task.SetRetry(0, m.opts.MaxRetry)
	}
	if sliceContains([]State{StateRunning}, task.GetState()) {
		task.SetState(StatePending)
	}
	if sliceContains([]State{StateCanceling}, task.GetState()) {
		task.SetState(StateCanceled)
		task.SetErr(context.Canceled)
	}
	if task.GetState() == StateFailing {
		task.SetState(StateFailed)
	}
	m.tasks.Store(task.GetID(), task)
	if !sliceContains([]State{StateSucceeded, StateCanceled, StateErrored, StateFailed}, task.GetState()) {
		m.queue.Push(task)
	}
	m.debouncePersist()
	m.next()
}

// get next task from queue and execute it
func (m *Manager[T]) next() {
	// if manager is not running, return
	if !m.running.Load() {
		return
	}
	// if workers is full, return
	worker := m.workers.Get()
	if worker == nil {
		return
	}
	m.logger.Debug("got worker", "id", worker.ID)
	task, err := m.queue.Pop()
	// if cannot get task, return
	if err != nil {
		m.workers.Put(worker)
		return
	}
	m.logger.Debug("got task", "id", task.GetID())
	go func() {
		defer func() {
			if task.GetState() == StateWaitingRetry {
				m.queue.Push(task)
			}
			m.workers.Put(worker)
			m.next()
		}()
		if task.GetState() == StateCanceling {
			task.SetState(StateCanceled)
			task.SetErr(context.Canceled)
			return
		}
		if m.opts.Timeout != nil {
			ctx, cancel := context.WithTimeout(task.Ctx(), *m.opts.Timeout)
			defer cancel()
			task.SetCtx(ctx)
		}
		m.logger.Info("worker execute task", "worker", worker.ID, "task", task.GetID())
		worker.Execute(task)
	}()
}

// Wait wait all tasks done, just for test
func (m *Manager[T]) Wait() {
	for {
		tasks, running := m.queue.Len(), m.workers.working
		if tasks == 0 && running == 0 {
			return
		}
		runtime.Gosched()
	}
}

// persist all tasks
func (m *Manager[T]) persist() error {
	if m.opts.PersistPath == "" && m.opts.PersistReadFunction == nil && m.opts.PersistWriteFunction == nil {
		return nil
	}
	// serialize tasks
	tasks := m.GetAll()
	var toPersist []T
	for _, task := range tasks {
		// only persist task which is not persistable or persistable and need persist
		if p, ok := Task(task).(Persistable); !ok || p.Persistable() {
			toPersist = append(toPersist, task)
		}
	}
	marshal, err := json.Marshal(toPersist)
	if err != nil {
		return err
	}
	if m.opts.PersistReadFunction != nil && m.opts.PersistWriteFunction != nil {
		err = m.opts.PersistWriteFunction(marshal)
		if err != nil {
			return err
		}
	}
	if m.opts.PersistPath != "" {
		// write to file
		err = os.WriteFile(m.opts.PersistPath, marshal, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// recover all tasks
func (m *Manager[T]) recover() error {
	var data []byte
	var err error
	if m.opts.PersistPath != "" {
		// read from file
		data, err = os.ReadFile(m.opts.PersistPath)
	} else if m.opts.PersistReadFunction != nil && m.opts.PersistWriteFunction != nil {
		data, err = m.opts.PersistReadFunction()
	} else {
		return nil
	}
	if err != nil {
		return err
	}
	// deserialize tasks
	var tasks []T
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return err
	}
	// add tasks
	for _, task := range tasks {
		// only recover task which is not recoverable or recoverable and need recover
		if r, ok := Task(task).(Recoverable); !ok || r.Recoverable() {
			m.Add(task)
		} else {
			task.SetState(StateFailed)
			task.SetErr(fmt.Errorf("the task is interrupted and cannot be recovered"))
			m.tasks.Store(task.GetID(), task)
		}
	}
	return nil
}

// Cancel a task by ID
func (m *Manager[T]) Cancel(id string) {
	if task, ok := m.tasks.Load(id); ok {
		task.Cancel()
		m.debouncePersist()
	}
}

// CancelAll cancel all tasks
func (m *Manager[T]) CancelAll() {
	m.tasks.Range(func(key string, value T) bool {
		value.Cancel()
		return true
	})
	m.debouncePersist()
}

// CancelByCondition cancel tasks under specific condition given by a function
func (m *Manager[T]) CancelByCondition(condition func(task T) bool) {
	m.tasks.Range(func(key string, value T) bool {
		if condition(value) {
			value.Cancel()
		}
		return true
	})
}

// GetAll get all tasks
func (m *Manager[T]) GetAll() []T {
	var tasks []T
	m.tasks.Range(func(key string, value T) bool {
		tasks = append(tasks, value)
		return true
	})
	return tasks
}

// GetByID get task by ID
func (m *Manager[T]) GetByID(id string) (T, bool) {
	return m.tasks.Load(id)
}

// GetByState get tasks by state
func (m *Manager[T]) GetByState(state ...State) []T {
	return m.GetByCondition(func(value T) bool {
		return sliceContains(state, value.GetState())
	})
}

// GetByCondition get tasks under specific condition given by a function
func (m *Manager[T]) GetByCondition(condition func(task T) bool) []T {
	var tasks []T
	m.tasks.Range(func(key string, value T) bool {
		if condition(value) {
			tasks = append(tasks, value)
		}
		return true
	})
	return tasks
}

// Remove a task by ID
func (m *Manager[T]) Remove(id string) {
	m.tasks.Delete(id)
	m.debouncePersist()
}

// RemoveAll remove all tasks
func (m *Manager[T]) RemoveAll() {
	tasks := m.GetAll()
	for _, task := range tasks {
		m.Remove(task.GetID())
	}
}

// RemoveByState remove tasks by state
func (m *Manager[T]) RemoveByState(state ...State) {
	tasks := m.GetByState(state...)
	for _, task := range tasks {
		m.Remove(task.GetID())
	}
}

// RemoveByCondition remove tasks under specific condition given by a function
func (m *Manager[T]) RemoveByCondition(condition func(task T) bool) {
	tasks := m.GetByCondition(condition)
	for _, task := range tasks {
		m.Remove(task.GetID())
	}
}

// Retry a task by ID
func (m *Manager[T]) Retry(id string) {
	if task, ok := m.tasks.Load(id); ok {
		task.SetState(StateWaitingRetry)
		task.SetErr(nil)
		task.SetRetry(0, m.opts.MaxRetry)
		m.queue.Push(task)
		m.next()
		m.debouncePersist()
	}
}

// RetryAllFailed retry all failed tasks
func (m *Manager[T]) RetryAllFailed() {
	tasks := m.GetByState(StateFailed)
	for _, task := range tasks {
		m.Retry(task.GetID())
	}
}

// Start manager
func (m *Manager[T]) Start() {
	m.running.Store(true)
	m.next()
}

// Pause manager
func (m *Manager[T]) Pause() {
	m.running.Store(false)
}

func (m *Manager[T]) SetWorkersNumActive(active int64) {
	oldActive := m.workers.numActive
	m.workers.SetNumActive(active)
	for i := 0; i < int(active-oldActive) && i < m.queue.Len(); i++ {
		m.next()
	}
}
