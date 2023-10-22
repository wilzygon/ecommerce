package main

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newHTTP(errorHandler echo.HTTPErrorHandler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())  //middleware para loggear todas las peticiones
	e.Use(middleware.Recover()) //middleware para recuperar errores

	corsConfig := middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowMethods: strings.Split(os.Getenv("ALLOWED_METHODS"), ","),
	}

	e.Use(middleware.CORSWithConfig(corsConfig)) //middleware para aceptar los Cors de acuerdo a lo que
	//configure el Administrador en .env.example

	e.HTTPErrorHandler = errorHandler //errorHandler se asigna a echo para configurar nuestro standar de
	//errores

	return e
}
