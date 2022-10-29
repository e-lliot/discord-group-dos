package main

import (
    "github.com/fatih/color"
    "time"
    "fmt"
    "net/http"
    "os"
    "bufio"
    "os/exec"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "github.com/tidwall/gjson"
)


var (
    variable string
    endpoint string
    gid string
    TOKEN_AUTH string
    logo = `
M""""""'YMM  oo                                           8P 
M  mmmm. 'M                                               88 
M  MMMMM  M  dP .d8888b. .d8888b. .d8888b. 88d888b. .d888b88 
M  MMMMM  M  88 Y8ooooo. 88'  '"" 88'  '88 88'  '88 88'  '88 
M  MMMM' .M  88       88 88.  ... 88.  .88 88       88.  .88 
M       .MM  dP '88888P' '88888P' '88888P' dP       '88888P8 
MMMMMMMMMMM 
MM'"""""'MM                                     
M' .mmm. 'M                                     
M  MMMMMMMM  88d888b. .d8888b. dP    dP 88d888b. 
M  MMM   'M  88'  '88 88'  '88 88    88 88'  '88 
M. 'MMM' .M  88       88.  .88 88.  .88 88.  .88 
MM.     .MM  dP       '88888P' '88888P' 88Y888P' 
MMMMMMMMMMM                             88
M""""""'YMM MMP"""""YMM MP""""""'MM 
M  mmmm. 'M M' .mmm. 'M M  mmmmm..M 
M  MMMMM  M M  MMMMM  M M.      'YM 
M  MMMMM  M M  MMMMM  M MMMMMMM.  M 
M  MMMM' .M M. 'MMM' .M M. .MMM'  M 
M       .MM MMb     dMM Mb.     .dM 
MMMMMMMMMMM MMMMMMMMMMM MMMMMMMMMMM 
`
blue = color.New(color.Bold, color.FgHiBlue)
red = color.New(color.Bold, color.FgHiRed)
yellow = color.New(color.Bold, color.FgHiYellow)
green = color.New(color.Bold, color.FgHiGreen)
)


func main() {
    file, err := os.Open("token.txt")
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        TOKEN_AUTH = scanner.Text()
    }
 
    if err := scanner.Err(); err != nil {
        fmt.Println(err)
    }
    clear()
    red.Println(logo)
    slow_print("by elliot")
    sleep(400)
    clear()
    
    blue.Print("group id: ")
    fmt.Scan(&gid)
    test_body, test_code := make_request(gid)
    if test_code == 401 {

        red.Println("ERROR BAD TOKEN")
        return
    }
    if test_code == 403 {
        red.Println("ERROR NOT IN GROUP")
    }
    if len(test_body) > 2 && gjson.Get(test_body, "message").String() == "You are being rate limited." {
        if gjson.Get(test_body, "retry_after").Int() <= 500 {
            green.Println("already locked")
            sleep(1000)
            check_time_body, _ := make_request(gid)
            check_time_body_retry := gjson.Get(check_time_body, "retry_after").Int()
            for check_time_body_retry <= 500 {
                blue.Println("waiting for ratelimit to end")
                check_time_body, _ = make_request(gid)
                check_time_body_retry = gjson.Get(check_time_body, "retry_after").Int()
                clear()
            }
            yellow.Println("ratelimit ended, locking group again")
        }else {
            green.Println("already locked")
            sleep(1000)
        }
    }else {
        yellow.Println("locking")
        spam(gid)

    }
    for {
        test_body, _ := make_request(gid)
        clear()
        if len(test_body) > 2 {
            if gjson.Get(test_body, "retry_after").Int() <= 500 {
                yellow.Println("locking again")
                spam(gid)
            }else if gjson.Get(test_body, "retry_after").Int() <= 5000 {
                clear()
                blue.Println("locking soon")
                spam(gid)
            }else if gjson.Get(test_body, "message").String() != "You are being rate limited." {
                yellow.Println("locking again")
                spam(gid)
                green.Println("locked")
            }else {
                green.Println("locked")
                blue.Print("remaining time: ")
                fmt.Println(gjson.Get(test_body, "retry_after").String())
            }
        }else{
            spam(gid)
        }
    }
}


func slow_print(s string){
    for c := 0; c < len(s);c++ {
        f := bufio.NewWriter(os.Stdout)
        f.Write([]byte(string(s[c])))
        f.Flush()
        sleep(70)
    }
}


func clear() {
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
    cmd = exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func sleep(t time.Duration){
    time.Sleep(t * time.Millisecond)
}


func make_request(gid string) (string, int) {
    gid = fmt.Sprintf(gid)
    httpputturl := "https://discord.com/api/v7/channels/" + gid + "/recipients/1337"

    empty, _ := json.Marshal("")

    request, err := http.NewRequest(http.MethodPut, httpputturl, bytes.NewBuffer(empty))
    if err != nil {
        panic(err)
    }

    request.Header.Set("Content-Type", "application/json; charset=UTF-8")
    request.Header.Set("authorization", TOKEN_AUTH)
    request.Header.Set("user-agent", "Discord/21295 CFNetwork/1128.0.1 Darwin/19.6.0")

    client := &http.Client{}

    response, err := client.Do(request)
    if err != nil {
        panic(err)
    }
    defer response.Body.Close()
    body_bytes, _ := ioutil.ReadAll(response.Body)
    body := string(body_bytes)
    return body, response.StatusCode

}

func spam(gid string) {
    for i := 1; i <= 50; i++ {
        go make_request(gid)
    }
    sleep(4600)
    return
}

