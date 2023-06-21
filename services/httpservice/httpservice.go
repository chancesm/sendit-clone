package httpservice

import (
	"log"

	"github.com/chancesm/sendit-clone/services/tunnel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type HttpService struct {
	ts *tunnel.TunnelService
	f  *fiber.App
}

func NewHttpService(t *tunnel.TunnelService) *HttpService {

	engine := html.New("./views", ".html")

	f := fiber.New(fiber.Config{
		Views: engine,
	})
	h := &HttpService{
		ts: t,
	}

	// Register middleware androute handlers and add server to service

	f.Use(logger.New())
	f.Use(favicon.New())
	f.Get("/", h.rootHandler)
	f.Get("/f/:id", h.webHandler)
	f.Get("/f/:id/raw", h.fileHandler)

	// 404 Handler
	f.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
	h.f = f
	return h
}

func (h *HttpService) Run() {
	log.Fatal(h.f.Listen(":1337"))
}

func (h *HttpService) fileHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	tnlchan, found := h.ts.GetTunnelChannel(id)
	if !found {
		return c.Status(404).SendString("Tunnel Not Found")
	}

	donech := make(chan struct{})
	tnlchan <- tunnel.Tunnel{
		Writer:   c.Response().BodyWriter(),
		DoneChan: donech,
	}
	<-donech

	return nil
}

func (h *HttpService) webHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	return c.Render("file", fiber.Map{
		"Title": "Get a File!",
		"ID":    id,
	})
}

func (h *HttpService) rootHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Sendit Home",
	})
}
