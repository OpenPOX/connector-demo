package confluence

import (
	"net/url"
	"testing"
)

var p *Provider
var state = "d9f4e3a4-5b1e-4c7a-8f1e-8b2c3d4e5f6g"

func Init() {
	clientID := "xxx"
	clientSecret := "xxx"

	scopes := []string{
		"read:page:confluence",
		"read:space:confluence",
		"read:space.permission:confluence",
		"read:content-details:confluence",
		"read:content:confluence",
		"read:space-details:confluence",
		"read:attachment:confluence",
		"read:content.metadata:confluence",
		"offline_access",
	}
	p = New(clientID, clientSecret, "http://localhost:8080/auth/confluence/callback", scopes...)
}

func TestSession(t *testing.T) {
	Init()
	session, err := p.BeginAuth(state)
	if err != nil {
		t.Fatal(err)
	}
	authUrl, err := session.GetAuthURL()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("authUrl: %s", authUrl)
}

func TestExchangeCode(t *testing.T) {
	Init()
	code := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJiZjMwNTgwNy00MDZhLTRhNzMtODRjMS02YzFhMzg0MTA2MGMiLCJzdWIiOiI3MTIwMjA6YWJhN2EwZjgtNDZkOC00YWFkLTgxMDUtNjVhNDk5ZDI1M2JlIiwibmJmIjoxNzU3NTg1Nzc3LCJpc3MiOiJhdXRoLmF0bGFzc2lhbi5jb20iLCJpYXQiOjE3NTc1ODU3NzcsImV4cCI6MTc1NzU4NjA3NywiYXVkIjoiSWlWQnZZR0xNQW1yUlJCRU1OS3F3R29rU2g1ZGxFUTkiLCJjbGllbnRfYXV0aF90eXBlIjoiUE9TVCIsImh0dHBzOi8vaWQuYXRsYXNzaWFuLmNvbS92ZXJpZmllZCI6dHJ1ZSwiaHR0cHM6Ly9pZC5hdGxhc3NpYW4uY29tL3VqdCI6ImJmMzA1ODA3LTQwNmEtNGE3My04NGMxLTZjMWEzODQxMDYwYyIsInNjb3BlIjpbInJlYWQ6Y29udGVudC5tZXRhZGF0YTpjb25mbHVlbmNlIiwicmVhZDpjb250ZW50LWRldGFpbHM6Y29uZmx1ZW5jZSIsInJlYWQ6cGFnZTpjb25mbHVlbmNlIiwib2ZmbGluZV9hY2Nlc3MiLCJyZWFkOnNwYWNlLWRldGFpbHM6Y29uZmx1ZW5jZSIsInJlYWQ6c3BhY2UucGVybWlzc2lvbjpjb25mbHVlbmNlIiwicmVhZDpzcGFjZTpjb25mbHVlbmNlIiwicmVhZDphdHRhY2htZW50OmNvbmZsdWVuY2UiLCJyZWFkOmNvbnRlbnQ6Y29uZmx1ZW5jZSJdLCJodHRwczovL2lkLmF0bGFzc2lhbi5jb20vYXRsX3Rva2VuX3R5cGUiOiJBVVRIX0NPREUiLCJodHRwczovL2lkLmF0bGFzc2lhbi5jb20vaGFzUmVkaXJlY3RVcmkiOnRydWUsImh0dHBzOi8vaWQuYXRsYXNzaWFuLmNvbS9zZXNzaW9uX2lkIjoiMjNhYjU2ZmItNDU5OC00ZWNjLWEwOGEtNDc2NTY3MDM4MGQxIiwiaHR0cHM6Ly9pZC5hdGxhc3NpYW4uY29tL3Byb2Nlc3NSZWdpb24iOiJ1cy13ZXN0LTIifQ.-D8lC41dgxtxIsL-s3IMES30faRnmsi6yrSJVKOtu0g"
	session, err := p.BeginAuth(state)
	if err != nil {
		t.Fatalf("begin auth error: %v", err)
	}
	params := url.Values{
		"code": {code},
	}
	token, err := session.Authorize(p, params)
	if err != nil {
		t.Fatalf("authorize error: %v", err)
	}
	t.Logf("token :%v", token)
	t.Logf("session =%+v", session)
}
