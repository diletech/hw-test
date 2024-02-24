package main

import (
	"fmt"
	"os"
)

func main() {
	// Проверяем на наличие обязательных двух аргументов.
	if len(os.Args) <= 2 {
		fmt.Printf("Usage: %s <pathToEnvDir> <runCommand> [<args1>|<args2>]\n", os.Args[0])
		return
	}

	// Формируем мапу из всех системных переменных окуржения,
	// изменяя и добавляя по правилам обработки указанной дирректории
	mapEnv, err := ReadDir(os.Args[1])
	if err != nil {
		// если ошибка обработки, то выводим в stderr и код возврата 111
		fmt.Fprintf(os.Stderr, "%s --> Break prepare '%s': %v\n", os.Args[0], os.Args[1], err)
		os.Exit(111)
	}

	// Запуск команнды с аргументами и новым окружением,
	// полчуем код возврата и с ним выходим.
	code := RunCmd(os.Args[2:], mapEnv)
	os.Exit(code)
}
