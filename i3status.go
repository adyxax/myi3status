package main
import (
    "bufio"
    "errors"
    "flag"
    "fmt"
    "log"
    "os/exec"
)

func main_loop(weather_code *string, rain_color *string) error {
    path, lookupErr := exec.LookPath("i3status")
    if lookupErr != nil {
        return lookupErr
    }
    cmd := exec.Command(path)
    stdout, pipeErr := cmd.StdoutPipe()
    if pipeErr != nil {
        return pipeErr
    }
    if err := cmd.Start(); err != nil {
        return err
    }
    scanner := bufio.NewScanner(stdout)

    /* expect '{"version": 1}' */
    scannerOk := scanner.Scan()
    if !scannerOk {
        return errors.New("scanner failed")
    }
    var t string
    t = scanner.Text()
    if t != "{\"version\":1}" {
        return errors.New("invalid header '"+t+"'")
    }
    fmt.Println(t)

    /* expect '[' */
    scannerOk2 := scanner.Scan()
    if !scannerOk2 {
        return errors.New("scanner failed")
    }
    t = scanner.Text()
    if t != "[" {
        return errors.New("invalid 2nd line '"+t+"'")
    }
    fmt.Println(t)

    first := true
    for scanner.Scan() {
        text := scanner.Text()
        if weather_code != nil {
            weather, errStatus := GetRainI3barFormat(weather_code, rain_color)
            if errStatus != nil {
                fmt.Println(text)
            } else if first {
                fmt.Println("["+weather+","+text[1:])
                first = false
            } else {
                fmt.Println(",["+weather+","+text[2:])
            }
        } else {
            fmt.Println(text)
        }
    }

    return nil
}

func main() {
    weatherCode := flag.String("weather_code", "",
        "a location code for the Pluie dans l'heure API")
    rainColor := flag.String("rain_color", "#268bd2",
        "Color to display text when it's raining")
    flag.Parse()
    if *weatherCode == "" {
        weatherCode = nil
    }
    err := main_loop(weatherCode, rainColor)
    if err != nil {
        log.Fatal(err)
    }
}

