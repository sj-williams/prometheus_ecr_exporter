package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	imagesInRepository = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_ecr_repository_image_count",
		Help: "The total number of images in an ECR repository",
	}, []string{"repository_name"})
)

func main() {
	go func() {
		svc := ecr.New(session.New())
		for {
			err := imageCounter(svc)
			if err != nil {
				log.Printf("Failed to count images: %v", err)
			}
			time.Sleep(time.Minute)
		}
	}()

	// This section will start the HTTP server and expose
	// any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Beginning to serve on port :9606")
	log.Fatal(http.ListenAndServe(":9606", nil))
}

func imageCounter(svc *ecr.ECR) error {
	return svc.DescribeRepositoriesPages(&ecr.DescribeRepositoriesInput{},
		func(page *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
			for _, repo := range page.Repositories {
				c, err := countImagesInRepository(svc, *repo.RepositoryName)
				if err == nil {
					imagesInRepository.WithLabelValues(*repo.RepositoryName).Set(float64(c))
				} else {
					log.Printf("Ignoring repository '%s', could not count images: %v", *repo.RepositoryName, err)
				}
			}
			return !lastPage
		})
}

func countImagesInRepository(svc *ecr.ECR, repositoryName string) (int, error) {
	count := 0
	err := svc.ListImagesPages(&ecr.ListImagesInput{RepositoryName: aws.String(repositoryName)},
		func(page *ecr.ListImagesOutput, lastPage bool) bool {
			count += len(page.ImageIds)
			return !lastPage
		})
	return count, err
}
