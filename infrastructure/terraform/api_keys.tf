# HaiLanGo - API Keys Configuration (Simplified)
# Only creates API key for accessing GCP APIs

# API Key for accessing Vision, TTS, STT, Translate APIs
resource "google_apikeys_key" "hailango_web_key" {
  count = var.create_api_keys ? 1 : 0

  name         = "hailango-web-key"
  display_name = "HaiLanGo Web API Key"
  project      = var.project_id

  restrictions {
    # API restrictions - only allow specific APIs
    api_targets {
      service = "vision.googleapis.com"
    }
    api_targets {
      service = "texttospeech.googleapis.com"
    }
    api_targets {
      service = "speech.googleapis.com"
    }
    api_targets {
      service = "translate.googleapis.com"
    }

    # Browser key restrictions (add your domains in production)
    # browser_key_restrictions {
    #   allowed_referrers = [
    #     "localhost:3000",
    #     "*.hailango.com",
    #   ]
    # }
  }

  depends_on = [time_sleep.wait_for_apis]
}

# Output
output "web_api_key" {
  description = "Web API Key (restricted to Vision, TTS, STT, Translate)"
  value       = var.create_api_keys ? google_apikeys_key.hailango_web_key[0].key_string : null
  sensitive   = true
}
