package hw06pipelineexecution

// Типы каналов для удобства использования в пайплайне
type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

// Stage представляет собой функцию-стейдж для пайплайна
type Stage func(in In) (out Out)

// ExecutePipeline запускает конкуррентный пайплайн, состоящий из стейджей
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	// Запуск каждого стейджа в пайплайне
	for _, stage := range stages {
		out = runStage(out, done, stage)
	}

	return out
}

// runStage обрабатывает выполнение отдельного стейджа в горутине
func runStage(in In, done In, stage Stage) Out {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case data, ok := <-in:
				if !ok {
					return
				}

				// Вызов обертки для стейджа с данными из входного канала
				result := stageWrapper(stage, data)
				select {
				case <-done:
					return
				case out <- result:
				}
			}
		}
	}()

	return out
}

// stageWrapper создает обертку для стейджа, позволяя ему принимать и возвращать данные через каналы
func stageWrapper(stage Stage, data interface{}) interface{} {
	out := make(chan interface{})
	in := make(chan interface{}, 1)

	in <- data

	go func() {
		defer close(out)
		close(in)

		// Выполнение стейджа с входным каналом
		stageOut := stage(in)

		// Передача данных из выходного канала стейджа в оберточный канал
		for data := range stageOut {
			out <- data
		}
	}()

	// Возврат результата из оберточного канала
	return <-out
}
