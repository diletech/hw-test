package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run выполняет задачи в n параллельных горутинах с ограничением на количество ошибок m
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var once sync.Once
	var errorsCount int32                 // счетчик ошибок, используем атомарные операции для безопасного доступа
	taskCh := make(chan Task, len(tasks)) // канал для передачи задач в горутины
	stopCh := make(chan struct{})         // канал для сигнала горутинам остановить выполнение
	doneCh := make(chan struct{})         // канал для сигнала, что все горутины завершили свою работу

	// Переменная с функцией для горутины по выполнению задания
	worker := func(once *sync.Once) {
		defer wg.Done()
		for {
			select {
			// Выполняем задание
			case task, ok := <-taskCh:
				if !ok {
					// taskCh закрыт, завершаем работу
					return
				}
				err := task()
				// Увеличиваем счетчик ошибок и проверяем, не превышен ли лимит
				if err != nil {
					atomic.AddInt32(&errorsCount, 1)
					if int(atomic.LoadInt32(&errorsCount)) >= m && m != 0 {
						once.Do(func() {
							close(stopCh)
						})
						return
					}
				}
			// остановка горутины при получении сигнала
			case <-stopCh:
				return
			}
		}
	}

	// Запускаем n горутин для выполнения заданий
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&once)
	}

	// Горутина для передачи заданий в канал и прерывния по сигналу stopCh в случае достижения m ошибок
	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case taskCh <- task:
			case <-stopCh:
				return
			}
		}
	}()

	// Горутина для ожидания завершения всех горутин и сигнализации о успешном их выполнении
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	// После запуска всех горутин ожидаем их завершения
	select {
	case <-doneCh:
		// Все горутины завершили свою работу успешно или нет если m=0
		return nil
	case <-stopCh:
		// Функция была остановлена из-за превышения лимита ошибок
		return ErrErrorsLimitExceeded
	}
}
