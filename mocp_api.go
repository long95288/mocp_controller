package main

import (
    "fmt"
    "log"
    "os/exec"
    "strings"
)

/**
Music On Console (version 2.6-alpha3)
Usage: mocp [OPTIONS] [FILE|DIR ...]

General options:
  -D, --debug                       Turn on logging to a file
  -M, --moc-dir=DIR                 Use the specified MOC directory instead of the default
  -m, --music-dir                   Start in MusicDir
  -C, --config=FILE                 Use the specified config file instead of the default (conflicts with '--no-config')
      --no-config                   Use program defaults rather than any config file (conflicts with '--config')
  -O, --set-option='NAME=VALUE'     Override the configuration option NAME with VALUE
  -F, --foreground                  Run the server in foreground (logging to stdout)
  -S, --server                      Only run the server
  -R, --sound-driver=DRIVERS        Use the first valid sound driver
  -A, --ascii                       Use ASCII characters to draw lines
  -T, --theme=FILE                  Use the selected theme file (read from ~/.moc/themes if the path is not absolute)
  -y, --sync                        Synchronize the playlist with other clients
  -n, --nosync                      Don't synchronize the playlist with other clients

Server commands:
  -P, --pause                       Pause
  -U, --unpause                     Unpause
  -G, --toggle-pause                Toggle between playing and paused
  -s, --stop                        Stop playing
  -f, --next                        Play the next song
  -r, --previous                    Play the previous song
  -k, --seek=N                      Seek by N seconds (can be negative)
  -j, --jump=N{%,s}                 Jump to some position in the current track
  -v, --volume=[+,-]LEVEL           Adjust the PCM volume
  -x, --exit                        Shutdown the server
  -a, --append                      Append the files/directories/playlists passed in the command line to playlist
  -e, --recursively                 Alias for --append
  -q, --enqueue                     Add the files given on command line to the queue
  -c, --clear                       Clear the playlist
  -p, --play                        Start playing from the first item on the playlist
  -l, --playit                      Play files given on command line without modifying the playlist
  -t, --toggle=CONTROL              Toggle a control (shuffle, autonext, repeat)
  -o, --on=CONTROL                  Turn on a control (shuffle, autonext, repeat)
  -u, --off=CONTROL                 Turn off a control (shuffle, autonext, repeat)
  -i, --info                        Print information about the file currently playing
  -Q, --format=FORMAT               Print formatted information about the file currently playing

Miscellaneous options:
  -V, --version                     Print version information
      --echo-args                   Print POPT-interpreted arguments
      --usage                       Print brief usage
  -h, --help                        Print extended usage

Environment variables:

  MOCP_OPTS                         Additional command line options
  MOCP_POPTRC                       List of POPT configuration files
*/

const execScript ="mocp"

type Info struct {
    State string `json:"state"`
    File  string `json:"file"`
    Title string `json:"title"`
    Artist string `json:"artist"`
    SongTitle string `json:"song_title"`
    Album string `json:"album"`
    TotalTime string `json:"total_time"`
    TimeLeft string `json:"time_left"`
    TotalSec string `json:"total_sec"`
    CurrentTime string `json:"current_time"`
    Bitrate string `json:"bitrate"`
    AvgBitrate string `json:"avg_bitrate"`
    Rate string `json:"rate"`
}

type MocpPlayer struct {
    CtlList  []string
    MocpInfo Info
}

func NewMocpPlayer() (MocpPlayer) {
    m := MocpPlayer{}
    m.CtlList = []string{"shuffle", "autonext", "repeat"}
    log.Println("init Mocp Player")
    return m
}



func (m *MocpPlayer) Pause() error {
    return exec.Command(execScript, "-P").Run()
}

func (m *MocpPlayer) UnPause() error{
    return exec.Command(execScript, "-U").Run()
}

func (m *MocpPlayer) TogglePause() error {
    return exec.Command(execScript, "-G").Run()
}

func (m *MocpPlayer) Stop() error {
    return exec.Command(execScript, "-s").Run()
}

func (m *MocpPlayer) Next() error {
    return exec.Command(execScript, "-f").Run()
}

func (m *MocpPlayer) Previous() error {
    return exec.Command(execScript, "-r").Run()
}

func (m *MocpPlayer) Seek(n int) error {
    return exec.Command(execScript, []string{"-k", fmt.Sprintf("%d", n)}...).Run()
}

func (m *MocpPlayer) Volume(level int) error {
    return exec.Command(execScript, []string{"-v", fmt.Sprintf("%d", level)}...).Run()
}

func (m *MocpPlayer) Exit() error {
    return exec.Command(execScript, "-x").Run()
}

func (m *MocpPlayer) Play() error {
    return exec.Command(execScript, "-p").Run()
}

func (m *MocpPlayer) ToggleCtl(ctl string) error {
    return exec.Command(execScript, []string{"-t", ctl}...).Run()
}

func (m *MocpPlayer) TurnOnCtl(ctl string) error {
    return exec.Command(execScript, []string{"-o", ctl}...).Run()
}

func (m *MocpPlayer) TurnOffCtl(ctl string) error {
    return exec.Command(execScript, []string{"-u", ctl}...).Run()
}

func (m *MocpPlayer) Info() (*Info, error) {
    cmd := exec.Command(execScript, "-i")
    infoStr, err := cmd.CombinedOutput()
    if nil != err {
        return nil,err
    }
    
    for _, str := range strings.Split(string(infoStr), "\n") {
        // 直接找到第一个":"作为分割点
        key := ""
        value := ""
        for i, v := range str {
            if v == ':' {
                key = strings.TrimSpace(str[:i])
                if i < len(str) {
                    value = strings.TrimSpace(str[i + 1:])
                }
                break
            }
        }
        // log.Printf("key:%s value:%s\n", key, value)
        
        if "State" == key {
            m.MocpInfo.State = value
        }else if "File" == key {
            m.MocpInfo.File = value
        }else if "Title" == key {
            m.MocpInfo.Title = value
        }else if "Artist" == key {
            m.MocpInfo.Artist = value
        }else if "SongTitle" == key {
            m.MocpInfo.SongTitle = value
        }else if "Album" == key {
            m.MocpInfo.Album = value
        }else if "TotalTime" == key {
            m.MocpInfo.TotalTime = value
        }else if "TimeLeft" == key {
            m.MocpInfo.TimeLeft = value
        }else if "TotalSec" == key {
            m.MocpInfo.TotalSec = value
        }else if "CurrentTime" == key {
            m.MocpInfo.CurrentTime = value
        }else if "Bitrate" == key {
            m.MocpInfo.Bitrate = value
        }else if "AvgBitrate" == key {
            m.MocpInfo.AvgBitrate = value
        }else if "Rate" == key {
            m.MocpInfo.Rate = value
        }else{
            log.Printf("Unknown key[%s], value[%s]\n", key, value)
        }
    }
    return &m.MocpInfo, nil
}




