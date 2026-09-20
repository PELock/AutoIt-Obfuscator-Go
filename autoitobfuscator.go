/******************************************************************************
 * AutoIt Obfuscator WebApi interface
 *
 * Version        : v1.5.0
 * Language       : Go
 * Author         : Bartosz Wójcik
 * Web page       : https://www.pelock.com
 *
 *****************************************************************************/

package autoitobfuscator

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	APIURL = "https://www.pelock.com/api/autoit-obfuscator/v1"

	ErrorSuccess     = 0
	ErrorInputSize   = 1
	ErrorInput       = 2
	ErrorParsing     = 3
	ErrorObfuscation = 4
	ErrorOutput      = 5
)

// Result is the parsed Web API JSON response.
type Result struct {
	Error        int    `json:"error"`
	Output       string `json:"output,omitempty"`
	Demo         bool   `json:"demo,omitempty"`
	CreditsLeft  int    `json:"credits_left,omitempty"`
	CreditsTotal int    `json:"credits_total,omitempty"`
	Expired      bool   `json:"expired,omitempty"`
	StringLimit  int    `json:"string_limit,omitempty"`
}

// Client is the AutoIt Obfuscator Web API client.
// Strategy flags default to false; set the ones you want before calling Obfuscate*.
type Client struct {
	APIKey     string
	APIURL     string
	UserAgent  string
	HTTPClient *http.Client

	EnableCompression bool

	AntiDebug    bool
	AntiVM       bool
	AntiSandbox  bool
	AntiEmulator bool

	RandomIntegers               bool
	RandomCharacters             bool
	RandomAntiRegex              bool
	RandomArrays                 bool
	RandomArraysMultidimensional bool
	RandomFunctions              bool
	RandomAutostarted            bool

	MixCodeFlow         bool
	RenameVariables     bool
	RenameFunctions     bool
	RenameFunctionCalls bool
	ShuffleFunctions    bool
	ResolveConstants    bool
	CryptNumbers        bool
	SplitStrings        bool
	ModifyStrings       bool
	CryptStrings        bool
	InsertTernaryOperators bool
}

// New creates a client. An empty or invalid key runs demo mode.
func New(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		APIURL:     APIURL,
		UserAgent:  "PELock AutoIt Obfuscator",
		HTTPClient: &http.Client{Timeout: 120 * time.Second},
	}
}

// Login returns license and quota information for the activation key.
func (c *Client) Login(ctx context.Context) (*Result, error) {
	return c.postRequest(ctx, map[string]string{"command": "login"})
}

// ObfuscateScriptFile reads a UTF-8 AutoIt script and obfuscates it.
func (c *Client) ObfuscateScriptFile(ctx context.Context, scriptFilePath string) (*Result, error) {
	source, err := os.ReadFile(scriptFilePath)
	if err != nil {
		return nil, err
	}
	if len(source) == 0 {
		return nil, fmt.Errorf("empty source file")
	}
	return c.ObfuscateScriptSource(ctx, string(source))
}

// ObfuscateScriptSource obfuscates AutoIt source code.
func (c *Client) ObfuscateScriptSource(ctx context.Context, scriptSource string) (*Result, error) {
	return c.postRequest(ctx, map[string]string{
		"command": "obfuscate",
		"source":  scriptSource,
	})
}

func (c *Client) postRequest(ctx context.Context, params map[string]string) (*Result, error) {
	if c.APIKey != "" {
		params["key"] = c.APIKey
	}

	if c.AntiDebug {
		params["anti_debug"] = "1"
	}
	if c.AntiVM {
		params["anti_vm"] = "1"
	}
	if c.AntiSandbox {
		params["anti_sandbox"] = "1"
	}
	if c.AntiEmulator {
		params["anti_emulator"] = "1"
	}

	if c.RandomIntegers {
		params["random_bucket_integers"] = "1"
	}
	if c.RandomCharacters {
		params["random_bucket_characters"] = "1"
	}
	if c.RandomAntiRegex {
		params["random_bucket_anti_regex"] = "1"
	}
	if c.RandomArrays {
		params["random_bucket_arrays"] = "1"
	}
	if c.RandomArraysMultidimensional {
		params["random_bucket_arrays_multidimensional"] = "1"
	}
	if c.RandomFunctions {
		params["random_bucket_functions"] = "1"
	}
	if c.RandomAutostarted {
		params["random_bucket_autostart"] = "1"
	}

	if c.MixCodeFlow {
		params["mix_code_flow"] = "1"
	}
	if c.RenameVariables {
		params["rename_variables"] = "1"
	}
	if c.RenameFunctions {
		params["rename_functions"] = "1"
	}
	if c.RenameFunctionCalls {
		params["rename_function_calls"] = "1"
	}
	if c.ShuffleFunctions {
		params["shuffle_functions"] = "1"
	}
	if c.ResolveConstants {
		params["resolve_const"] = "1"
	}
	if c.CryptNumbers {
		params["crypt_numbers"] = "1"
	}
	if c.SplitStrings {
		params["split_strings"] = "1"
	}
	if c.ModifyStrings {
		params["modify_strings"] = "1"
	}
	if c.CryptStrings {
		params["crypt_strings"] = "1"
	}
	if c.InsertTernaryOperators {
		params["insert_ternary_operators"] = "1"
	}

	if c.EnableCompression && params["source"] != "" {
		compressed, err := zlibCompressBase64(params["source"])
		if err != nil {
			return nil, err
		}
		params["source"] = compressed
		params["compression"] = "1"
	}

	body, err := c.postMultipart(ctx, params, nil)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("empty API response")
	}

	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if c.EnableCompression && result.Error == ErrorSuccess && result.Output != "" {
		plain, err := zlibDecompressBase64(result.Output)
		if err != nil {
			return nil, err
		}
		result.Output = plain
	}

	return &result, nil
}

type formFile struct {
	Field    string
	FileName string
	Data     []byte
}

func (c *Client) postMultipart(ctx context.Context, fields map[string]string, files []formFile) ([]byte, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	for _, f := range files {
		part, err := w.CreateFormFile(f.Field, f.FileName)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(f.Data); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	url := c.APIURL
	if url == "" {
		url = APIURL
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func zlibCompressBase64(s string) (string, error) {
	var buf bytes.Buffer
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err := zw.Write([]byte(s)); err != nil {
		_ = zw.Close()
		return "", err
	}
	if err := zw.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func zlibDecompressBase64(s string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	zr, err := zlib.NewReader(bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	defer zr.Close()
	plain, err := io.ReadAll(zr)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
