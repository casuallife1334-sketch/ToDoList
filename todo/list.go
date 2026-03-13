package todo

import "sync"

type List struct {
	tasks map[string]Task
	// На этом уровне добавляем mutex, чтобы избежать гонки данных
	mtx sync.RWMutex
}

func NewList() *List {
	return &List{
		tasks: make(map[string]Task),
	}
}

func (l *List) AddTask(task Task) error {
	// Используем мапу мы в КАЖДОМ методе структуры list,
	// поэтому, конкурентный доступ будет именно к ней
	// нужно его отрегулировать
	l.mtx.Lock()

	// При этом, независимо от того, когда у нас закончится выполнение
	// данного метода, удобнее будет задеферить mtx.Unlock()
	defer l.mtx.Unlock()

	if _, ok := l.tasks[task.Title]; ok {
		return ErrTaskAlreadyExists
	}

	l.tasks[task.Title] = task

	return nil
}

func (l *List) GetTask(title string) (Task, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (l *List) ListTask() map[string]Task {
	// Тут не нужно брать Lock, потому что область видимости переменной
	// tmp - только этот метод
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	tmp := make(map[string]Task, len(l.tasks))

	for k, v := range l.tasks {
		tmp[k] = v
	}

	return tmp
}

func (l *List) ListUncompletedTasks() map[string]Task {
	// То же, что и в прошлом методе
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	UncompletedTasks := make(map[string]Task)

	for title, task := range l.tasks {
		if !task.Completed {
			UncompletedTasks[title] = task
		}
	}

	return UncompletedTasks
}

func (l *List) CompleteTask(title string) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	task.Complete()

	l.tasks[title] = task

	return l.tasks[title], nil
}

func (l *List) UncompleteTask(title string) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	task.Uncomplete()

	l.tasks[title] = task

	return l.tasks[title], nil
}

func (l *List) DeleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.tasks[title]; !ok {
		return ErrTaskNotFound
	}

	delete(l.tasks, title)

	return nil
}
