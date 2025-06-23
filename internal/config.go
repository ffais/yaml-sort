package internal

type Config struct {
	CustomSort  []string `mapstructure:"custom-sort"`
	Reverse     bool     `mapstructure:"reverse"`
	Indent      int      `mapstructure:"indent"`
	SpaceTopKey bool     `mapstructure:"space-top-key"`
}
