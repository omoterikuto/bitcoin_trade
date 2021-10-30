package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

const ValidKey = "valid_key"
const ValidSecret = "valid_secret"

type ResponseBody struct {
	Text string `json:"text"`
}

// レスポンスにかかる時間と、*http.Responseを引数で取れるようにしました
// *http.Responseがnilであれば、正常な*http.Responseを返します
func client(t *testing.T, respTime time.Duration, resp *http.Response) *http.Client {
	t.Helper()

	body := ResponseBody{
		Text: "hello",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	failedBody := ResponseBody{
		Text: "failed",
	}
	fb, err := json.Marshal(failedBody)
	if err != nil {
		t.Fatal(err)
	}

	return NewTestClient(func(req *http.Request) *http.Response {
		time.Sleep(respTime)

		if resp != nil {
			return resp
		}

		if req.Header.Get("ACCESS-KEY") != ValidKey {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       ioutil.NopCloser(bytes.NewBuffer(fb)),
				Header:     make(http.Header),
			}
		}

		timestamp := strconv.FormatInt(time.Now().Unix(), 10)

		message := timestamp + "GET" + "/v1/me/getbalance"

		mac := hmac.New(sha256.New, []byte(ValidSecret))
		mac.Write([]byte(message))
		sign := hex.EncodeToString(mac.Sum(nil))
		if req.Header.Get("ACCESS-SIGN") != sign {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       ioutil.NopCloser(bytes.NewBuffer(fb)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})
}

func TestDoRequest(t *testing.T) {
	testCases := map[string]struct {
		key                  string
		secret               string
		client               *http.Client
		expectHasError       bool
		expectedErrorMessage string
		expectedText         string
	}{
		"normal": { // 正常なRequest & Response
			key:            ValidKey,
			secret:         ValidSecret,
			client:         client(t, 0, nil),
			expectHasError: false,
			expectedText:   "hello",
		},
		"invalid key": { // 無効なAuth Keyを投げたとき
			key:                  "invalid_key",
			secret:               ValidSecret,
			client:               client(t, 0, nil),
			expectHasError:       true,
			expectedErrorMessage: "bad response status code 401",
		},
		"invalid secret": { // 無効なsecretで投げたとき
			key:                  ValidKey,
			secret:               "invalid_secret",
			client:               client(t, 0, nil),
			expectHasError:       true,
			expectedErrorMessage: "bad response status code 401",
		},
		"internal server error response": { // 5xxが返ってくるとき
			key:    ValidKey,
			secret: ValidSecret,
			client: client(t, 0, &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       nil,
				Header:     make(http.Header),
			}),
			expectHasError:       true,
			expectedErrorMessage: "bad response status code 500",
		},
		"plain text response": { // プレーンテキストが返って来る場合はjson.Unmarshalできずに死ぬはず
			key:    ValidKey,
			secret: ValidSecret,
			client: client(t, 0, &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString("bad")),
				Header:     make(http.Header),
			}),
			expectHasError:       true,
			expectedErrorMessage: "cannot unmarshal to ResponseBody",
		},
		"long response time": { // レスポンスに3秒かかるのでタイムアウトになるはず
			key:                  ValidKey,
			secret:               ValidSecret,
			client:               client(t, 3*time.Second, nil),
			expectHasError:       true,
			expectedErrorMessage: "bad response status code 401",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			api := New(testCase.key, testCase.secret, OptionHTTPClient(testCase.client))

			resp, err := api.doRequest("GET", "me/getbalance", map[string]string{}, nil)

			var body ResponseBody
			// プレーンテキストが返って来る場合の検証
			if unmarshalErr := json.Unmarshal(resp, &body); unmarshalErr != nil {
				if testCase.expectedErrorMessage == "cannot unmarshal to ResponseBody" {
					return
				}
			}

			// エラーになることを期待していた場合は期待したエラーメッセージかどうか検証する
			if testCase.expectHasError {
				if err.Error() != testCase.expectedErrorMessage {
					t.Errorf("unexpected error message. expected '%s', actual '%s'", testCase.expectedErrorMessage, err.Error())
				}

				if err == nil {
					t.Errorf("expected error but no errors ouccured")
					return
				}
				return
			}

			if body.Text != testCase.expectedText {
				t.Errorf("unexpected response's text. expected '%s', actual '%s'", testCase.expectedText, resp)
				return
			}

			if err != nil {
				t.Errorf(err.Error())
				return
			}

		})
	}
}
