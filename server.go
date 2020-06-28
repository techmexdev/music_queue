package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", handleHealth)
	e.GET("/tnd", handleTND)
	e.GET("/sub-mu-core", handleSubMuCore)
	e.GET("/health", handleHealth)
	return e
}

func handleHealth(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func handleSubMuCore(c echo.Context) error {
	return c.JSON(http.StatusOK, subMuCoreAlbums)
}

func handleTND(c echo.Context) error {
	if rating, err := strconv.Atoi(c.QueryParam("rating")); err == nil {
		aa := ratingAlbums[rating]
		return c.JSON(http.StatusOK, aa)
	}

	var rr []int
	for r, _ := range ratingAlbums {
		rr = append(rr, r)
	}
	sort.Ints(rr)

	var aa []Album
	for i := len(rr) - 1; i >= 0; i-- {
		albs := ratingAlbums[rr[i]]
		aa = append(aa, albs...)
	}

	return c.JSON(http.StatusOK, aa)
}
