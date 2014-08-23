package ocrworker

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/couchbaselabs/logg"
)

type OcrHttpMultipartHandler struct {
	RabbitConfig RabbitConfig
}

func NewOcrHttpMultipartHandler(r RabbitConfig) *OcrHttpHandler {
	return &OcrHttpMultipartHandler{
		RabbitConfig: r,
	}
}

func (s *OcrHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	logg.LogTo("OCR_HTTP", "request to ocr-file-upload")

	switch r.Method {
	case "POST":
		h := r.Header.Get("Content-Type")
		logg.LogTo("OCR_HTTP", "content type: %v", h)

		contentType, attrs, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		logg.LogTo("OCR_HTTP", "content type: %v", contentType)

		if !strings.HasPrefix(h, "multipart/related") {
			http.Error(w, "Expected multipart related", 500)
			return
		}

		reader := multipart.NewReader(r.Body, attrs["boundary"])

		ocrReq := OcrRequest{}

		for {
			part, err := reader.NextPart()
			defer part.Close()

			if err == io.EOF {
				break
			}
			var body Body
			contentTypeOuter := mainPart.Header["Content-Type"][0]
			contentType, attrs, _ := mime.ParseMediaType(contentTypeOuter)
			switch contentType {
			case "application/json":
				decoder := json.NewDecoder(part)
				err := decoder.Decode(&ocrReq)
				if err != nil {
					logg.LogError(err)
					http.Error(w, "Unable to unmarshal json", 500)
					return
				}

			default:
				if !strings.HasPrefix(contentType, "image") {

					http.Error(w, "Expected content-type to start with image/", 500)
					return
				}

				// dump part to output (for now ..)

				partContents, err := ioutil.ReadAll(part)
				if err != nil {
					logg.LogTo("OCR_HTTP", "failed to read mime part: %v", part)
					return err
				}
				logg.LogTo("OCR_HTTP", "partContents: %v", partContents)

			}

		}

	}

	/*
		logg.LogTo("OCR_HTTP", "serveHttp called")
		defer req.Body.Close()

		ocrReq := OcrRequest{}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&ocrReq)
		if err != nil {
			logg.LogError(err)
			http.Error(w, "Unable to unmarshal json", 500)
			return
		}

		ocrClient, err := NewOcrRpcClient(s.RabbitConfig)
		if err != nil {
			logg.LogError(err)
			http.Error(w, "Unable to create rpc client", 500)
			return
		}

		decodeResult, err := ocrClient.DecodeImage(ocrReq)

		if err != nil {
			logg.LogError(err)
			http.Error(w, "Unable to perform OCR decode", 500)
			return
		}

		logg.LogTo("OCR_HTTP", "decodeResult: %v", decodeResult)

		logg.LogTo("OCR_HTTP", "ocrReq: %v", ocrReq)
		fmt.Fprintf(w, decodeResult.Text)
	*/

}
