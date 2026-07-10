package requests

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}