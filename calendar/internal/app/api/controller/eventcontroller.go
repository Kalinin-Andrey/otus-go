package controller

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/go-ozzo/ozzo-routing/v2"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/errorshandler"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/apperror"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
)

type eventController struct {
	Controller
	Service event.IService
	Logger  log.ILogger
}

// RegisterEventHandlers sets up the routing of the HTTP handlers.
//	GET /api/event/
//	GET /api/event/{ID}
//	POST /api/event/
//	PUT /api/event/{ID}
//	DELETE /api/event/{ID}
func RegisterEventHandlers(r *routing.RouteGroup, service event.IService, logger log.ILogger) {
	c := eventController{
		Service:		service,
		Logger:			logger,
	}

	r.Get("/event", c.list)
	r.Get(`/event/<id:\d+>`, c.get)
	r.Post("/event", c.create)
	r.Put(`/event/<id:\d+>`, c.update)
	r.Delete(`/event/<id:\d+>`, c.delete)
}

// get method is for a getting a one enmtity by ID
func (c eventController) get(ctx *routing.Context) error {
	id, err := c.parseUint(ctx, "id")
	if err != nil {
		c.Logger.With(ctx.Request.Context()).Info(errors.Wrapf(err, "Can not parse uint64 from %q", ctx.Param("id")))
		return errorshandler.BadRequest("id mast be a uint")
	}
	entity, err := c.Service.Get(ctx.Request.Context(), uint(id))
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.Write(entity)
}

// list method is for a getting a list of entities
func (c eventController) list(ctx *routing.Context) error {

	items, err := c.Service.List(ctx.Request.Context())
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}
	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.Write(items)
}

func (c eventController) create(ctx *routing.Context) error {
	entity := c.Service.NewEntity()
	if err := ctx.Read(entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	if err := entity.Validate(); err != nil {
		return errorshandler.BadRequest("event invalid: " + err.Error())
	}

	if err := c.Service.Create(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus(entity, http.StatusCreated)
}

func (c eventController) update(ctx *routing.Context) error {
	entity := c.Service.NewEntity()
	if err := ctx.Read(entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	if err := entity.Validate(); err != nil {
		return errorshandler.BadRequest("event invalid: " + err.Error())
	}

	if err := c.Service.Update(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus(entity, http.StatusCreated)
}


func (c eventController) delete(ctx *routing.Context) error {
	id, err := c.parseUint(ctx, "id")
	if err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest("id must be uint")
	}

	if err := c.Service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return errorshandler.Success()
}


