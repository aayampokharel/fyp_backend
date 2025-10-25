package sse

import (
	"net/http"
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
)

func RegisterRoutes(mux *http.ServeMux, module *Module, ch <-chan entity.Institution) common.SSERouteWrapper {
	var prefix = "/sse"

	// /sse/institution
	var route common.SSERouteWrapper = common.SSERouteWrapper{
		Mux:       mux,
		Prefix:    prefix,
		Route:     "/institution",
		Method:    enum.METHODGET,
		InnerFunc: module.Controller.SendInstitutionsToBeVerified,
		Ch:        ch,
	}
	return route
}
