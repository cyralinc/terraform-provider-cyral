resource "cyral_integration_logging" "cloud_watch_integration" {
  name = "my-cloudwatch-integration"
  cloud_watch {
    region = "us-east-1"
    group  = "everyone"
    stream = "abc"
  }
}

resource "cyral_integration_logging" "fluent_bit_integration" {
  name = "my-fluentbit-integration"
  fluent_bit {
    # Configuring a s3 bucket
    config = <<-EOF
    [OUTPUT]
      Name s3
      Match *
      Region us-east-2
      Bucket user-us-east-2
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

