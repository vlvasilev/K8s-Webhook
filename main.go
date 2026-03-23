package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/slackhq/simple-kubernetes-webhook/pkg/configuration"
	"github.com/slackhq/simple-kubernetes-webhook/pkg/servlet"
)

func main() {
	setLogger()

	curatorConfigPath := flag.String("config", "/etc/admission-webhook/config.yaml", "A path to the webhooks's configuration file")
	flag.Parse()

	// read the configuration file
	cfg, err := configuration.ReadConfigurationFromFile(*curatorConfigPath)
	if err != nil {
		logrus.Fatalf("could not read configuration file %q: %v", *curatorConfigPath, err)
	}

	// create the servlet
	servlet := servlet.NewServlet(cfg)

	// handle our core application
	http.HandleFunc("/validate-pods", servlet.ServeValidatePods)
	http.HandleFunc("/mutate-pods", servlet.ServeMutatePods)
	http.HandleFunc("/health", servlet.ServeHealth)

	// start the server
	// listens to clear text http on port 8080 unless TLS env var is set to "true"
	if os.Getenv("TLS") == "true" {
		cert := "/etc/admission-webhook/tls/tls.crt"
		key := "/etc/admission-webhook/tls/tls.key"
		logrus.Print("Listening on port 443...")
		logrus.Fatal(http.ListenAndServeTLS(":443", cert, key, nil))
	} else {
		logrus.Print("Listening on port 8080...")
		logrus.Fatal(http.ListenAndServe(":8080", nil))
	}
}

// setLogger sets the logger using env vars, it defaults to text logs on
// debug level unless otherwise specified
func setLogger() {
	logrus.SetLevel(logrus.DebugLevel)

	lev := os.Getenv("LOG_LEVEL")
	if lev != "" {
		llev, err := logrus.ParseLevel(lev)
		if err != nil {
			logrus.Fatalf("cannot set LOG_LEVEL to %q", lev)
		}
		logrus.SetLevel(llev)
	}

	if os.Getenv("LOG_JSON") == "true" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
