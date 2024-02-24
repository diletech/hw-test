package main

import (
	"os"
	"path/filepath"
	"testing"
)

var dirtst = "mytestdata"

func TestReadDir(t *testing.T) {
	// Подготовка тестового окружения
	testEnvSetup(t, dirtst)
	defer testEnvTeardown(t, dirtst)
	// defer os.RemoveAll(dirtst)

	// Добавление переменных в окружение
	tEnvMap := map[string]string{
		"EMPTY": "replace",
		"UNSET": "some_value",
	}
	for k, v := range tEnvMap {
		err := os.Setenv(k, v)
		if err != nil {
			t.Fatalf("Error setting environment variable: %v", err)
		}
		defer os.Unsetenv(k) // valid use go v1.22
	}

	// Вызов функции ReadDir
	env, err := ReadDir(dirtst)
	if err != nil {
		t.Fatalf("ReadDir(%s) returned error: %v", dirtst, err)
	}

	// Проверка результатов
	assertEnvValue(t, env, "ADD", "value1")
	assertEmptyValue(t, env, "EMPTY")
	assertNeedRemove(t, env, "UNSET")
}

func testEnvSetup(t *testing.T, dir string) {
	t.Helper()
	// Проверка существования директории и создание, если не существует
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0o755)
		if err != nil {
			t.Fatalf("Error creating test directory: %v", err)
		}
	} else if err != nil {
		t.Fatalf("Error checking test directory: %v", err)
	}

	// Создание тестовых файлов с данными
	createTestFile(t, dir, "ADD", "value1")
	createTestFile(t, dir, "EMPTY", "\n")
	createTestFile(t, dir, "UNSET", "")
}

func createTestFile(t *testing.T, dir string, name, content string) {
	t.Helper()
	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("Error creating test file %s: %v", name, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Error writing content to test file %s: %v", name, err)
	}
}

func testEnvTeardown(t *testing.T, dir string) {
	t.Helper()
	// Удаление тестовой директории и всех файлов в ней
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Error cleaning up test environment: %v", err)
	}
}

func assertEnvValue(t *testing.T, env Environment, key, expectedValue string) {
	t.Helper()
	if value, ok := env[key]; !ok || value.Value != expectedValue {
		t.Errorf("Expected value for key '%s' to be '%s', got '%s'", key, expectedValue, value.Value)
	}
}

func assertNeedRemove(t *testing.T, env Environment, key string) {
	t.Helper()
	if value, ok := env[key]; !ok || !value.NeedRemove {
		t.Errorf("Expected NeedRemove to be true for key '%s'", key)
	}
}

func assertEmptyValue(t *testing.T, env Environment, key string) {
	t.Helper()
	if value, ok := env[key]; !ok || value.Value != "" {
		t.Errorf("Expected value for key '%s' to be empty", key)
	}
}
