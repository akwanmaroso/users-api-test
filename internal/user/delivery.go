package user

import "github.com/labstack/echo/v4"

type Handlers interface {
	Create() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Detail() echo.HandlerFunc
	List() echo.HandlerFunc
}
