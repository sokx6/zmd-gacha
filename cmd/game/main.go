package main

func main() {
	server := NewServer("game.toml")
	server.Start()
}
