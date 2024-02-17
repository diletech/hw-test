package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// проверяем, существует ли файл
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	// и является ли файл регулярным
	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// открываем копируемый файл на чтение
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	// получаем значение размера файла
	stat, err := fromFile.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	fileSize := stat.Size()

	fmt.Printf("%v \n", fileSize)

	// проверяем на вилидность offset
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	// выставляем смещение
	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// создаем целевой файл
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	// начинаем определять копируемые данные
	var reader io.Reader // реадер для лимита
	var copySize int64   // расчетный размер копируемых данных для pb
	if limit > 0 {
		// если установлен лимит байтов считываем их колличество
		reader = io.LimitReader(fromFile, limit)
		// а также оприделяем размер копируемых данных для прогресс-бара
		if offset+limit > fileSize {
			copySize = fileSize - offset
		} else {
			copySize = limit
		}
	} else {
		// иначе берём весь файл
		reader = fromFile
		copySize = fileSize
	}

	// инициализируем прогресс-бар
	bar := pb.Full.Start64(copySize)
	defer bar.Finish()
	// и содаем прокси-реадер для прогресс-бара
	barReader := bar.NewProxyReader(reader)

	// копируем содержимое reader в файл
	_, err = io.Copy(toFile, barReader)
	if err != nil {
		return err
	}

	return nil
}
