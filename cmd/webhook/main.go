package main

import (
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
	"os"
	"solarland/infra/annunciation/nimitz/pkg/config"
	"solarland/infra/annunciation/nimitz/pkg/utils"
	"solarland/infra/annunciation/nimitz/pkg/webhook/mutating"

	"github.com/oklog/oklog/pkg/group"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
)

// Main is the main program
type Main struct {
	flags *Flags
	cfg   *viper.Viper
}

func (m *Main) Run() error {
	// Create our mutator
	mt := kwhmutating.MutatorFunc(mutating.ImagePodMutator)

	mcfg := kwhmutating.WebhookConfig{
		ID:      "nimitz",
		Obj:     &corev1.Pod{},
		Mutator: mt,
	}
	wh, err := kwhmutating.NewWebhook(mcfg)

	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("error creating webhook: %s, %s \n", err, os.Stderr)
		os.Exit(1)
	}

	// Get the handler for our webhook.
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh})
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("error creating webhook handler: %s, %s \n", err, os.Stderr)
		os.Exit(1)
	}

	// Create the servers and set them listening
	mux := http.NewServeMux()
	mux.Handle("/", whHandler)

	var g group.Group
	{
		// Serve webhooks
		g.Add(func() error {
			log.Printf("Listening on %s \n", m.flags.ListenAddress)
			return http.ListenAndServeTLS(
				m.flags.ListenAddress,
				m.flags.CertFile,
				m.flags.KeyFile,
				mux,
			)
		}, func(err error) {
			if err != nil {
				sentry.CaptureException(err)
				log.Printf("error received: %s\n", err)
			}
			log.Println("app finished successfully")
		},
		)
	}

	utils.HandleSignal(&g)
	log.Println("exit ", g.Run())
	return nil
}

func main() {
	cfg, _ := config.Initialize()
	flags := NewFlags(cfg)
	utils.SentrySetup()
	m := Main{
		flags: flags,
		cfg:   cfg,
	}
	err := m.Run()
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("%s, %s\n", err, os.Stderr)
		os.Exit(1)
	}
	os.Exit(0)
}
