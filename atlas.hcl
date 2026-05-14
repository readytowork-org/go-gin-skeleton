variable "DATABASE_URL" {
  type = string
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path",
    "./domain/models",
    "--dialect",
    "mysql", // | postgres | sqlite | sqlserver
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://mysql/8/dev"
  migration {
    dir = "file://database/migration"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "mysql" {
  url = var.DATABASE_URL
  src = "file://database/schema.sql"
  dev = "docker://mysql/8/dev"
  migration {
    dir = "file://database/migration"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
