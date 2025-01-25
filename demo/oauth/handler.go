package main

type Handler interface {
	GenerateCodeChallenges() (*GenerateCodeChallengesResponse, error)
}

type handler struct {
	oauth Oauth
}

func NewHandler(
	oauth Oauth,
) Handler {
	return &handler{
		oauth: oauth,
	}
}

type GenerateCodeChallengesResponse struct {
	State         string `json:"state" example:"state"`
	CodeVerifier  string `json:"code_verifier" example:"code_verifier"`
	CodeChallenge string `json:"code_challenge" example:"code_challenge"`
}

func (h *handler) GenerateCodeChallenges() (*GenerateCodeChallengesResponse, error) {
	state, err := h.oauth.GenerateState()
	if err != nil {
		return nil, err
	}

	codeVerifier, err := h.oauth.GenerateCodeVerifier()
	if err != nil {
		return nil, err
	}

	codeChallenge, err := h.oauth.GenerateCodeChallenge(codeVerifier)
	if err != nil {
		return nil, err
	}

	res := &GenerateCodeChallengesResponse{
		State:         state,
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
	}

	return res, nil
}
