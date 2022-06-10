package utils

import (

 "fmt"

 "testing"

 "time"

 "golang.org/x/time/rate"

)

func TestLimter(t *testing.T) {

 limiter := rate.NewLimiter(rate.Every(time.Millisecond*31), 2)

 //time.Sleep(time.Second)

 for i := 0; i < 10; i++ {

  var ok bool

  if limiter.Allow() {

   ok = true

  }

  time.Sleep(time.Millisecond * 20)

  fmt.Println(ok, limiter.Burst())

 }

}