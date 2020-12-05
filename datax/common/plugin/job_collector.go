package plugin

type JobCollector interface {
	MessageMap() map[string][]string
	MessageByKey(key string) []string
}
