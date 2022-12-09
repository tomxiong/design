package http

import (
	_ "design/docs"
	design "design/internal"
	"design/internal/conf"
	"design/internal/model"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server is http server.
type Server struct {
	engine *gin.Engine
	design *design.Design
}

// New a http server.
func New(c *conf.HTTPServer, d *design.Design) *Server {
	engine := gin.New()
	engine.Use(loggerHandler, recoverHandler)
	go func() {
		if err := engine.Run(c.Addr); err != nil {
			panic(err)
		}
	}()
	s := &Server{
		engine: engine,
		design: d,
	}
	s.initRouter()
	return s
}

func (s *Server) initRouter() {
	//s.engine.Use(cors())
	group := s.engine.Group("/v1")

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	group.POST("/member/register", s.register) // registration
	group.POST("/member/setemail", s.setEmail) // update e-mail for registration
	group.GET("/member/list", s.list)          // list
}

func (s *Server) HandleResult(c *gin.Context, err error, data interface{}) {
	if err != nil {
		failResult(c, ForbiddenErr, err.Error())
		return
	}

	successResult(c, data)
	return
}

// Close close the server.
func (s *Server) Close() {

}

// @Summary member register API
// @Description register action for new member(include normal member, validator member and manager member).
// @Tags member
// @Produce  json
// @Param input body model.Member true "member"
// @Success 200 {object} resp
// @Failure 400 {object} resp
// @Router /v1/member/register [post]
func (s *Server) register(c *gin.Context) {
	input := model.Member{}
	if err := c.BindJSON(&input); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	done, err := s.design.Register(input)
	handleResult(c, err, done)
}

// @Summary member set e-mail API
// @Description update e-mail after registration
// @Tags member
// @Produce  json
// @Param input body model.Member true "member"
// @Success 200 {object} resp
// @Failure 400 {object} resp
// @Router /v1/member/setemail [post]
func (s *Server) setEmail(c *gin.Context) {
	input := model.Member{}
	if err := c.BindJSON(&input); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	done, err := s.design.SetEmail(input)
	handleResult(c, err, done)
}

// @Summary list member by role and status
// @Description list member by role (include normal member, validator member and manager member) and status: new or other
// @Tags member
// @Produce  json
// @Param role query string false "role"
// @Param status query string false "status"
// @Success 200 {object} []model.Member
// @Failure 400 {object} resp
// @Router /v1/member/list [get]
func (s *Server) list(c *gin.Context) {
	role := c.Query("role")
	status := c.Query("status")
	data, err := s.design.ListMember(role, status)
	handleResult(c, err, data)
}
