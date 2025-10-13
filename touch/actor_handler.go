package touch

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	log "github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-go/touch/actor"
	"github.com/dirty-bro-tech/peers-touch-go/touch/auth"
	"github.com/dirty-bro-tech/peers-touch-go/touch/model"
)

// ActorHandlerInfo represents a single handler's information
type ActorHandlerInfo struct {
	RouterURL RouterPath
	Handler   func(context.Context, *app.RequestContext)
	Method    server.Method
	Wrappers  []server.Wrapper
}

// GetActorHandlers returns all actor handler configurations
func GetActorHandlers() []ActorHandlerInfo {
	return []ActorHandlerInfo{
		{
			RouterURL: RouterURLActorSignUP,
			Handler:   ActorSignup,
			Method:    server.POST,
			Wrappers:  []server.Wrapper{CommonAccessControlWrapper("Actor")},
		},
		{
			RouterURL: RouterURLActorLogin,
			Handler:   ActorLogin,
			Method:    server.POST,
			Wrappers:  []server.Wrapper{CommonAccessControlWrapper("Actor")},
		},
		{
			RouterURL: RouterURLActorProfile,
			Handler:   GetActorProfile,
			Method:    server.GET,
			Wrappers:  []server.Wrapper{CommonAccessControlWrapper("Actor")},
		},
		{
			RouterURL: RouterURLActorProfile,
			Handler:   UpdateActorProfile,
			Method:    server.PUT,
			Wrappers:  []server.Wrapper{CommonAccessControlWrapper("Actor")},
		},
	}
}

// Handler implementations

func ActorSignup(c context.Context, ctx *app.RequestContext) {
	var params model.ActorSignParams
	if err := ctx.Bind(&params); err != nil {
		log.Warnf(c, "Signup bound params failed: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := params.Check(); err != nil {
		log.Warnf(c, "Signup checked params failed: %v", err)
		FailedResponse(ctx, err)
		return
	}

	err := actor.SignUp(c, &params)
	if err != nil {
		log.Warnf(c, "Signup failed: %v", err)
		FailedResponse(ctx, err)
		return
	}

	SuccessResponse(ctx, "Actor signup successful", nil)
}

func ActorLogin(c context.Context, ctx *app.RequestContext) {
	var params model.ActorLoginParams
	if err := ctx.Bind(&params); err != nil {
		log.Warnf(c, "Login bound params failed: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := params.Check(); err != nil {
		log.Warnf(c, "Login checked params failed: %v", err)
		FailedResponse(ctx, err)
		return
	}

	// Prepare credentials
	credentials := &auth.Credentials{
		Email:    params.Email,
		Password: params.Password,
	}

	// Get client IP and user agent
	clientIP := ctx.ClientIP()
	userAgent := string(ctx.GetHeader("User-Agent"))

	// Use auth service to handle login with session
	result, err := auth.LoginWithSession(c, credentials, clientIP, userAgent)
	if err != nil {
		log.Warnf(c, "Login failed: %v", err)
		FailedResponse(ctx, err)
		return
	}

	// Set session cookie
	ctx.SetCookie("session_id", result.SessionID, int(24*time.Hour.Seconds()), "/", "", protocol.CookieSameSiteDisabled, false, true)

	// Return success response
	SuccessResponse(ctx, "Login successful", result)
}

func GetActorProfile(c context.Context, ctx *app.RequestContext) {
	// Get current actor from context
	currentActor := auth.GetCurrentActor(ctx)
	if currentActor == nil {
		FailedResponse(ctx, auth.ErrNotLoggedIn)
		return
	}

	// Get actor profile
	profile, err := actor.GetProfile(c, currentActor.ID)
	if err != nil {
		log.Warnf(c, "Get actor profile failed: %v", err)
		FailedResponse(ctx, err)
		return
	}

	SuccessResponse(ctx, "Actor profile retrieved", profile)
}

func UpdateActorProfile(c context.Context, ctx *app.RequestContext) {
	var params model.ProfileUpdateParams
	if err := ctx.Bind(&params); err != nil {
		log.Warnf(c, "Update profile bound params failed: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Get current actor from context
	currentActor := auth.GetCurrentActor(ctx)
	if currentActor == nil {
		FailedResponse(ctx, auth.ErrNotLoggedIn)
		return
	}

	if err := actor.UpdateProfile(c, currentActor.ID, &params); err != nil {
		log.Warnf(c, "Update profile failed: %v", err)
		FailedResponse(ctx, err)
		return
	}

	SuccessResponse(ctx, "Profile updated successfully", nil)
}
