package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/domain/user"
	"github.com/wilzygon/ecommerce/infrastructure/handler/middle"
	storageUser "github.com/wilzygon/ecommerce/infrastructure/postgres/user"
)

// newRoutes función principal de este route, que se instancie un route completo
func NewRouter(e *echo.Echo, dbPool *pgxpool.Pool) {
	h := buildHandler(dbPool) //Creamos el handler

	authMiddleWare := middle.New()

	adminRoutes(e, h, authMiddleWare.IsValid, authMiddleWare.IsAdmin) //registramos las rutas administrativas
	publicRoutes(e, h)
}

func buildHandler(dbPool *pgxpool.Pool) handler {
	storage := storageUser.New(dbPool)
	useCase := user.New(storage)

	return newHandler(useCase)
}

func adminRoutes(e *echo.Echo, h handler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/users", middlewares...)
	g.GET("", h.GetAll)

}

func publicRoutes(e *echo.Echo, h handler) {
	g := e.Group("/api/v1/public/users")
	g.POST("", h.Create)
}
