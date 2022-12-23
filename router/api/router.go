package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine // gin router
)

type Controllers struct {
	R *gin.Engine
}

func init() {
	// Default With the Logger and Recovery middleware already attached gin.New() would be without middleware
	// Router = gin.Default()

	// Create router and add middleware
	Router = gin.New()
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	Router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	Router.Use(gin.Recovery())
	//   // Per route middleware, you can add as many as you desire.
	//   r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := Router.Group("/auth")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(AuthRequired) // May be AuthRequired() needed
	// {
	// 	authorized.POST("/login", loginEndpoint)
	// 	authorized.POST("/submit", submitEndpoint)
	// 	authorized.POST("/read", readEndpoint)

	// 	// nested group
	// 	testing := authorized.Group("testing")
	// 	// visit 0.0.0.0:8080/testing/analytics
	// 	testing.GET("/analytics", analyticsEndpoint)
	// }

	// TODO: add TLS
	// TODO: HTTP2 add
	// custom router setup
	s := &http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        Router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
	// Simple group: v2
	// v2 := router.Group("/auth")
	// {
	// 	v2.POST("/login", loginEndpoint)
	// 	v2.POST("/submit", submitEndpoint)
	// 	v2.POST("/read", readEndpoint)
	// }
	//   // Creates a router without any middleware by default
	//   r := gin.New()

	//   // Global middleware
	//   // Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	//   // By default gin.DefaultWriter = os.Stdout
	//   r.Use(gin.Logger())

	//   // Recovery middleware recovers from any panics and writes a 500 if there was one.
	//   r.Use(gin.Recovery())

	//   // Per route middleware, you can add as many as you desire.
	//   r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	//   // Authorization group
	//   // authorized := r.Group("/", AuthRequired())
	//   // exactly the same as:
	//   authorized := r.Group("/")
	//   // per group middleware! in this case we use the custom created
	//   // AuthRequired() middleware just in the "authorized" group.
	//   authorized.Use(AuthRequired())
	//   {
	//     authorized.POST("/login", loginEndpoint)
	//     authorized.POST("/submit", submitEndpoint)
	//     authorized.POST("/read", readEndpoint)

	//	  // nested group
	//	  testing := authorized.Group("testing")
	//	  // visit 0.0.0.0:8080/testing/analytics
	//	  testing.GET("/analytics", analyticsEndpoint)
	//	}

	// custom recovery from 500 and other errors
	//   // Creates a router without any middleware by default
	//   r := gin.New()

	//   // Global middleware
	//   // Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	//   // By default gin.DefaultWriter = os.Stdout
	//   r.Use(gin.Logger())

	//   // Recovery middleware recovers from any panics and writes a 500 if there was one.
	//   r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
	//     if err, ok := recovered.(string); ok {
	//       c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	//     }
	//     c.AbortWithStatus(http.StatusInternalServerError)
	//   }))

	//   r.GET("/panic", func(c *gin.Context) {
	//     // panic with a string -- the custom middleware could save this to a database or report it to the user
	//     panic("foo")
	//   })

	//	r.GET("/", func(c *gin.Context) {
	//	  c.String(http.StatusOK, "ohai")
	//	})
}

// Middleware custom auth
func AuthRequired(c *gin.Context) {
	c.String(http.StatusUnauthorized, "Authorization required...")
}
