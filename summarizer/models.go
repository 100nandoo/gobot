package summarizer

type SmmryResponse struct {
	SmAPICharacterCount string `json:"sm_api_character_count"`
	SmAPIContentReduced string `json:"sm_api_content_reduced"`
	SmAPITitle          string `json:"sm_api_title"`
	SmAPIContent        string `json:"sm_api_content"`
	SmAPILimitation     string `json:"sm_api_limitation"`
}

type ErrorResponse struct {
	SmAPIError   int    `json:"sm_api_error"`
	SmAPIMessage string `json:"sm_api_message"`
}
