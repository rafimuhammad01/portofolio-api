package jwt

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) VerifyToken(token string) (*Payload, error) {
	return h.service.VerifyToken(token)
}
