//go:build ignore

package main

import (
	"encoding/base64"
	"os"
)

func main() {
	ico := "AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAABILAAASCwAAAAAAAAAAAAD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/AP///wD///8A////AP///wD///8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8A////AP///wD///8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/AP///wD///8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wD///8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wD///8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/AP///wD///8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8A////AP///wD///8A////AP///wD///8AJoT/ACaE/wAmhP8AJoT/ACaE/wAmhP8AJoT/AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A//8AAP//AAD8PwAA+B8AAPAPAADgBwAAwAMAAIABAACAAQAAwAMAAOAHAADwDwAA+B8AAPw/AAD//wAA//8AAA=="
	data, _ := base64.StdEncoding.DecodeString(ico)
	os.MkdirAll("build/windows", 0755)
	os.WriteFile("build/windows/icon.ico", data, 0644)
}