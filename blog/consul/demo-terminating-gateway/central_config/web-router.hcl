Kind = "service-router"
Name = "upstream_service"
Routes = [
  {
    Match {
      HTTP {
        PathPrefix = "/v1"
      }
    }

    Destination {
      ServiceSubset = "v1"
    }
  },
  {
    Match {
      HTTP {
        PathPrefix = "/v2"
      }
    }

    Destination {
      ServiceSubset = "v2"
    }
  },
]