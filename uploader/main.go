package main

import (
    "fmt"
    "net/http"
    "os"
    "uploader/handler"

    "github.com/sirupsen/logrus"
)

const (
    DefaultPort                   = "80"
    DefaultUploadPath             = "/tmp"
    PortEnvironmentVariable       = "C3SR_PORT"
    UploadPathEnvironmentVariable = "C3SR_UPLOAD_PATH"
)

var (
    log *logrus.Entry
)

func main() {
    logrus.SetLevel(logrus.DebugLevel)
    log = logrus.New().WithField("pkg", "api/uploader")

    port, found := os.LookupEnv(PortEnvironmentVariable)
    if !found {
        port = DefaultPort
    }
    uploadPath, found := os.LookupEnv(UploadPathEnvironmentVariable)
    if !found {
        uploadPath = DefaultUploadPath
    }
    uploadHandler, err := handler.MakeUploadHandler(uploadPath, log)

    if err != nil {
        return
    }

    http.Handle("/", http.StripPrefix("/", uploadHandler))
    err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
    if err != nil {
        panic(fmt.Errorf("unable to listen: %s", err))
    }
}
