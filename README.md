# go-physics-units
go-physics-units is a Golang library for working with physical units such as distance and time.
This library provides functionalities to create, convert, and manipulate units easily.

## <span style="color:red">Warning!</span>
This library is under active development now and still unstable

## Instalation

```shell
go get -u github.com/asaphin/go-physics-units
```

## Usage

```go
package main

import (
	"fmt"
	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
)

func main() {
    // Creating a distance of 5 kilometers
    dist, err := units.NewDistance(5, distance.Kilometer)
    if err != nil {
        panic(err)
    }

    fmt.Println(dist) // 5 km

    // Parsing a distance from a string "10 km"
    ms, err := units.ParseString("10 km")
    if err != nil {
        panic(err)
    }

    // Converting distance to a specific unit (nautical miles in this case)
    dist, err = units.ToDistance(ms)
    if err != nil {
        panic(err)
    }

    convDist, err := dist.ConvertTo("nmi")
    if err != nil {
        panic(err)
    }

    fmt.Println(convDist) // 5.399568034557235 nmi

    // Accessing the value and unit of the distance
    fmt.Println(dist.Value()) // 10
    fmt.Println(dist.Unit()) // km

    // Creating a time of 600 seconds
    tm, err := units.NewTime(600, "s")
    if err != nil {
        panic(err)
    }

    // Calculating velocity by dividing distance by time
    vel := dist.DivideByTime(tm)

    fmt.Println(vel) // 16.666666666666668 m/s
}
```

You can also create a factory for your own measure type:

```go
package main

import (
	"fmt"
	units "github.com/asaphin/go-physics-units"
)

func main() {
	type Currency interface {
		units.Measurement
	}

	currencyFactors := conversion.Factors{
		"EUR": 1.1,
		"USD": 1,
		"UAH": 0.027,
	}

	mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", currencyFactors)
	if err != nil {
		panic(err)
	}

	cur, err := mc.New(1, "USD")
	if err != nil {
		panic(err)
	}

	cur2, err := cur.ConvertToMeasurement("UAH")
	if err != nil {
		panic(err)
	}

	fmt.Println(cur2) // 37.03703703703704 UAH
}
```

## Documentation

## License

MIT License