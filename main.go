package main

// main
func main() {
	a := App{}
	a.Initialize("data.json")
	a.Run(8080)
}
