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
			"thanos": {
			  "type": "object",
			  "properties": {
				"name": {
				  "type": "string"
				},
				"querier": {
				  "type": "object",
				  "properties": {
					"name": {
					  "type": "string"
					},
					"targets": {
					  "type": "string"
					},
					"dedup-enbaled": {
					  "type": "string"
					},
					"autoDownSample": {
					  "type": "string"
					},
					"partial_response": {
					  "type": "string"
					}
				  },
				  "required": ["name"]
				},
				"querier-fe": {
				  "type": "object",
				  "properties": {
					"name": {
					  "type": "string"
					},
					"cacheOption": {
					  "mode": { "enum": ["in-memory", "memcached"] }
					},
					"memcached-options": {
					  "type": "object",
					  "properties": {
						"enabled": {
						  "type": "boolean"
						},
						"key1": {
						  "type": "string"
						}
					  },
					  "required": ["key1"]
					}
				  },
				  "if": {
					"properties": { "cacheOption": { "const": "memcached" } },
					"required": ["cacheOption"]
				  },
				  "then": {
					"required": ["memcached-options"]
				  },
				  "else": {
					"not": { "required": ["memcached-options"] }
				  },
				  "required": ["name", "cacheOption"]
				},
				"receiver": {
				  "type": "object",
				  "properties": {
					"name": {
					  "type": "string"
					},
					"httpPort": {
					  "type": "string"
					},
					"httpNodePort": {
					  "type": "string"
					},
					"remoteWritePort": {
					  "type": "string"
					},
					"remoteWriteNodePort": {
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
					"alertmanagers": {
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
				  "required": ["alertmanagers", "config"]
				}
			  },
			  "required": ["name", "querier"]
			},
			"grafana": {
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
			  "required": ["name", "setup"]
			},
			"prometheus": {
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
		  },
		  "required": ["name"],
		  "anyOf": [
			{
			  "required": ["thanos"]
			},
			{
			  "required": ["prometheus"]
			},
			{
			  "required": ["grafana"]
			}
		  ]
		}
	  },
	  "s3configslist": {
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
				  "access_key": {
					"type": "string"
				  },
				  "secret_key": {
					"type": "string"
				  }
				},
				"required": ["bucket", "endpoint", "access_key", "secret_key"]
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
	"required": ["clusters"],
	"additionalProperties": true
  }`
