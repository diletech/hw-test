package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Создаем новый командный объект.
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	// Устанавливаем пользовательские переменные среды, если они есть.
	for key, value := range env {
		if !value.NeedRemove {
			command.Env = append(command.Env, key+"="+value.Value)
		}
	}

	// Устанавливаем стандартные потоки ввода/вывода/ошибок.
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	// Запускаем команду.
	err := command.Run()
	if err != nil {
		// Если возникла ошибка, проверяем, является ли это ошибкой выполнения команды.
		if exitError, ok := err.(*exec.ExitError); ok { //nolint:errorlint
			// Если это ошибка выполнения, возвращаем код возврата.
			return exitError.ExitCode()
		}
		// Если это другая ошибка, то падаем с 111 без объяснения причин
		os.Exit(111)
	}

	// Если команда выполнена успешно, возврат 0.
	return 0
}
