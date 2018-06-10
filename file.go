package main

import (
	"os"

	"github.com/go-mego/mego"
)

const (
	// B 為 Byte。
	B = 1
	// KB 為 Kilobyte。
	KB = 1024 * B
	// MB 為 Megabyte。
	MB = 1024 * KB
	// GB 為 Gigabyte。
	GB = 1024 * MB
	// TB 為 Terabyte。
	TB = 1024 * GB
)

//
func processor(chunk []byte) error {

}

// New 會建立一個檔案處理模組，可供安插於 Mego 引擎中。
func New() mego.HandlerFunc {
	return func(c *mego.Context) {
		c.Map(&Store{
			Processor: processor,
		})
	}
}

// Store 是檔案處理的主要存儲建構體。
type Store struct {
	Processor func(chunk []byte) error
}

// Get 會取得指定欄位的檔案資訊。
func (s *Store) Get(field string) (file *File, err error) {

}

// MustGet 和 `Get` 相同，但沒有該欄位可取得檔案時會呼叫 `panic`。
func (s *Store) MustGet(field string) (file *File) {

}

// GetMulti 會取得指定欄位中的多個檔案。
func (s *Store) GetMulti(field string) (files []*File, err error) {

}

// MustGetMulti 和 `GetMulti` 相同，但沒有該欄位可取得檔案時會呼叫 `panic`。
func (s *Store) MustGetMulti(field string) (files []*File) {

}

// Serve 可以接收靜態檔案路徑、*os.File、ioutil.ReadFile 並提供靜態檔案給瀏覽器，就像點擊按鈕會下載檔案那樣。
func (s *Store) Serve(file interface{}) (err error) {

}

// File 呈現了一個檔案與其資訊。
type File struct {
	// Name 為此檔案的原生名稱。
	Name string
	// Size 是這個檔案的總位元組大小。
	Size int
	// Extension 是這個檔案的副檔名。
	Extension string
	// Path 為此檔案上傳後的本地路徑。
	Path string
	// Keys 為此檔案的鍵值組，可供開發者存放自訂資料。
	Keys map[string]interface{}
}

// Remove 會移除這個檔案。
func (f *File) Remove() error {
	return os.Remove(f.Path)
}

// Move 會移動接收到的檔案到指定路徑。
func (f *File) Move(dest string) error {
	return os.Rename(f.Path, dest)
}
