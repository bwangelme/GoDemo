package main

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type MyStruct struct {
	foo string
}

func MyFunction(arg string) {

}

/*使用 go-cache 的示例程序*/
func main() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Set("foo", "bar", cache.DefaultExpiration)

	// Set the value of the key "baz" to 42, with no expiration time
	// (the item won't be removed until it is re-set, or removed using
	// c.Delete("baz")
	c.Set("baz", 42, cache.NoExpiration)

	// Get the string associated with the key "foo" from the cache
	foo, found := c.Get("foo")
	if found {
		fmt.Println(foo)
	}

	// Since Go is statically typed, and cache values can be anything, type
	// assertion is needed when values are being passed to functions that don't
	// take arbitrary types, (i.e. interface{}). The simplest way to do this for
	// values which will only be used once--e.g. for passing to another
	// function--is:
	fooI, found := c.Get("foo")
	if found {
		MyFunction(fooI.(string))
	}

	// This gets tedious if the value is used several times in the same function.
	// You might do either of the following instead:
	if x, found := c.Get("foo"); found {
		foo := x.(string)
		fmt.Println(foo)
	}
	// or
	if x, found := c.Get("foo"); found {
		foo = x.(string)
		fmt.Println(foo)
	}
	// ...
	// foo can then be passed around freely as a string

	// Want performance? Store pointers!
	var s = MyStruct{}
	c.Set("foo", &s, cache.DefaultExpiration)
	if x, found := c.Get("foo"); found {
		fooS := x.(*MyStruct)
		fmt.Println(fooS)
	}
}
