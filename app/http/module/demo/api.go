package demo

import (
	demoService "github.com/PengLei-Adam/lade/app/provider/demo"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/PengLei-Adam/lade/framework/gin"
)

type DemoApi struct {
	service *Service
}

func Register(r *gin.Engine) error {
	api := NewDemoApi()
	r.Bind(&demoService.DemoProvider{})

	r.GET("/demo/demo", api.Demo)
	r.GET("/demo/demo2", api.Demo2)
	r.GET("/demo/demo3", api.Demo3)
	r.POST("/demo/demo_post", api.DemoPost)
	return nil

}

func NewDemoApi() *DemoApi {
	service := NewService()
	return &DemoApi{service: service}
}

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo [get]
func (api *DemoApi) Demo(c *gin.Context) {
	appService := c.MustMake(contract.AppKey).(contract.App)
	baseFolder := appService.BaseFolder()
	//users := api.service.GetUsers()
	//usersDTO := UserModelsToUserDTOs(users)
	//c.JSON(200, usersDTO)
	c.JSON(200, baseFolder)
}

// Demo godoc
// @Summary 获取所有学生
// @Description 获取所有学生
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo2 [get]
func (api *DemoApi) Demo2(c *gin.Context) {
	demoProvider := c.MustMake(demoService.DemoKey).(demoService.IService)
	students := demoProvider.GetAllStudent()
	usersDTO := StudentsToUserDTOs(students)
	c.JSON(200, usersDTO)
}

func (api *DemoApi) DemoPost(c *gin.Context) {
	type Foo struct {
		Name string
	}
	foo := &Foo{}
	err := c.BindJSON(&foo)
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.JSON(200, nil)
}

// Demo godoc
// @Summary 获取密码
// @Description 获取数据库密码
// @Produce  json
// @Tags demo
// @Success 200 string
// @Router /demo/demo3 [get]
func (api *DemoApi) Demo3(c *gin.Context) {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	password := configService.GetString("database.mysql.password")
	c.JSON(200, password)
}
