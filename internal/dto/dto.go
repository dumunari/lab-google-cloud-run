package dto

type ViaCEPAPIResponse struct {
	StatusCode int    `json:"status"`
	Location   string `json:"localidade"`
}

type WeatherAPI struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}
