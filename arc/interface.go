package arc

import "github.com/gin-gonic/gin"

type Handle interface {
	Build(*gin.RouterGroup)
	Name() string
}
