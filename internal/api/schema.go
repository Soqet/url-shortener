package api

type CreateRequest struct {
	Url string `json:"url"`
}

type CreateResponse struct {
	ShortUrl string `json:"shortUrl"`
}
