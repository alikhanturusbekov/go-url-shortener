package model

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URLPair struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}
