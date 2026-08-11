package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Sensor struct {
	ID        string
	StartDate time.Time
	DaysLeft  int
}

type DebugInfo struct {
	BaseURL      string
	Token        string
	AccountHash  string
	PatientID    string
	RawJSON      string
	LastHTTPCode int
}

type GlucoseData struct {
	Current float64
	Trend   int
	History []float64
	Sensors []Sensor
	Debug   DebugInfo
}

// Estructuras con las claves EXACTAS de LibreLinkUp API v4
type GlucoseMeasurement struct {
	ValueInMgPerDl float64 `json:"ValueInMgPerDl"`
	Value          float64 `json:"Value"`
	TrendArrow     int     `json:"TrendArrow"`
}

type LoginResponse struct {
	Status int `json:"status"`
	Data   struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		AuthTicket struct {
			Token string `json:"token"`
		} `json:"authTicket"`
	} `json:"data"`
}

type ConnectionsResponse struct {
	Status int `json:"status"`
	Data   []struct {
		PatientId string `json:"patientId"`
	} `json:"data"`
}

type GraphResponse struct {
	Status int `json:"status"`
	Data   struct {
		Connection struct {
			GlucoseItem GlucoseMeasurement `json:"glucoseItem"` // ¡Dato en vivo actual!
		} `json:"connection"`
		GraphData []GlucoseMeasurement `json:"graphData"` // Historial de puntos
	} `json:"data"`
}

type apiCache struct {
	Token       string
	AccountHash string
	PatientID   string
	ExpiresAt   time.Time
}

var localCache apiCache

func clearCache() {
	localCache = apiCache{}
}

func getBaseURL(region string) string {
	region = strings.ToLower(strings.TrimSpace(region))
	switch region {
	case "cl", "chile", "la", "latam":
		return "https://api-la.libreview.io"
	case "eu", "europa", "es":
		return "https://api-eu.libreview.io"
	case "us", "usa", "eeuu":
		return "https://api-us.libreview.io"
	case "ae":
		return "https://api-ae.libreview.io"
	default:
		return "https://api.libreview.io"
	}
}

func getAccountHash(userID string) string {
	hash := sha256.Sum256([]byte(userID))
	return hex.EncodeToString(hash[:])
}

func setHeaders(req *http.Request, token string, accountHash string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header["version"] = []string{"4.16.0"}
	req.Header["product"] = []string{"llu.ios"}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if accountHash != "" {
		req.Header["account-id"] = []string{accountHash}
	}
}

func Fetch(email, password, region string) (GlucoseData, error) {
	baseURL := getBaseURL(region)
	client := &http.Client{Timeout: 12 * time.Second}
	var debug DebugInfo
	debug.BaseURL = baseURL

	if localCache.Token == "" || time.Now().After(localCache.ExpiresAt) {
		loginBody, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
		})

		reqLogin, _ := http.NewRequest("POST", baseURL+"/llu/auth/login", bytes.NewBuffer(loginBody))
		setHeaders(reqLogin, "", "")

		respLogin, err := client.Do(reqLogin)
		if err != nil {
			return GlucoseData{}, fmt.Errorf("error de red en PASO 1 (Login): %w", err)
		}
		defer respLogin.Body.Close()

		loginBytes, _ := io.ReadAll(respLogin.Body)
		if respLogin.StatusCode != http.StatusOK {
			return GlucoseData{}, fmt.Errorf("PASO 1 (Login) falló (HTTP %d): %s", respLogin.StatusCode, string(loginBytes))
		}

		var loginData LoginResponse
		if err := json.Unmarshal(loginBytes, &loginData); err != nil {
			return GlucoseData{}, fmt.Errorf("error decodificando login: %w", err)
		}

		localCache.Token = loginData.Data.AuthTicket.Token
		rawUserID := loginData.Data.User.ID
		if rawUserID == "" {
			return GlucoseData{}, fmt.Errorf("error: User ID no encontrado")
		}
		localCache.AccountHash = getAccountHash(rawUserID)

		reqConn, _ := http.NewRequest("GET", baseURL+"/llu/connections", nil)
		setHeaders(reqConn, localCache.Token, localCache.AccountHash)

		respConn, err := client.Do(reqConn)
		if err != nil {
			return GlucoseData{}, fmt.Errorf("error de red en PASO 2 (Conexiones): %w", err)
		}
		defer respConn.Body.Close()

		connBytes, _ := io.ReadAll(respConn.Body)
		if respConn.StatusCode != http.StatusOK {
			return GlucoseData{}, fmt.Errorf("PASO 2 (Conexiones) falló (HTTP %d): %s", respConn.StatusCode, string(connBytes))
		}

		var connData ConnectionsResponse
		if err := json.Unmarshal(connBytes, &connData); err != nil {
			return GlucoseData{}, fmt.Errorf("error decodificando conexiones: %w", err)
		}

		if len(connData.Data) == 0 {
			return GlucoseData{}, fmt.Errorf("sin pacientes conectados")
		}

		localCache.PatientID = connData.Data[0].PatientId
		localCache.ExpiresAt = time.Now().Add(4 * time.Hour)
	}

	debug.Token = localCache.Token
	debug.AccountHash = localCache.AccountHash
	debug.PatientID = localCache.PatientID

	graphURL := fmt.Sprintf("%s/llu/connections/%s/graph", baseURL, localCache.PatientID)
	reqGraph, _ := http.NewRequest("GET", graphURL, nil)
	setHeaders(reqGraph, localCache.Token, localCache.AccountHash)

	respGraph, err := client.Do(reqGraph)
	if err != nil {
		return GlucoseData{}, fmt.Errorf("error de red en PASO 3 (Gráfico): %w", err)
	}
	defer respGraph.Body.Close()

	if respGraph.StatusCode == http.StatusUnauthorized {
		clearCache()
		return GlucoseData{}, fmt.Errorf("sesión expirada, reconectando...")
	}

	graphBytes, _ := io.ReadAll(respGraph.Body)
	debug.RawJSON = string(graphBytes)
	debug.LastHTTPCode = respGraph.StatusCode

	if respGraph.StatusCode != http.StatusOK {
		return GlucoseData{}, fmt.Errorf("PASO 3 (Gráfico) falló (HTTP %d): %s", respGraph.StatusCode, string(graphBytes))
	}

	var graphResponse GraphResponse
	if err := json.Unmarshal(graphBytes, &graphResponse); err != nil {
		return GlucoseData{}, fmt.Errorf("error procesando json de glucosa: %w", err)
	}

	// 1. Extraer glucosa actual priorizando ValueInMgPerDl y luego Value
	// Extraer la tendencia en vivo desde la conexión actual
	liveItem := graphResponse.Data.Connection.GlucoseItem
	current := liveItem.ValueInMgPerDl
	if current == 0 {
		current = liveItem.Value
	}
	trend := liveItem.TrendArrow

	// 2. Extraer historial preservando el orden cronológico estricto (de más antiguo a más reciente)
	var history []float64
	for _, item := range graphResponse.Data.GraphData {
		v := item.ValueInMgPerDl
		if v == 0 {
			v = item.Value
		}
		if v > 0 {
			history = append(history, v)
		}
	}

	// Fallback si la glucosa actual no llegó en GlucoseItem
	if current == 0 && len(history) > 0 {
		current = history[len(history)-1]
	}

	mockSensors := []Sensor{
		{ID: "Sensor Conectado", StartDate: time.Now().Add(-24 * time.Hour * 4), DaysLeft: 10},
	}

	return GlucoseData{
		Current: current,
		Trend:   trend,
		History: history,
		Sensors: mockSensors,
		Debug:   debug,
	}, nil
}
