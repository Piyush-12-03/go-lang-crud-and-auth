package response

type TagsResponse struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Neches []NecheResponse `json:"neches"`
}