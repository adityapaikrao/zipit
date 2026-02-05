package handler

import pb "zipit/gen/url"

type GatewayHandler struct {
	urlSvc pb.URLServiceClient
}

func NewGatewayHandler(urlSvc pb.URLServiceClient) *GatewayHandler {
	return &GatewayHandler{urlSvc: urlSvc}
}
