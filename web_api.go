package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "log"
    "net/http"
)

func HttpServer() {
    e := gin.Default()
    e.POST("/cmd", CMDHandler)
    err := e.Run(":9898")
    if err != nil {
        log.Println(err)
    }
}
/**
RequestBody:{
    "cmd":"",
    "args":{},
    "page":xxx,
    "pageSize":xxx
}
*/

const (
    REQUEST_SUCCESS   = 20000
    READBODY_ERR    = 20001
    GETCMDBYBODY_ERR = 20002
    
    PALY_RUN_CMD_ERR = 20003
)

var player = MocpPlayer{}

func CMDHandler(context *gin.Context) {
    body,err := ioutil.ReadAll(context.Request.Body)
    if err != nil {
        ResponseError(context, READBODY_ERR, "获得Body失败")
        return
    }
    cmd, err := GetCMDByBody(body)
    if err != nil {
        ResponseError(context, GETCMDBYBODY_ERR, "")
    }
    
    if cmd == "Pause" {
        err = player.Pause()
    }else if cmd == "UnPause" {
        err = player.UnPause()
    }else if cmd == "TogglePause" {
        err = player.TogglePause()
    }else if cmd == "Stop" {
        err = player.Stop()
    }
    if err != nil {
        ResponseError(context, PALY_RUN_CMD_ERR ,err.Error())
    }
    ResponseSuccess(context, cmd, nil)
}

func GetCMDByBody(body []byte) (string, error) {
    header := struct {
        CMD string `json:"cmd"`
    }{}
    err := json.Unmarshal(body, &header)
    if err != nil {
        return "", err
    }
    return header.CMD, nil
}

func ResponseError(c *gin.Context, status int, message string) {
    c.JSON(http.StatusOK, gin.H{
        "Status": status,
        "Message": message,
    })
    return
}

func ResponseSuccess(c *gin.Context, message string, data interface{}){
    c.JSON(http.StatusOK, gin.H{
        "Status": REQUEST_SUCCESS,
        "Message":message,
        "Data":data})
}
