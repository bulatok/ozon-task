package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"

	"github.com/bulatok/ozon-task/pkg/pb"
)

const (
	originalLink = "https://vk.com/feed?section=updates"
)

type Config struct {
	Addr string `yaml:"address"`
}

func getConfig() (*Config, error) {
	config := &Config{}

	f, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}

	d := yaml.NewDecoder(f)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(conf.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewLinksClient(conn)

	ctx := context.Background()

	shortResp, err := client.ShortLink(ctx, &pb.ShortLinkRequest{
		OriginalLink: originalLink,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OriginalLink response", shortResp.String())

	getOrigResp, err := client.GetOriginalLink(ctx, &pb.GetOriginalRequest{
		ShortLink: shortResp.ShortLink,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetOriginalLink response", getOrigResp.String())

	if getOrigResp.OriginalLink != originalLink {
		panic("does not work")
	}

	// shutdown
	sigMain := make(chan os.Signal, 1)
	signal.Notify(sigMain, syscall.SIGTERM, syscall.SIGINT)
	<-sigMain

	if err := conn.Close(); err != nil {
		log.Println("could not close grpc client connection", err)
	}
}
