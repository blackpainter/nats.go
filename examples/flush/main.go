package main

import (
    "fmt"
    "net"
    "net/http"
    "time"
    )

func main(){
    //var urls = "nats://127.0.0.1:4221,nats://127.0.0.1:4222,nats://127.0.0.1:4223"

    buf:=make(chan int)
    flg := make(chan int)
    // go producer(buf)
    go consumer(buf, flg)
    <-flg
}

func producer(c chan int){
    defer close(c)
    for i:= 0; i <= 1000000; i++{
        c <- i
    }
}

func postToServer(name string, client http.Client) {
    resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/nats?name=%s&url=https://www.douban.com/", name))
    if err != nil {
        fmt.Println("Error ", resp.StatusCode, " : ", err)
        return
    }
    defer resp.Body.Close()
}

func consumer(c, f chan int){
    connTimeout := 2*time.Millisecond
    readTimeout := 3*time.Millisecond

    client := http.Client{
        Transport: &http.Transport{
            Dial: func(netw, addr string) (net.Conn, error) {
                c, err := net.DialTimeout("tcp4", addr, connTimeout)
                if err != nil {
                    return nil, err
                }
                c.SetDeadline(time.Now().Add(readTimeout))
                return c, nil
            },
        },
    }

    // for{
    //    if v, ok := <-c; ok{
    //        postToServer(fmt.Sprintf("%d: %d", f, v), client)
    //    }else{
    //        break
    //    }
    // }
    // f<-1
    
    i := 0
    for {
	postToServer(fmt.Sprintf("%d", i), client)
	i ++
    }
}

