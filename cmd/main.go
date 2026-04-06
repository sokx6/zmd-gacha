package main

func main() {
	server := NewServer("config.toml")
	server.Start()
}
