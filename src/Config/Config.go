package Config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	HideExtensions bool `json:"HideExtensions"`
}

func Load(Path string) AppConfig {
	Cfg := AppConfig{HideExtensions: true}
	Bytes, Err := os.ReadFile(Path)
	if Err == nil {
		json.Unmarshal(Bytes, &Cfg)
	}
	return Cfg
}