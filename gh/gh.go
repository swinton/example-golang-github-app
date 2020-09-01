package gh

// App encapsulates the fields needed to define a GitHub App
type App struct {
	ID  int
	Key []byte
}

// Installation encapsulates the fields needed to define an installation of a GitHub App
type Installation struct {
	ID int
}
