package config

import "testing"

func TestGenerateConfig(t *testing.T) {
	GenerateConfig("crypto-config.yaml", cryptoCfg)
}
