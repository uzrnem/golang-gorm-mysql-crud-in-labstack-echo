package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/sync/errgroup"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlDB *gorm.DB
)

var gr errgroup.Group

type Employee struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    string `json:"age"`
}

func gormMain() {
	dsn := "root:root@tcp(127.0.0.1:3306)/maven_contact_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
	MysqlDB = db

	e := setupEcho()

	setupRoutes(e)

	port := getServerPort()
	mainServer := &http.Server{
		Addr:    port,
		Handler: e,
	}

	startServer(mainServer, port)
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	return e
}

func setupRoutes(e *echo.Echo) {
	e.POST("/employee", createEmployeeHandler)
	e.DELETE("/employee/:id", deleteEmployeeHandler)
	e.GET("/employee/:id", getEmployeeHandler)
	e.PUT("/employee/:id", updateEmployeeHandler)
	e.GET("/employee", listEmployeeHandler)
}

func createEmployeeHandler(c echo.Context) error {
	emp := &Employee{}
	if err := c.Bind(emp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := MysqlDB.Create(emp).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, emp)
}

func deleteEmployeeHandler(c echo.Context) error {
	empId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = MysqlDB.Delete(&Employee{}, empId).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func getEmployeeHandler(c echo.Context) error {
	empId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	emp := &Employee{Id: empId}
	err = MysqlDB.First(&emp).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, emp)
}

func updateEmployeeHandler(c echo.Context) error {
	empId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	emp := &Employee{}
	if err := c.Bind(emp); err != nil {
		return err
	}
	emp.Id = empId
	err = MysqlDB.Save(emp).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, string(emp.Id))
}

func listEmployeeHandler(c echo.Context) error {
	emp := &[]Employee{}
	err := MysqlDB.Find(emp).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, emp)
}

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9055"
	}
	return ":" + port
}

func startServer(mainServer *http.Server, port string) {
	g.Go(func() error {
		return mainServer.ListenAndServe()
	})
	log.Println("Service Running at: " + port)
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
