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
							"mode": { "enum": ["inMemory", "memcached"] }
						  },
						  "memcachedOptions": {
							"type": "object",
							"properties": {
							  "enabled": {
								"type": "boolean"
							  },
							  "key": {
								"type": "string"
							  }
							}
						  }
						},
						"if": {
						  "properties": {
							"cacheOption": { "const": "memcached" }
						  },
						  "required": ["cacheOption"]
						},
						"then": {
						  "required": ["memcachedOptions"]
						},
						"else": {
						  "not": { "required": ["memcachedOptions"] }
						},
						"required": ["name", "cacheOption"]
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
						  "config": {
							"type": "string"
						  }
						},
						"required": ["alertManagers", "config"]
					  }
					},
					"additionalProperties": false,
					"required": ["name", "querier"]
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
					"if": {
					  "properties": { "mode": { "const": "receiver" } },
					  "required": ["mode"]
					},
					"then": {
					  "required": ["receiveReference"]
					},
					"else": {
					  "not": { "required": ["receiveReference"] }
					},
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
			  "type": {
				"type": "string"
			  },
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
				"required": ["bucket", "endpoint", "accessKey", "secretKey"]
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
