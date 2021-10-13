package web

type Link struct {
	LinkID string `uri:"link_id" binding:"required"`
}

type AddLinkParams struct {
	ShortLinkLength uint8  `json:"short_link_length"`
	LinkContent     string `json:"link_content"`
}

type DeleteLinkParams struct {
	LinkID string `json:"link_id"`
}

type CheckShortLink struct {
	LinkID  string `json:"link_id"`
	Content string `json:"content"`
}
