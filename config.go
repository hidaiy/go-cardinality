package dbindex

type Config struct {
	User          string
	Password      string
	Host          string
	Port          int
	Dialect       string
	Database      string
	Threshold     int
	IgnoreTables  []string
	IgnoreColumns map[string]map[string]int
}
