package handlers

// TODO to be determined
//type LaunchHandler interface {
//	HandleStart(rq *gorp.StartLaunchRQ) (*gorp.EntryCreatedRS, error)
//	HandleFinish(rq *gorp.FinishExecutionRQ) (*gorp.FinishExecutionRQ, error)
//}
//
//type HTTPLaunchHandler struct {
//	handler LaunchHandler
//}
//
//func (h *HTTPLaunchHandler) HandleStart(rq *http.Request) (interface{}, error) {
//	var body gorp.StartLaunchRQ
//
//	if err := ReadJSON(rq.Body, &body); err != nil {
//		return nil, err
//	}
//
//	return h.handler.HandleStart(&body)
//}
//func (h *HTTPLaunchHandler) HandleFinish(rq *http.Request) (interface{}, error) {
//	var body gorp.FinishExecutionRQ
//
//	if err := ReadJSON(rq.Body, &body); err != nil {
//		return nil, err
//	}
//
//	return h.handler.HandleFinish(&body)
//}
