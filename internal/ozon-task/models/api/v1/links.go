package models

type (
	ApiNewLinkRequest struct {
		OriginalLink string `json:"original_link"`
	}

	ApiNewLinkResponse struct {
		ShortLink string `json:"short_link"`
	}

	ApiGetOriginalLinkResponse struct {
		OriginalLink string `json:"short_link"`
	}
)
