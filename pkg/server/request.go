package server

type PredictionRequest struct {
	Inputs []float64 `json:"inputs"`
}
