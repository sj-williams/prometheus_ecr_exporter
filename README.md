# ecr-exporter

This is a prometheus exporter for AWS ECR metrics.

## Building and running

Prerequisites:
- Docker

Building with docker:
```
$ docker build -t ecr-exporter .
```

Running:
```
$ docker run \
    -e AWS_REGION=eu-west-2 \
    -e AWS_ACCESS_KEY_ID=xxxxxxxxxxx \
    -e AWS_SECRET_ACCESS_KEY=yyyyyyyyyyy \
    -p 9606:9606 \
    ecr-exporter
```

The exporter requires an IAM role with the following permissions: `ecr:DescribeRepositories` and `ecr:ListImages` on all
ECR repositories.

## Exported Metrics

| Metric | Meaning | Labels |
| ------ | ------- | ------ |
| aws_ecr_repository_image_count | Total number of images in a repository | repository_name |
