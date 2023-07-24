package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/sync/errgroup"

	"go-mysql/pkg"
)

var g errgroup.Group

func main() {
	err := pkg.MysqlDBLoad()
	if err != nil {
		log.Fatal(err)
		return
	}

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
	emp := &pkg.Employee{}
	if err := c.Bind(emp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := pkg.SaveEmployee(emp)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, emp)
}

func deleteEmployeeHandler(c echo.Context) error {
	empId := c.Param("id")
	err := pkg.DeleteEmployee(empId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func getEmployeeHandler(c echo.Context) error {
	empId := c.Param("id")
	emp, err := pkg.GetEmployee(empId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, emp)
}

func updateEmployeeHandler(c echo.Context) error {
	empId := c.Param("id")
	emp := &pkg.Employee{}
	if err := c.Bind(emp); err != nil {
		return err
	}
	err := pkg.UpdateEmployee(empId, emp)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, emp.Id)
}

func listEmployeeHandler(c echo.Context) error {
	emp, err := pkg.ListEmployee()
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
		defer pkg.MysqlDB.Close()
		log.Fatal(err)
	}
}
