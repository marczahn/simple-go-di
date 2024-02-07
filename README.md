# Simple Go DI

A simple but powerful, functional and type-safe Dependency Injection solution for Golang applications

### Introduction

The idea behind this project is that I see all over the place "microservices" with more and more complex architectures while this was one of the reasons why we introduced *micro*services. With this dependency injection solution, I'd like to bring back the idea of what the architecture of microservices should be: Simple and flat. This is also the reason why this introduction is short - There is not much to explain. Neither does it introduce more complexity, nor does it need any other third-party libraries.

The principle is absolute simple: To your application, each dependency is provided by a service provider which is nothing else than a plain Go function. Since it's based on generics, a dependency can be anything - an object, a slice, an array, or even a scalar.

### Usage

To continue on keeping things simple, let's try it out. For this documentation, I'll use the http example from [examples/http](examples/http). To separate concerns, I recommend to put all the di logic into a separate package.

First, we need an global variable of type `*di.instance`. For this, simple-go-di provides a constructor function. Let's see what it looks like:

```Go
var httpSingleton = di.NewInstance[*http.Client]()
```
This is our actual storage for this dependency. As you can see, with `[*http.Client]`, we define that the value it stores is of type `*http.Client`. That was the first step.

Now, we will add the actual service provider that is later called by your application:

```Go
func HttpClient() *http.Client {}
```
As you can see, nothing else than a function. Now, let's implement the functionality:

```Go 
func HttpClient() *http.Client {
    return httpSingleton.GetOrSet(
        func () *http.Client { 
            return &http.Client{Timeout: 10 * time.Second}
        },
        false,
    )
}
```

All creation logic is encapsulated inside of the callback provided to GetOrSet. The second argument defines if you want to overwrite the value even if it is already set. And that's it. Now, from everywhere in your application, you can receive your client with:

```Go
client := di.HttpClient()
```

To make a small addition, why not create a config that contains the timeout instead of hard coding it at client creation?

```Go
package config

type Config struct{
    timeout time.Duration
}
```

Let's add the dependency logic for our configuration

```Go
var configSingleton = di.NewInstance[*config.Config]()

func Config() *config.Config {
    return configSingleton.GetOrSet(
        func () *config.Config {
            tenv, _ := os.LookupEnv("SIMPLE_GO_DI_TIMOUT")
            timeout, _ := strconv.Atoi(tenv)

            return &config.Config{timout: time.Duration(timeout) * time.Second}
        },
        false,
    )
}
```

And now, altogether:

```Go
var httpSingleton = di.NewInstance[*http.Client]()
var configSingleton = di.NewInstance[*config.Config]()

func Config() *config.Config {
    return configSingleton.GetOrSet(
        func () *config.Config {
            tenv, _ := os.LookupEnv("SIMPLE_GO_DI_TIMOUT")
            timeout, _ := strconv.Atoi(tenv)
            
            return &config.Config{timout: time.Duration(timeout) * time.Second}
        },
        false,
    )
}

func HttpClient() *http.Client {
    return httpSingleton.GetOrSet(
        func () *http.Client {
            return &http.Client{Timeout: Config().timeout}
        },
        false,
    )
}
```

Of course, this example is not really useful, it illustrates how dependency graphs work.