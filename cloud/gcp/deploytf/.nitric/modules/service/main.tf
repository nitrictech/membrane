terraform {
  required_providers {
    docker = {
      source = "kreuzwerker/docker"
    }
  }
}

# Create a GCR repository for the service image
data "google_container_registry_repository" "repo" {
  project = var.project_id
}

locals {
  service_image_url = "${var.artifact_registry_repository}/${var.service_name}"
}

# Tag the provided docker image with the repository url
resource "docker_tag" "tag" {
  source_image = var.image
  target_image = local.service_image_url
}

# Push the tagged image to the repository
resource "docker_registry_image" "push" {
  name = local.service_image_url
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
}

locals {
  ids_prefix = "nitric-"
}

resource "random_string" "unique_id" {
  length = 4
  special = false
  upper   = false
}

# Create a random ID for the service name, so that it confirms to regex restrictions
resource "random_string" "service_account_id" {
  length  = 30 - length(local.ids_prefix)
  special = false
  upper   = false
}

# Create a service account for the google cloud run instance
resource "google_service_account" "service_account" {
  account_id   = "${local.ids_prefix}${random_string.service_account_id.id}"
  project      = var.project_id
  display_name = "${var.service_name} service account"
  description  = "Service account which runs the ${var.service_name} service"
}

# Create a random password for events that will target this service
resource "random_password" "event_token" {
  length  = 32
  special = false
  keepers = {
    "name" = var.service_name
  }
}

locals{
  service_name = replace(var.service_name, "_", "-")
}

# Create a cloud run service
resource "google_cloud_run_v2_service" "service" {
  name = "${local.service_name}-${random_string.unique_id.result}"

  location = var.region
  project  = var.project_id
  # set launch_stage to BETA if gpus set otherwise GA
  # launch_stage        = var.gpus > 0 ? "BETA" : "GA"
  launch_stage        = "GA"
  deletion_protection = false

  ingress = var.internal_ingress == true ? "INGRESS_TRAFFIC_INTERNAL_ONLY" : "INGRESS_TRAFFIC_ALL"

  template {
    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }

    dynamic "vpc_access" {
      for_each = var.vpc != null ? [1] : []

     
      content {
        egress = var.vpc.all_traffic ? "ALL_TRAFFIC" : "PRIVATE_RANGES_ONLY"
        network_interfaces {
          network    = var.vpc.network
          subnetwork = var.vpc.subnet
          tags       = var.vpc.network_tags
        }
      }
    }

    encryption_key = var.kms_key != "" ? var.kms_key : null

    # dynamic "node_selector" {
    #   for_each = var.gpus > 0 ? [1] : []
    #   content {
    #     accelerator = "nvidia-l4"
    #   }
    # }
    containers {
      image = "${local.service_image_url}@${docker_registry_image.push.sha256_digest}"
      resources {
        limits = {
          cpu    = var.cpus
          memory = "${var.memory_mb}Mi"
        }

        # limits = merge({
        #   cpu    = "${var.cpus}"
        #   memory = "${var.memory_mb}Mi"
        # }, var.gpus > 0 ? { "nvidia.com/gpu" = var.gpus } : {})
      }

      ports {
        container_port = 9001
      }
      env {
        name  = "EVENT_TOKEN"
        value = random_password.event_token.result
      }
      env {
        name  = "SERVICE_ACCOUNT_EMAIL"
        value = google_service_account.service_account.email
      }
      env {
        name  = "GCP_REGION"
        value = var.region
      }

      dynamic "env" {
        for_each = var.environment
        content {
          name  = env.key
          value = env.value
        }
      }
    }

    service_account = google_service_account.service_account.email
    timeout         = "${var.timeout_seconds}s"
  }

  depends_on = [
    docker_registry_image.push,
    google_service_account_iam_member.account_member,
    google_service_account_iam_member.service_account_iam_member,
    google_service_account_iam_member.service_account_invoker_iam_member
  ]
}

# Create a random ID for the service name, so that it confirms to regex restrictions
resource "random_string" "service_id" {
  length  = 30 - length(local.ids_prefix)
  special = false
  upper   = false
}


data "google_client_openid_userinfo" "deployer" {
}

locals {
  deployer_email = data.google_client_openid_userinfo.deployer.email
  deployer_type  = endswith(local.deployer_email, "gserviceaccount.com") ? "serviceAccount" : "user"
}

# If we're impersonation a service account, we need to grant that account the service account user role on the service account
resource "google_service_account_iam_member" "service_account_iam_member" {
  service_account_id = google_service_account.service_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "${local.deployer_type}:${local.deployer_email}"
}

# If we're impersonation a service account, we need to grant that account the service account user role on the service account
resource "google_service_account_iam_member" "service_account_invoker_iam_member" {
  service_account_id = google_service_account.invoker_service_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "${local.deployer_type}:${local.deployer_email}"
}

# Create an invoker service account for the google cloud run instance
resource "google_service_account" "invoker_service_account" {
  project      = var.project_id
  account_id   = "${local.ids_prefix}${random_string.service_id.id}"
  display_name = "${var.service_name} invoker"
  description  = "Service account which allows other resources to invoke the ${var.service_name} service"
}

# Give the above service account permissions to execute the CloudRun service
resource "google_cloud_run_service_iam_member" "invoker" {
  service  = google_cloud_run_v2_service.service.name
  location = google_cloud_run_v2_service.service.location
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.invoker_service_account.email}"
}

resource "google_project_iam_member" "project_member" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.service_account.email}"
  # p.BaseComputeRole.Name,
  role = var.base_compute_role
}

# Give the above service account permissions to act as itself
resource "google_service_account_iam_member" "account_member" {
  service_account_id = google_service_account.invoker_service_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.invoker_service_account.email}"
}
