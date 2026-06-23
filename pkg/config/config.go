package config

type NetworkConfig struct {
	Layers       []int   `json:"layers"`
	Epochs       int     `json:"epochs"`
	RandomRate   float64 `json:"random_rate"`
	SkipRate     float64 `json:"skip_rate"`
	LearningRate float64 `json:"learning_rate"`
}
