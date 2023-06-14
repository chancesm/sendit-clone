package httpservice

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/chancesm/sendit-clone/services/tunnel"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type HttpService struct {
	ts *tunnel.TunnelService
	e  *echo.Echo
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewHttpService(t *tunnel.TunnelService) *HttpService {
	e := echo.New()

	h := &HttpService{
		ts: t,
	}
	tmp := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	// Register middleware androute handlers and add server to service
	e.Use(middleware.Logger())
	e.Renderer = tmp
	e.GET("/", h.rootHandler)
	e.GET("/:id", h.webHandler)
	e.GET("/:id/raw", h.fileHandler)
	e.GET("/favicon.ico", echo.NotFoundHandler)
	h.e = e

	return h
}

func (h *HttpService) Run() {
	h.e.Logger.Fatal(h.e.Start(":1337"))
}

func (h *HttpService) fileHandler(c echo.Context) error {

	idstr := c.Param("id")
	id, _ := strconv.Atoi(idstr)

	tnlchan, found := h.ts.GetTunnelChannel(id)
	if !found {
		return c.String(http.StatusNotFound, "Tunnel Not Found")
	}

	donech := make(chan struct{})
	tnlchan <- tunnel.Tunnel{
		Writer:   c.Response().Writer,
		DoneChan: donech,
	}
	<-donech

	return nil
}

func (h *HttpService) webHandler(c echo.Context) error {
	idstr := c.Param("id")
	id, _ := strconv.Atoi(idstr)
	pagedata := struct{ ID int }{ID: id}

	return c.Render(http.StatusOK, "file.html", pagedata)
}
func (h *HttpService) rootHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
