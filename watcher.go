package confd

// Watcher watch config change
type Watcher interface {
	OnChange()
}
