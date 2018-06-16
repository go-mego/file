package file

import (
	"errors"
	"io"
	"io/ioutil"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-mego/mego"
)

var (
	ErrNotFound          = errors.New("file: no such file")
	ErrUnknownFileReader = errors.New("file: unknown file reader")
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

// New 會建立一個檔案處理模組，可供安插於 Mego 引擎中。
func New(option ...*Options) mego.HandlerFunc {
	o := &Options{
		MaxMemory: 24 * KB,
	}
	if len(option) > 0 {
		o = option[0]
	}
	return func(c *mego.Context) {
		s := &Store{
			context: c,
			files:   make(map[string][]*File),
			options: o,
		}
		c.Map(s)
	}
}

//
type Options struct {
	MaxMemory int
}

// Store 是檔案處理的主要存儲建構體。
type Store struct {
	context *mego.Context
	parsed  bool
	options *Options
	files   map[string][]*File
}

//
func (s *Store) parse() error {
	if s.parsed {
		return nil
	}
	err := s.context.Request.ParseMultipartForm(int64(s.options.MaxMemory))
	if nil != err {
		return err
	}
	for field, headers := range s.context.Request.MultipartForm.File {
		for _, header := range headers {
			file, err := header.Open()
			if err != nil {
				return err
			}
			tmpFile, err := ioutil.TempFile(os.TempDir(), "")
			if err != nil {
				return err
			}
			_, err = io.Copy(tmpFile, file)
			if err != nil {
				return err
			}
			path := tmpFile.Name()
			if err != nil {
				return err
			}
			s.files[field] = append(s.files[field], &File{
				Name:      strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)),
				Size:      int(header.Size),
				Extension: strings.TrimLeft(filepath.Ext(header.Filename), "."),
				Headers:   header.Header,
				Path:      path,
			})
		}
	}
	s.parsed = true
	return nil
}

// Get 會取得指定欄位的檔案資訊。
func (s *Store) Get(key string) (*File, error) {
	if !s.parsed {
		err := s.parse()
		if err != nil {
			return &File{}, err
		}
	}
	v, err := s.GetMulti(key)
	if err != nil {
		return nil, err
	}
	return v[0], nil
}

// GetMulti 會取得指定欄位中的多個檔案。
func (s *Store) GetMulti(key string) ([]*File, error) {
	if !s.parsed {
		err := s.parse()
		if err != nil {
			return []*File{}, err
		}
	}
	v, ok := s.files[key]
	if !ok {
		return []*File{}, ErrNotFound
	}
	return v, nil
}

// Serve 可以接收靜態檔案路徑、*os.File、ioutil.ReadFile 並提供靜態檔案給瀏覽器，就像點擊按鈕會下載檔案那樣。
func (s *Store) Serve(filename string, file interface{}) (err error) {
	switch v := file.(type) {
	case *os.File:
		s.context.Header("Content-Disposition", "attachment; filename="+filename)
		s.context.Header("Content-Type", s.context.Request.Header.Get("Content-Type"))
		_, err := io.Copy(s.context.Writer, v)
		return err
	case []byte:
		s.context.Header("Content-Disposition", "attachment; filename="+filename)
		s.context.Header("Content-Type", s.context.Request.Header.Get("Content-Type"))
		_, err := s.context.Writer.Write(v)
		return err
	case string:
		b, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}
		s.context.Header("Content-Disposition", "attachment; filename="+filename)
		s.context.Header("Content-Type", s.context.Request.Header.Get("Content-Type"))
		_, err = s.context.Writer.Write(b)
		return err
	default:
		return ErrUnknownFileReader
	}
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
	//
	Headers textproto.MIMEHeader
	// keys 為此檔案的鍵值組，可供開發者存放自訂資料。
	keys map[string]interface{}
}

//
func (f *File) Set(key string, value interface{}) {
	f.keys[key] = value
}

//
func (f *File) Get(key string) (interface{}, bool) {
	v, ok := f.keys[key]
	if !ok {
		return nil, false
	}
	return v, true
}

// Remove 會移除這個檔案。
func (f *File) Remove() error {
	return os.Remove(f.Path)
}

// Move 會移動接收到的檔案到指定路徑。
func (f *File) Move(dest string) error {
	return os.Rename(f.Path, dest)
}
