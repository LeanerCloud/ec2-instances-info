package ec2instancesinfo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"
)

const secretPath = "/signed-download-token-abc123"

// stubAPI serves the API step (which returns the data file URL) and the data
// file itself, and points the package at it. dataHandler serves the data file.
func stubAPI(t *testing.T, dataHandler http.HandlerFunc) (host string) {
	t.Helper()
	var srv *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc(defaultAPIPath, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(srv.URL + secretPath))
	})
	mux.HandleFunc(secretPath, dataHandler)
	srv = httptest.NewTLSServer(mux)
	oldClient := httpClient
	httpClient = srv.Client()
	t.Cleanup(func() {
		httpClient = oldClient
		srv.Close()
	})
	return strings.TrimPrefix(srv.URL, "https://")
}

func resetData(t *testing.T) {
	t.Helper()
	oldData, oldBackup := dataBody, backupDataBody
	dataBody, backupDataBody = []byte("current"), []byte("backup")
	t.Cleanup(func() { dataBody, backupDataBody = oldData, oldBackup })
}

func inTempDir(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
}

func TestUpdateDataRejectsNon200AndKeepsData(t *testing.T) {
	resetData(t)
	inTempDir(t)
	host := stubAPI(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "link expired", http.StatusForbidden)
	})
	key := "k"

	err := UpdateData(&host, &key)
	if err == nil {
		t.Fatal("expected an error for a 403 data download")
	}
	if string(dataBody) != "current" || string(backupDataBody) != "backup" {
		t.Errorf("data was replaced by a non-200 body: data=%q backup=%q", dataBody, backupDataBody)
	}
	if strings.Contains(err.Error(), secretPath) {
		t.Errorf("error leaks the download URL: %v", err)
	}
	if _, statErr := os.Stat("instances_dump.json"); statErr == nil {
		t.Error("a dump file was written for a failed download")
	}
}

func TestUpdateDataTransportErrorOmitsURL(t *testing.T) {
	resetData(t)
	inTempDir(t)
	srv := httptest.NewTLSServer(http.NotFoundHandler())
	dead := srv.URL + secretPath
	srv.Close() // the data URL now refuses connections
	api := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(dead))
	}))
	defer api.Close()
	oldClient := httpClient
	httpClient = api.Client()
	defer func() { httpClient = oldClient }()
	host, key := strings.TrimPrefix(api.URL, "https://"), "k"

	err := UpdateData(&host, &key)
	if err == nil {
		t.Fatal("expected a transport error")
	}
	if strings.Contains(err.Error(), secretPath) {
		t.Errorf("error leaks the download URL: %v", err)
	}
	if string(dataBody) != "current" {
		t.Errorf("data changed on a failed download: %q", dataBody)
	}
}

func TestUpdateDataSucceedsAndSwapsData(t *testing.T) {
	resetData(t)
	inTempDir(t)
	host := stubAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("fresh"))
	})
	key := "k"

	if err := UpdateData(&host, &key); err != nil {
		t.Fatal(err)
	}
	if string(dataBody) != "fresh" || string(backupDataBody) != "current" {
		t.Errorf("data=%q backup=%q", dataBody, backupDataBody)
	}
	if b, err := os.ReadFile("instances_dump.json"); err != nil || string(b) != "fresh" {
		t.Errorf("dump file = %q, %v", b, err)
	}
}

func TestUpdateDataDoesNotPanicWhenDumpCannotBeWritten(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs a read-only directory the process cannot write to")
	}
	resetData(t)
	dir := t.TempDir()
	wd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chmod(dir, 0o700) }()
	host := stubAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("fresh"))
	})
	key := "k"

	if err := UpdateData(&host, &key); err != nil {
		t.Fatalf("a dump-file failure must not fail the update: %v", err)
	}
	if string(dataBody) != "fresh" {
		t.Errorf("data = %q", dataBody)
	}
}

func TestUpdateDataBadDataURLOmitsURL(t *testing.T) {
	resetData(t)
	inTempDir(t)
	api := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("https://example.invalid" + secretPath + "\x7f"))
	}))
	defer api.Close()
	oldClient := httpClient
	httpClient = api.Client()
	defer func() { httpClient = oldClient }()
	host, key := strings.TrimPrefix(api.URL, "https://"), "k"

	err := UpdateData(&host, &key)
	if err == nil {
		t.Fatal("expected an error for a malformed data URL")
	}
	if strings.Contains(err.Error(), secretPath) {
		t.Errorf("error leaks the download URL: %v", err)
	}
}
