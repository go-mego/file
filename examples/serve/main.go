package main

import (
	"github.com/go-mego/file"
	"github.com/go-mego/mego"
)

func main() {
	e := mego.Default()
	e.GET("/", file.New(), func(c *mego.Context, s *file.Store) {
		s.Serve("myFile.png", "./example.png")
		return
	})
	e.Run()
}
