package main

// main
func main() {
	a := App{}
	a.Initialize("data.json", "sitemap.json")
	a.Run(8080)
}
