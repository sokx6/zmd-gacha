package main

func main() {
	server := NewServer("management.toml")
	server.Start()
}
