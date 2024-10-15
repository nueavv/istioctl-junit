package config

const CLIName = "istioctl-junit"

type Config struct {
	CheckResult bool
	TestName    string
	Format      string
	Output      string
}