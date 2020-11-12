package types

type CreateGetResponse struct {
	Error            bool              `json:"error"`
	ErrorText        string            `json:"errorText"`
	Data             *Items            `json:"data"`
	AdditionalErrors map[string]string `json:"additionalErrors"`
}

type AddDeleteResponse struct {
	Error            bool              `json:"error"`
	ErrorText        string            `json:"errorText"`
	Data             *Id               `json:"data"`
	AdditionalErrors map[string]string `json:"additionalErrors"`
}

type Items struct {
	Items []ListItem `json:"items"`
}

type Id struct {
	Id int `json:"id"`
}

type ListItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type RawItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
