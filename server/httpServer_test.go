package main

import (
	"github.com/gin-gonic/gin"
	"lms/middleware"
	"log"
	"strconv"
)

func dummyRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//router.LoadHTMLFiles(fmt.Sprintf("%v/index.html", path))
	//router.Use(static.Serve("/static", static.LocalFile(staticPath, true)))
	log.Println("running dummyRouter for testing purpose")

	g1 := router.Group("/")
	g1.Use(middleware.UserAuth())
	{
		g1.POST("/getUserBooks", getUserBooksHandler)
		g1.POST("/borrowBook", borrowBookHandler)
		g1.POST("/returnBook", returnBookHandler)
	}

	g2 := router.Group("/")
	g2.Use(middleware.AdminAuth())
	{
		g2.POST("/updateBookStatus", updateBookStatusHandler)
		g2.POST("/deleteBook", deleteBookHandler)
		g2.POST("/addBook", addBookHandler)
	}

	router.POST("/login", loginHandler)
	router.POST("/admin", adminLoginHandler)
	router.POST("/register", registerHandler)
	router.GET("/getCount", getCountHandler)
	router.GET("/getBooks", getBooksHandler)
	router.POST("/getBooks", getBooksHandler)

	//router.StaticFile("/favicon.ico", fmt.Sprintf("%v/favicon.ico", staticPath))

	err := router.Run(":" + strconv.Itoa(80))
	if err != nil {
		log.Println(err)
		return nil
	} else {
		return router
	}
}

//func Test_addBookHandler(t *testing.T) {
//	router := startService()
//	w := httptest.NewRecorder()
//	testContext, _ := gin.CreateTestContext(w)
//	addBookHandler(testContext)
//}
