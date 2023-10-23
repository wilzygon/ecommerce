package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/domain/user"
	storageUser "github.com/wilzygon/ecommerce/infrastructure/postgres/user"
)

// newRoutes función principal de este route, que se instancie un route completo
func newRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	h := buildHandler(dbPool) //Creamos el handler

	adminRoutes(e, h) //registramos las rutas administrativas
	publicRoutes(e, h)
}

func buildHandler(dbPool *pgxpool.Pool) handler {
	storage := storageUser.New(dbPool)
	useCase := user.New(storage)

	return newHandler(useCase)
}

func adminRoutes(e *echo.Echo, h handler) {
	g := e.Group("/api/v1/admin/users")
	g.GET("", h.GetAll)

}

func publicRoutes(e *echo.Echo, h handler) {
	g := e.Group("/api/v1/public/users")
	g.POST("", h.Create)
}
