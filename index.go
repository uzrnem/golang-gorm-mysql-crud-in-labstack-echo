package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"go-mysql/pkg"
)

func main2() {
	err := pkg.MysqlDBLoad()
	if err != nil {
		log.Fatal(err)
		return
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.POST("/employee", func(c echo.Context) error {
		emp := &pkg.Employee{}
		if err := c.Bind(emp); err != nil {
			c.Echo().Logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err := pkg.SaveEmployee(emp)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusCreated, emp)
	})

	e.DELETE("/employee/:id", func(c echo.Context) error {
		empId := c.Param("id")
		err := pkg.DeleteEmployee(empId)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, "Deleted")
	})

	e.GET("/employee/:id", func(c echo.Context) error {
		empId := c.Param("id")
		emp, err := pkg.GetEmployee(empId)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, emp)
	})

	e.PUT("/employee/:id", func(c echo.Context) error {
		empId := c.Param("id")
		emp := &pkg.Employee{}
		if err := c.Bind(emp); err != nil {
			return err
		}
		err := pkg.UpdateEmployee(empId, emp)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, emp.Id)
	})

	e.GET("/employee", func(c echo.Context) error {
		emp, err := pkg.ListEmployee()
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, emp)
	})
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":9055"
	} else {
		port = ":" + port
	}
	mainServer := &http.Server{
		Addr:    port,
		Handler: e,
	}
	g.Go(func() error {
		return mainServer.ListenAndServe()
	})
	log.Println("Service Running at: " + port)
	if err := g.Wait(); err != nil {
		defer pkg.MysqlDB.Close()
		log.Fatal(err)
	}
}
