package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type urlStoreStruct struct {
	Uuid      int    `json:"uuid"`
	Short_url string `json:"short_url"`
	Long_url  string `json:"long_url"`
}

var (
	idCounter int
	counterMu sync.Mutex
)

func generateID() int {
	counterMu.Lock()
	defer counterMu.Unlock()
	idCounter++

	return idCounter
}

type FileStorage struct {
	urlStorage map[string]string
	mu         sync.Mutex
	fileName   string
}

func NewFileStorage(fName string) *FileStorage {
	fs := FileStorage{urlStorage: make(map[string]string), fileName: fName}
	fs.LoadFromFile()
	return &fs
}

func (s *FileStorage) SaveShortURL(shortURL64, longURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urlStorage[shortURL64] = string(longURL)
	return SaveToFile(shortURL64, longURL, s.fileName)
}

func (s *FileStorage) GetOriginalURL(shortURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	longURL, ok := s.urlStorage[shortURL]
	if !ok {
		return "", fmt.Errorf("Not found")
	}
	return longURL, nil
}

func (s *FileStorage) LoadFromFile() error {
	file, err := os.OpenFile(s.fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if len(s.urlStorage) > 0 {
		for k := range s.urlStorage {
			delete(s.urlStorage, k)
		}
	}
	scan := bufio.NewScanner(file)
	s.mu.Lock()
	defer s.mu.Unlock()

	for scan.Scan() {
		storeData := urlStoreStruct{}
		err = json.Unmarshal(scan.Bytes(), &storeData)
		if err != nil {
			return err
		}
		s.urlStorage[storeData.Short_url] = storeData.Long_url
	}

	return nil
}

func SaveToFile(shortURL64, longURL, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	storeData := urlStoreStruct{}
	storeData.Long_url = longURL
	storeData.Short_url = shortURL64
	storeData.Uuid = generateID()

	encoder := json.NewEncoder(file)
	return encoder.Encode(&storeData)
}
