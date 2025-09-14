package buscacep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lab-google-cloud-run/internal/dto"
	"net/http"

	_ "golang.org/x/crypto/x509roots/fallback"
)

// HTTPDoer permite injetar http.Client ou mocks
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ViaCepClient struct {
	httpClient HTTPDoer
	baseURL    string
}

func NewViaCepClient(c HTTPDoer) *ViaCepClient {
	if c == nil {
		c = http.DefaultClient
	}
	return &ViaCepClient{httpClient: c, baseURL: "https://viacep.com.br/ws"}
}

func (v *ViaCepClient) GetLocationByCep(ctx context.Context, cep string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/json/", v.baseURL, cep), nil)
	if err != nil {
		return "", err
	}
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var viaCepAPI dto.ViaCEPAPIResponse
	if err := json.Unmarshal(body, &viaCepAPI); err != nil {
		return "", err
	}
	return viaCepAPI.Location, nil
}

// Função legado mantendo assinatura original usando client default
func GetLocationByCep(cep string) (string, error) {
	return NewViaCepClient(http.DefaultClient).GetLocationByCep(context.Background(), cep)
}
