package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run выполняет задачи в n параллельных горутинах с ограничением на количество ошибок m.
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
			// получаем функцию на выполнени
			case task, ok := <-taskCh:
				if !ok {
					// taskCh закрыт, завершаем работу
					fmt.Println("<- stop taskCh закрыт")
					return
				}
				fmt.Println("-> start")
				// проверяем, не превышен ли лимит
				if int(atomic.LoadInt32(&errorsCount)) >= m && m != 0 {
					once.Do(func() {
						close(stopCh)
					})
					fmt.Println("<- stop лимит ошибок")
					return
				}
				// выполняем задание
				err := task()
				// если ошибка, то увеличиваем счетчик ошибок
				if err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
			// остановка горутины при получении сигнала
			case <-stopCh:
				fmt.Println("<- stop канал stopCh закрыт")
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
		// и тут необходимо дождаться пока все воркеры завершаться
		// иначе возможны racecond в вызывающей таске\тесте если вдруг что-то не атомарно и не потокобезапасно
		// (хотя как кажется тут этот wg.Wait искуственно тормозит общее выполнение программы)
		wg.Wait()
		return ErrErrorsLimitExceeded
	}
}
