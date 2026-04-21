# basic resource with AZURE_STORAGE_QUEUE
resource "snowflake_notification_integration" "azure" {
  name    = "notification"
  enabled = true

  notification_provider           = "AZURE_STORAGE_QUEUE"
  azure_storage_queue_primary_uri = "https://myaccount.queue.core.windows.net/myqueue"
  azure_tenant_id                 = "a123b4c5-1234-123a-a12b-1a23b45678c9"
}

# basic resource with AWS_SNS
resource "snowflake_notification_integration" "aws" {
  name    = "notification"
  enabled = true

  notification_provider = "AWS_SNS"
  aws_sns_topic_arn     = "arn:aws:sns:us-east-1:001234567890:mytopic"
  aws_sns_role_arn      = "arn:aws:iam::001234567890:role/myrole"
}

# basic resource with GCP_PUBSUB (subscription)
resource "snowflake_notification_integration" "gcp_subscription" {
  name    = "notification"
  enabled = true

  notification_provider        = "GCP_PUBSUB"
  gcp_pubsub_subscription_name = "projects/myproject/subscriptions/mysubscription"
}

# basic resource with GCP_PUBSUB (topic)
resource "snowflake_notification_integration" "gcp_topic" {
  name    = "notification"
  enabled = true

  notification_provider = "GCP_PUBSUB"
  gcp_pubsub_topic_name = "projects/myproject/topics/mytopic"
}

# resource with all non-provider-specific fields set
resource "snowflake_notification_integration" "complete" {
  name    = "notification"
  enabled = true
  comment = "A notification integration."

  notification_provider           = "AZURE_STORAGE_QUEUE"
  azure_storage_queue_primary_uri = "https://myaccount.queue.core.windows.net/myqueue"
  azure_tenant_id                 = "a123b4c5-1234-123a-a12b-1a23b45678c9"
}
