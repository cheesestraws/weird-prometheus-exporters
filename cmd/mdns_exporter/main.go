package main

func main() {
	m := newmdnsWatcher()

	m.watchServices()
}
