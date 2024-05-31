# Configures `my-sidecar-fluent-bit` to push logs to a bucket named
# `example-bucket` in AWS S3.
resource "cyral_sidecar" "sidecar_fluent_bit" {
  name                        = "my-sidecar-fluent-bit"
  deployment_method           = "terraform"
  activity_log_integration_id = cyral_integration_logging.s3.id
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

# Configures a raw Splunk integration with no sidecar associated.
resource "cyral_integration_logging" "splunk_integration" {
  name = "my-splunk-integration"
  splunk {
    hostname = "http://splunk.com"
    hec_port = "8088"
    access_token = "XXXXXXXXXXX"
  }
}
