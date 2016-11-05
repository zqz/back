package lib

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	null "gopkg.in/nullbio/null.v5"

	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models"
)

func extractIP(remoteAddr string) (string, error) {
	ip, _, err := net.SplitHostPort(remoteAddr)

	if err != nil {
		return "", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", errors.New("Failed to parse IP")
	}

	// forward := req.Header.Get("X-Forwarded-For")

	// fmt.Fprintf(w, "<p>IP: %s</p>", ip)
	// fmt.Fprintf(w, "<p>Port: %s</p>", port)
	// fmt.Fprintf(w, "<p>Forwarded for: %s</p>", forward)

	return ip, nil
}

// TrackDownload stores a record for the download.
func TrackDownload(db db.Executor, fileID string, r *http.Request, hit bool) {
	ip, err := extractIP(r.RemoteAddr)

	if err != nil {
		fmt.Println("Failed to figure out IP", err.Error())
		return
	}

	d := models.Download{
		FileID:   null.StringFrom(fileID),
		Ip:       null.StringFrom(ip),
		CacheHit: hit,
	}

	if err := d.Insert(db); err != nil {
		fmt.Println("Failed to track download", err.Error())
	}
}
