package main

import (
	"net/http"

	"github.com/go-mego/file"
	"github.com/go-mego/mego"
)

func main() {
	e := mego.Default()
	e.POST("/", file.New(), func(c *mego.Context, s *file.Store) {
		f, err := s.Get("file")
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, f)
	})
	e.Run()
}
