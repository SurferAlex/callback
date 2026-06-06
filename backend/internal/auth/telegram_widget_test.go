package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestWidgetFieldsFromMap_integerFields(t *testing.T) {
	raw := `{"id":863631849,"first_name":"Artie","username":"artie","auth_date":1735689600,"hash":"placeholder"}`
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatal(err)
	}
	fields := WidgetFieldsFromMap(body)
	if fields["id"] != "863631849" {
		t.Fatalf("id=%q want 863631849", fields["id"])
	}
	if fields["auth_date"] != "1735689600" {
		t.Fatalf("auth_date=%q want 1735689600", fields["auth_date"])
	}
}

func TestVerifyWidget_roundTrip(t *testing.T) {
	botToken := "123456:TEST-BOT-TOKEN"
	authDate := time.Now().Unix()
	fields := map[string]string{
		"id":         "863631849",
		"first_name": "Artie",
		"username":   "artie",
		"auth_date":  strconv.FormatInt(authDate, 10),
	}

	var pairs []string
	for k, v := range fields {
		if strings.TrimSpace(v) == "" {
			continue
		}
		pairs = append(pairs, k+"="+v)
	}
	sort.Strings(pairs)
	dataCheck := strings.Join(pairs, "\n")

	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheck))
	fields["hash"] = hex.EncodeToString(mac.Sum(nil))

	raw, err := json.Marshal(map[string]interface{}{
		"id":         863631849,
		"first_name": "Artie",
		"username":   "artie",
		"auth_date":  authDate,
		"hash":       fields["hash"],
	})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}

	user, err := VerifyWidget(botToken, WidgetFieldsFromMap(body), time.Hour)
	if err != nil {
		t.Fatalf("VerifyWidget: %v", err)
	}
	if user.ID != 863631849 || user.FirstName != "Artie" || user.Username != "artie" {
		t.Fatalf("unexpected user: %+v", user)
	}
}
