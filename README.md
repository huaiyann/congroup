# congroup

Congroup provides a wrap of sync.WaitGroup to run and wait functions concurrently.
It is safely on panic. A panic will be wrapped as an error by github.com/pkg/errors. 

### Usage
```go
	cg := congroup.New()
	cg.Add(func() error {
		// do something
	})
	cg.Add(func() error {
		// do something other
	})
	err := cg.Wait()
	if err != nil {
		// handler the error
	}
```

