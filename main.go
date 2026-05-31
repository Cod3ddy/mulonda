package main

import "github.com/cod3ddy/mulonda/cmd"

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
