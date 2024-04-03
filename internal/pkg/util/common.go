package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"starland-backend/configs"
	"starland-backend/internal/pkg/httpclientutil"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const ttsPath = "/api/tts"

type GenerateVoiceRequest struct {
	Text   string `json:"text"`
	RoleId string `json:"role_id"`
}

type GenerateVoiceResponse struct {
	Code   int    `json:"code"`
	ErrMsg string `json:"err_msg"`
	Result string `json:"result"`
}

func ContainsChinese(str string) bool {
	for _, char := range str {
		if unicode.Is(unicode.Scripts["Han"], char) {
			return true
		}
	}
	return false
}

func GenerateVoice(text, roleId string) (string, error) {
	var (
		err     error
		reqData struct {
			Text   string `json:"text"`
			RoleId string `json:"role_id"`
		}

		res struct {
			Code   int    `json:"code"`
			ErrMsg string `json:"err_msg"`
			Result string `json:"result"`
		}
	)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()
	reqData.Text = text
	reqData.RoleId = roleId
	zap.S().Infof("GenerateVoice: req: %+v", reqData)
	reqBuf := new(bytes.Buffer)
	err = json.NewEncoder(reqBuf).Encode(reqData)
	if err != nil {
		return "", fmt.Errorf("req encode: %w", err)
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "POST", configs.GetConfig().Voice.Endpoint+ttsPath, reqBuf)
	if err != nil {
		return "", fmt.Errorf("post tts svc err: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	zap.S().Infof("tts text:%s", text)

	var resp *http.Response
	cl := httpclientutil.GetHttpClient()
	resp, err = cl.Do(req)
	if err != nil {
		return "", fmt.Errorf("cl.do err: %w", err)
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			fmt.Println(e)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil && errors.Is(err, http.ErrAbortHandler) {
		return "", fmt.Errorf("json decode err: %w", err)
	}

	if res.Code != 0 {
		return "", fmt.Errorf("res err: %s", res.ErrMsg)
	}
	return Base64ToVoiceFile(res.Result)
}

func Base64ToVoiceFile(base64str string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return "", fmt.Errorf("Base64ToVoiceFile: base64 decode err: %w", err)
	}
	fileName := uuid.New().String() + ".mp3"
	path := filepath.Join(configs.GetConfig().File.VoicePath, fileName)

	Mkdir(configs.GetConfig().File.VoicePath)
	tmpFile, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("Base64ToVoiceFile: make temp dir err: %w", err)
	}
	zap.S().Infof("voice file size:%d", len(decodedData))
	_, err = tmpFile.Write(decodedData)
	if err != nil {
		return "", fmt.Errorf("Base64ToVoiceFile: file write err: %w", err)
	}
	defer tmpFile.Close()
	return fileName, nil
}

func URL2FileName(url string) string {
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]

	if strings.Contains(fileName, "?") {
		queryParts := strings.Split(fileName, "?")
		fileName = queryParts[0]
	}
	zap.S().Infof("URL2FileName: %s", fileName)
	return fileName
}

func Mkdir(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, fs.ModeDir)
		if err != nil {
			return fmt.Errorf("Mkdir: Error creating directory: %w", err)
		} else {
			return nil
		}
	} else {
		zap.S().Errorf("Mkdir: Directory already exists : %s", directory)
		return nil
	}
}
