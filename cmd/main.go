package main

import (
	"log"
	"os"

	"github.com/wilzygon/ecommerce/infrastructure/handler"
	"github.com/wilzygon/ecommerce/infrastructure/handler/response"
)

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = validateEnvironments()
	if err != nil {
		log.Fatal(err)
	}

	//Crear nuestro echo
	//echo la librería que utilizamos para routing
	e := newHTTP(response.HTTPErrorHandler)

	//Instanciamos la conexión a la base de datos
	dbPool, err := newDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	handler.InitRoutes(e, dbPool)
	//Subimos nuestro Servicio al puerto configurado
	err = e.Start(":" + os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
