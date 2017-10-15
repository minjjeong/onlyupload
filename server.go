package main

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nfnt/resize"
)

func file_upload(c echo.Context) error {
	// Read form fields
	name := c.FormValue("name")
	email := c.FormValue("email")

	//------------
	// Read files
	//------------

	// Multipart form
	form, err := c.MultipartForm()

	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(file.Filename)
		if err != nil {
			fmt.Print("111")
			return err
		}

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			fmt.Print("222")
			return err
		}

		// create thumbnail
		img, err := jpeg.Decode(dst)
		if err != nil {
			fmt.Print("333\n")
			return err
		}
		defer dst.Close()

		outImg := resize.Thumbnail(100, 100, img, resize.Lanczos3)
		if err != nil {
			fmt.Print("444")
			return err
		}

		// Destination
		dst1, err := os.Create("thumbnail_test.jpg")
		if err != nil {
			fmt.Print("555")
			return err
		}
		defer dst1.Close()

		//if _, err = io.Copy(dst1, outImg); err != nil {
		//	return err
		//}
		jpeg.Encode(dst1, outImg, nil)
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Uploaded successfully %d files with fields name=%s and email=%s.</p>", len(files), name, email))
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthcheck OK!!")
	})
	e.POST("/upload", file_upload)

	e.Logger.Fatal(e.Start(":10304"))
}
