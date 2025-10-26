package sse

import (
	"fmt"
	"net/http"
	"project/internals/usecase"
	"project/package/enum"
	"project/package/utils/common"
)

type Controller struct {
	sqlUseCase *usecase.SqlUseCase
	sseUseCase *usecase.SSEUseCase
}

func NewController(sqlUseCase *usecase.SqlUseCase, sseUseCase *usecase.SSEUseCase) *Controller {
	return &Controller{sqlUseCase: sqlUseCase, sseUseCase: sseUseCase}
}

func (c *Controller) SendInstitutionsToBeVerified(w http.ResponseWriter, r *http.Request) {

	token := common.GenerateUUID(25)
	newInstitutionCh := c.sseUseCase.SSEManager.AddClient(token)
	defer c.sseUseCase.SSEManager.RemoveClient(token)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Connection-Token", token)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "event: connection_established\ndata: {\"token\":\"%s\"}\n\n", token)
	flusher.Flush()

	ctx := r.Context()
	c.sqlUseCase.Logger.Infoln("[send_institutions_to_be_verified] Info: sendInstitutionsToBeVerified", "started")
	for {
		select {
		case newInstitution := <-newInstitutionCh:
			HandleSSEResponse(newInstitution, enum.SSESINGLEFORM, w)
		case <-ctx.Done():
			c.sqlUseCase.Logger.Infoln("[send_institutions_to_be_verified] Info: sendInstitutionsToBeVerified::CLIENt disconected ! ", fmt.Sprint(ctx.Err()))
			return
		}
	}
}
