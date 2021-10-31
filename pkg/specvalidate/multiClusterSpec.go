package specvalidate

var RuleSchema = `
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
	  "clusters": {
		"type": "array",
		"items": {
		  "type": "object",
		  "properties": {
			"name": {
			  "type": "string"
			},
			"type": { "enum": ["thanos", "grafana", "prometheus"] }
		  },
		  "required": ["name", "type"],
		  "allOf": [
			{
			  "if": {
				"properties": { "type": { "const": "thanos" } },
				"required": ["type"]
			  },
			  "then": {
				"properties": {
				  "data": {
					"type": "object",
					"properties": {
					  "name": {
						"type": "string"
					  },
					  "namespace": {
						"type": "string"
					  },
					  "install": {
						"type": "boolean"
					  },
					  "objStoreConfig": {
						"type": "string"
					  },
					  "querier": {
						"type": "object",
						"properties": {
						  "name": {
							"type": "string"
						  },
						  "targets": {
							"type": "array",
							"items": [
							  {
								"type": "string"
							  }
							]
						  },
						  "dedupEnbaled": {
							"type": "boolean"
						  },
						  "autoDownSample": {
							"type": "boolean"
						  },
						  "partialResponse": {
							"type": "boolean"
						  }
						},
						"required": ["name"]
					  },
					  "querierFE": {
						"type": "object",
						"properties": {
						  "name": {
							"type": "string"
						  },
						  "cacheOption": {
							"enum": ["inMemory", "memcached"],
							"type": "string"
						  },
						  "config": {
							"type": "object"
						  }
						},
						"required": ["name", "cacheOption", "config"]
					  },
					  "receiver": {
						"type": "object",
						"properties": {
						  "name": {
							"type": "string"
						  }
						},
						"required": ["name"]
					  },
					  "compactor": {
						"type": "object",
						"properties": {
						  "name": {
							"type": "string"
						  }
						},
						"required": ["name"]
					  },
					  "ruler": {
						"type": "object",
						"properties": {
						  "alertManagers": {
							"type": "array",
							"items": [
							  {
								"type": "string"
							  }
							]
						  },
						  "name": {
							"type": "string"
						  },
						  "config": {
							"type": "string"
						  }
						},
						"required": ["alertManagers", "config", "name"]
					  }
					},
					"additionalProperties": false,
					"required": ["name", "install", "namespace"]
				  }
				}
			  }
			},
			{
			  "if": {
				"properties": { "type": { "const": "grafana" } },
				"required": ["type"]
			  },
			  "then": {
				"properties": {
				  "data": {
					"type": "object",
					"properties": {
					  "name": {
						"type": "string"
					  },
					  "setup": {
						"type": "object",
						"properties": {
						  "enabled": {
							"type": "boolean"
						  },
						  "name": {
							"type": "string"
						  },
						  "namespace": {
							"type": "string"
						  }
						},
						"required": ["name", "namespace"]
					  }
					},
					"additionalProperties": false,
					"required": ["name", "setup"]
				  }
				}
			  }
			},
			{
			  "if": {
				"properties": { "type": { "const": "prometheus" } },
				"required": ["type"]
			  },
			  "then": {
				"properties": {
				  "data": {
					"type": "object",
					"properties": {
					  "name": {
						"type": "string"
					  },
					  "install": {
						"type": "boolean"
					  },
					  "namespace": {
						"type": "string"
					  },
					  "mode": { "enum": ["sidecar", "receiver"] },
					  "receiveReference": {
						"type": "string"
					  },
					  "objStoreConfig": {
						"type": "string"
					  }
					},
					"additionalProperties": false,
					"required": ["name", "namespace", "mode", "objStoreConfig"]
				  }
				}
			  }
			}
		  ]
		}
	  },
	  "objStoreConfigslist": {
		"type": "array",
		"items": [
		  {
			"type": "object",
			"properties": {
			  "name": {
				"type": "string"
			  },
			  "type": { "enum": ["S3", "GCS", "AZURE"] },
			  "config": {
				"type": "object",
				"properties": {
				  "bucket": {
					"type": "string"
				  },
				  "endpoint": {
					"type": "string"
				  },
				  "accessKey": {
					"type": "string"
				  },
				  "secretKey": {
					"type": "string"
				  }
				},
				"required": ["bucket"]
			  },
			  "bucketweb": {
				"type": "object",
				"properties": {
				  "enabled": {
					"type": "boolean"
				  }
				}
			  }
			},
			"required": ["name", "type", "config"]
		  }
		]
	  }
	},
	"required": ["clusters", "objStoreConfigslist"],
	"additionalProperties": false
  }   
`
