resource "cyral_integration_logging" "cloud_watch_integration" {
  name = "my-cloudwatch-integration"
  cloud_watch {
    region = "us-east-1"
    group  = "example-loggroup"
    stream = "example-logstream"
  }
}

resource "cyral_integration_logging" "fluent_bit_integration" {
  name = "my-fluentbit-integration"
  fluent_bit {
    # Configures a custom Fluent Bit output config to write logs to an S3 bucket
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
