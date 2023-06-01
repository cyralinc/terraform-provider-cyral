# Configures `my-sidecar-cloud-watch` to push logs to CloudWatch to a log stream named `cyral-sidecar`
# in a log group named after `cyral-example-loggroup`.
resource "cyral_sidecar" "sidecar_cloud_watch" {
  name               = "my-sidecar-cloud-watch"
  deployment_method  = "terraform"
  log_integration_id = cyral_integration_logging.cloud_watch.id
}

resource "cyral_integration_logging" "cloud_watch" {
  name = "my-cloudwatch"
  cloud_watch {
    region = "us-east-1"
    group  = "cyral-example-loggroup"
    stream = "cyral-sidecar"
  }
}

# Configures `my-sidecar-fluent-bit` to push logs to S3 to a bucket named  `example-bucket`.
resource "cyral_sidecar" "sidecar_fluent_bit" {
  name               = "my-sidecar-fluent-bit"
  deployment_method  = "terraform"
  log_integration_id = cyral_integration_logging.s3.id
}

resource "cyral_integration_logging" "s3" {
  name = "my-s3"
  fluent_bit {
    config = <<-EOF
    [OUTPUT]
      Name s3
      Match *
      Region us-east-2
      Bucket example-bucket
      Total_file_size 1M
    EOF
  }
}

# Configures a raw Elk integration with no sidecar associated.
resource "cyral_integration_logging" "elk_integration" {
  name = "my-elk-integration"
  elk {
    es_url     = "http://es.com"
    kibana_url = "http://kibana.com"
    es_credentials {
      username = "another-user"
      password = "123"
    }
  }
}
