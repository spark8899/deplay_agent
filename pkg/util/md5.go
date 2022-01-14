package util

import (
    "os"
    "io"
    "bufio"
    "crypto/md5"
    "encoding/hex"
)

func EncodeMD5(value string) string {
    m := md5.New()
    m.Write([]byte(value))

    return hex.EncodeToString(m.Sum(nil))
}

func FileMD5(file string) string {
    f, err := os.Open(file)
    if err != nil {
        return ""
    }
    defer f.Close()
    r := bufio.NewReader(f)
    m := md5.New()

    _, err = io.Copy(m, r)
    if err != nil {
        return ""
    }

    return hex.EncodeToString(m.Sum(nil))
}
