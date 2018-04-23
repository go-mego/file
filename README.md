# File [![GoDoc](https://godoc.org/github.com/go-mego/file?status.svg)](https://godoc.org/github.com/go-mego/file)

File 檔案套件可以協助你在路由中處理客戶端所上傳的單、多檔案，同時也能夠將檔案提供給客戶端進行直接性地下載。

# 索引

* [安裝方式](#安裝方式)
* [使用方式](#使用方式)
	* [取得檔案](#取得檔案)
    * [多個檔案](#多個檔案)
    * [提供檔案](#提供檔案)

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/go-mego/file
```

# 使用方式

將 `file.New` 傳入 Mego 引擎中的 `Use` 來將檔案中介軟體作為全域中介軟體，即能在所有路由中使用與檔案有關的功能。

```go
package main

import (
	"github.com/go-mego/file"
	"github.com/go-mego/mego"
)

func main() {
	m := mego.New()
	// 將檔案中介軟體作為全域中介軟體就可以在所有路由中使用。
	m.Use(file.New())
	m.Run()
}
```

檔案中介軟體也能夠僅用於單個路由。

```go
func main() {
	m := mego.New()
	// 檔案中介軟體可以獨立用於某路由。
	m.POST("/", file.New(), func(f *file.Store) {
		// ...
	})
	m.Run()
}
```

## 取得檔案

透過 `Get` 來取得指定欄位的檔案資訊。

```go
func main() {
	m := mego.New()
	m.POST("/", file.New(), func(f *file.Store) {
		// 透過 `Get` 來取得指定欄位的檔案，如果該檔案為空則會回傳錯誤。
		if file, err := f.Get("Photo"); err == nil {
			// 輸出檔案的名稱。
			fmt.Println(file.Name)
		}
	})
	m.Run()
}
```

如果一個檔案欄位是必要的，那麼就可以透過 `MustGet` 來在沒有檔案時自動產生 `panic` 阻斷請求。這種情況需要配合 Mego 的 Recovery 中介軟體以避免單次 `panic` 造成整個伺服器終止。

```go
func main() {
	m := mego.New()
	m.POST("/", file.New(), func(f *file.Store) {
		// 透過 `MustGet` 可以自動中斷必要檔案但為空的請求。
		file := f.MustGet("Photo")
		fmt.Println(file.Name)
	})
	m.Run()
}
```

## 多個檔案

透過 `GetMulti` 可以取得指定欄位中的多個檔案。

```go
func main() {
	m := mego.New()
	m.POST("/", file.New(), func(f *file.Store) {
		// `GetMulti` 會回傳一個檔案切片，如果該欄位為空則會回傳錯誤。
		if files, err := f.GetMulti("Photos"); err == nil {
			for _, file := range files {
				fmt.Println(file.Name)
			}
		}
	})
	m.Run()
}
```

如果該檔案欄位是必要的，那麼就可以透過 `MustGetMulti` 來在檔案欄位為空時自動以 `panic` 中斷請求。

```go
func main() {
	m := mego.New()
	m.POST("/", file.New(), func(f *file.Store) {
		// 透過 `MustGetMulti` 來確保必要的多檔案欄位是有檔案的，若無則中斷請求。
		files := f.MustGetMulti("Photos")
		for _, file := range files {
			fmt.Println(file.Name)
		}
	})
	m.Run()
}
```

## 提供檔案

透過 `Serve` 可以傳入靜態檔案路徑、`*os.File`、`ioutil.ReadFile` 提供檔案給瀏覽器。就像點擊按鈕會下載檔案那樣。

```go
func main() {
	m := mego.New()
	m.GET("/", file.New(), func(f *file.Store) {
		// 透過 `os.Open` 來開啟一個檔案並取得 `*os.File`。
		dat, err := os.Open("/tmp/dat")
		if err != nil {
			panic(err)
		}
		// 透過 `Serve` 來將 `*os.File` 的檔案提供給瀏覽器下載。
		err := f.Serve(dat)
		if err != nil {
			panic(err)
		}
		// 同時也可以使用下列用法：
		// f.Serve(ioutil.ReadFile("./example.png"))
		// f.Serve("./example.png")
	})
	m.Run()
}
```