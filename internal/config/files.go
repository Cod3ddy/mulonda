package config

const (
	DefaultConfigPath    = "mulonda.yaml"
	DefaultWatchlistPath = "data/watchlist.yaml"
)

type FileConfig struct {
	WatchlistDirectory string `mapstructure:"watchlist_dir"`
}
