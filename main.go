package main

import (
	cryptoRand "crypto/rand"
	"log/slog"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	setupLogger()

	sourceAddress, targetAddress := parseArgs()

	conn, err := net.DialUDP("udp", sourceAddress, targetAddress)
	if err != nil {
		slog.Error("Failed to dial UDP", "source", sourceAddress, "target", targetAddress, "err", err)
		os.Exit(1)
	}

	packetsCount := rand.N(40) + 20
	slog.Info("Starting to send packets", "packetsCount", packetsCount)

	for range packetsCount {
		payload := cryptoRand.Text()
		delay := rand.N(100) + 50

		bytesSent, _, err := conn.WriteMsgUDP([]byte(payload), []byte{}, nil)
		if err != nil {
			slog.Info("Error sending packet to target", "err", err)
			os.Exit(1)
		}

		slog.Info("Sent packet to target", "bytesSent", bytesSent, "payload", payload)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

func parseArgs() (sourceAddress *net.UDPAddr, targetAddress *net.UDPAddr) {
	if len(os.Args) < 3 {
		slog.Error("Not all arguments were passed. Usage: gofuncrkn sourcePort targetHost:targetPort")
		os.Exit(1)
	}

	sourcePort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		slog.Error("Failed to parse source port", "source", os.Args[1], "err", err)
		os.Exit(1)
	}

	target := strings.Split(os.Args[2], ":")
	targetHost := target[0]
	targetPort, err := strconv.Atoi(target[1])
	if err != nil {
		slog.Error("Failed to parse target port", "target", os.Args[2], "err", err)
		os.Exit(1)
	}

	slog.Info("Starting", "source", sourcePort, "targetHost", targetHost, "targetPort", targetPort)

	sourceAddress, err = net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(sourcePort))
	if err != nil {
		slog.Error("Failed to resolve local address", "host", "0.0.0.0", "port", sourcePort, "err", err)
		os.Exit(1)
	}
	targetAddress, err = net.ResolveUDPAddr("udp", targetHost+":"+strconv.Itoa(targetPort))
	if err != nil {
		slog.Error("Failed to resolve target address", "host", targetHost, "port", targetPort, "err", err)
		os.Exit(1)
	}

	slog.Info("Parsed and resolved addresses", "source", sourceAddress, "target", targetAddress)
	return sourceAddress, targetAddress
}

func setupLogger() {
	w := os.Stderr

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		}),
	))
}
