package structs

type WordCounterResponse struct {
	Status         int         `json:"status"`
	WordCollection []WordCount `json:"words"`
}
