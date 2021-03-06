package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-web/src/controllers"
	cf "github.com/red-gold/telar-web/src/controllers/notifications/config"
	"github.com/red-gold/telar-web/src/controllers/notifications/handlers"
)

func init() {

	cf.InitConfig()
}

// Cache state
var server *coreServer.ServerRouter
var db interface{}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Start
	if db == nil {
		var startErr error
		db, startErr = controllers.Start(ctx)
		if startErr != nil {
			fmt.Printf("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
	}

	// Server Routing
	if server == nil {
		server = coreServer.NewServerRouter()
		server.POST("/", handlers.CreateNotificationHandle(db), coreServer.RouteProtectionHMAC)
		server.PUT("/", handlers.UpdateNotificationHandle(db), coreServer.RouteProtectionHMAC)
		server.PUT("/seen/:notificationId", handlers.SeenNotificationHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/id/:notificationId", handlers.DeleteNotificationHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/my", handlers.DeleteNotificationByUserIdHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/", handlers.GetNotificationsByUserIdHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/:notificationId", handlers.GetNotificationHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
