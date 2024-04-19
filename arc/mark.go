package arc

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type Mark struct {
	*gin.Engine
	currentGroup *gin.RouterGroup
}

/*
Launch 启动 mark
*/
func (m *Mark) Launch(address ...string) error {
	addr := resolveAddress(address)
	return endless.ListenAndServe(addr, m)
}

func (m *Mark) New() *Mark {
	m.Engine = gin.New()
	return m
}

func (m *Mark) Default() *Mark {
	m.Engine = gin.Default()
	return m
}

func (m *Mark) SwitchGroup(group string, router *gin.RouterGroup, handles ...gin.HandlerFunc) *Mark {
	if router == nil {
		m.currentGroup = m.Engine.Group(group, handles...)
	} else {
		m.currentGroup = router.Group(group, handles...)
	}

	return m
}

func (m *Mark) GUse(middlewares ...gin.HandlerFunc) *Mark {
	if m.currentGroup == nil {
		m.Engine.Use(middlewares...)
	} else {
		m.currentGroup.Use(middlewares...)
	}
	return m
}

func (m *Mark) Mount(args ...Handle) {
	for _, handle := range args {
		handle.Build(m.currentGroup)
	}
}
