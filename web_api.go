package main

import (
    "encoding/json"
    "errors"
    "fmt"
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

var player = NewMocpPlayer()

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
    var resp interface{} = nil
    
    if cmd == "Pause" {
        err = player.Pause()
    }else if cmd == "UnPause" {
        err = player.UnPause()
    }else if cmd == "TogglePause" {
        err = player.TogglePause()
    }else if cmd == "Stop" {
        err = player.Stop()
    }else if cmd == "Volume" {
        if level,err := GetVolumeByBody(body);nil == err {
            err = player.Volume(level)
        }
    }else if cmd == "Play" {
        err = player.Play()
    } else if cmd == "GetCtlList" {
        resp = player.CtlList
    } else if cmd == "ToggleCtl" {
        if ctl, err := GetCtlByBody(body);nil == err {
            err = player.ToggleCtl(ctl)
        }
    }else if cmd == "TurnOnCtl" {
        if ctl, err := GetCtlByBody(body);nil == err {
            err = player.TurnOnCtl(ctl)
        }
    }else if cmd == "TurnOffCtl" {
        if ctl,err := GetCtlByBody(body);nil == err {
            err = player.TurnOffCtl(ctl)
        }
    }else if cmd == "GetInfo" {
        resp, err = player.Info()
        resp = string(resp.([]byte))
    }else{
        err = errors.New(fmt.Sprintf("Unknow cmd %s\n", cmd))
    }
    if err != nil {
        ResponseError(context, PALY_RUN_CMD_ERR ,err.Error())
    }
    ResponseSuccess(context, cmd, resp)
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

func GetVolumeByBody(body []byte) (level int, err error) {
    _body := struct {
        Level int `json:"level"`
    }{}
    if err = json.Unmarshal(body, &_body);err != nil {
        return 0,err
    }
    return _body.Level, nil
}

func GetCtlByBody(body []byte)(ctl string, err error) {
    _body := struct {
        Ctl string `json:"ctl"`
    }{}
    if err = json.Unmarshal(body, &body);err!=nil{
        return "", err
    }
    return _body.Ctl, nil
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
