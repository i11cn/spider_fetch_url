package main

import (
    "bytes"
    "fmt"
    "github.com/gdamore/mangos"
    "github.com/gdamore/mangos/protocol/pull"
    "github.com/gdamore/mangos/transport/tcp"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
)

func fetch_url(reqs <-chan *http.Request, resps chan<- *http.Response, wg *sync.WaitGroup) {
    defer func() {
        recover()
        wg.Done()
    }()
    client := http.Client{}
    for req := range reqs {
        resp, err := client.Do(req)
        if err == nil {
            resps<- resp
        }
    }
    close(resps)
}

func gen_request(urls <-chan string, reqs chan<- *http.Request, wg *sync.WaitGroup) {
    defer func() {
        recover()
        wg.Done()
    }()
    for url := range urls {
        req, err := http.NewRequest("GET", url, nil)
        if err == nil {
            req.Header.Add("Agent", "Snower Search Enging 0.1")
            reqs<- req
        }
    }
    close(reqs)
}

func do_response(resps <-chan *http.Response, wg *sync.WaitGroup) {
    defer func() {
        recover()
        wg.Done()
    }()
    for resp := range resps {
        buff := bytes.NewBuffer(make([]byte, 0))
        resp.Write(buff)
        fmt.Println(buff)
    }
}

func listen_for_url_list(socket mangos.Socket, urls chan<- string, wg *sync.WaitGroup) {
    defer func() {
        recover()
        wg.Done()
    }()
    for {
        if url, err := socket.Recv(); err == nil {
            urls<- string(url)
        } else if err == mangos.ErrClosed {
            close(urls)
            return
        }
    }
}

func wait_signal(socket mangos.Socket) {
    defer socket.Close()
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGTERM)
    <-c
    signal.Stop(c)
}

func create_socket(ep string) (socket mangos.Socket, err error) {
    if socket, err = pull.NewSocket(); err != nil {
        return
    }
    socket.AddTransport(tcp.NewTransport())
    if err = socket.Listen(ep); err != nil {
        socket.Close()
        return
    }
    return
}

func main() {
    urls := make(chan string)
    reqs := make(chan *http.Request)
    resps := make(chan *http.Response)
    var wg sync.WaitGroup

    socket, err := create_socket("tcp://127.0.0.1:8888")
    if err != nil {
        return
    }
    wg.Add(1)
    go do_response(resps, &wg)
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go fetch_url(reqs, resps, &wg)
    }
    wg.Add(1)
    go gen_request(urls, reqs, &wg)

    wg.Add(1)
    go listen_for_url_list(socket, urls, &wg)

    go wait_signal(socket)

    wg.Wait()
}

