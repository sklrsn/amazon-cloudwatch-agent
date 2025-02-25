// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package config

import (
	"regexp"
	"strings"
)

// Keep a copy of schema.json in case we need to directly use it.

var schema = `{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "description": "Amazon CloudWatch Agent JSON Schema",
  "properties": {
    "agent": {
      "$ref": "#/definitions/agentDefinition"
    },
    "metrics": {
      "$ref": "#/definitions/metricsDefinition"
    },
    "logs": {
      "$ref": "#/definitions/logsDefinition"
    },
    "csm": {
      "$ref": "#/definitions/csmDefinition"
    }
  },
  "additionalProperties": true,
  "definitions": {
    "agentDefinition": {
      "type": "object",
      "description": "General configuration for Amazon CloudWatch Agent",
      "properties": {
        "metrics_collection_interval": {
          "description": "How often the metrics defined will be collected",
          "$ref": "#/definitions/timeIntervalDefinition"
        },
        "logfile": {
          "description": "Specifies the location to where the CloudWatch agent writes log messages. If you specify an empty string, the log goes to stdout",
          "type": "string",
          "maxLength": 4096
        },
        "region": {
          "description": "Specifies the region to use for the CloudWatch endpoint",
          "type": "string",
          "minLength": 1,
          "maxLength": 64
        },
        "debug": {
          "description": "Specifies running the CloudWatch agent with debug log messages",
          "type": "boolean"
        },
        "aws_sdk_log_level": {
          "description": "Specifies running the CloudWatch agent with AWS SDK debug logging. Multiple options must be separated by vertical bars.",
          "type": "string"
        },
        "credentials": {
          "description": "The credentials with which agent can access aws resources",
          "$ref": "#/definitions/credentialsDefinition"
        },
        "omit_hostname": {
          "description": "Hostname will be tagged by default unless you specifying append_dimensions, this flag allow you to omit hostname from tags without specifying append_dimensions",
          "type": "boolean"
        }
      },
      "additionalProperties": true
    },
    "csmDefinition": {
      "type": "object",
      "description": "Configuration for AWS SDK client-side monitoring",
      "properties": {
        "service_addresses": {
          "description": "List of addresses and ports to listen for UDP events on",
          "type": "array",
          "minItems": 1,
          "items": {
            "type": "string",
            "minLength": 1,
            "maxLength": 255
          }
        },
        "port": {
          "description": "Localhost UDP port to listen to client-side monitoring events on",
          "$ref": "#/definitions/userPortDefinition"
        },
        "log_level": {
          "description": "Determines the verbosity of logging with 0 being disabled.",
          "type": "integer"
        },
        "memory_limit_in_mb": {
          "description": "Approximate amount of memory, in MB, to use for unpublished monitoring records",
          "type": "integer"
        },
        "endpoint_override": {
          "description": "Override CSM service endpoint.",
          "type": "string",
          "format": "uri"
        }
      },
      "additionalProperties": false
    },
    "metricsDefinition": {
      "type": "object",
      "description": "configuration for metrics to be collected",
      "properties": {
        "namespace": {
          "type": "string",
          "description": "The namespace to use for the metrics collected by the agent. The default is CWAgent",
          "minLength": 1,
          "maxLength": 255
        },
        "aggregation_dimensions": {
          "description": "Specifies the dimensions on which collected metrics are to be aggregated",
          "type": "array",
          "items": {
            "type": "array",
            "items": {
              "type": "string",
              "minLength": 1,
              "maxLength": 1024
            },
            "uniqueItems": true,
            "maxItems": 30
          },
          "uniqueItems": true,
          "minItems": 1,
          "maxItems": 1024
        },
        "append_dimensions": {
          "type": "object",
          "description": "Adds Amazon EC2 metric dimensions to all metrics collected by the agent, we only support fixed key value pair now: ImageId:{aws:ImageId},InstanceId:{aws:InstanceId},InstanceType:{aws:InstanceType},AutoScalingGroupName:{aws:AutoScalingGroupName}. ",
          "maxProperties": 30,
          "additionalProperties": {
            "type": "string",
            "minLength": 1,
            "maxLength": 1024
          }
        },
        "metrics_collected": {
          "type": "object",
          "properties": {
            "collectd": {
              "$ref": "#/definitions/metricsDefinition/definitions/collectdDefinitions"
            },
            "cpu": {
              "$ref": "#/definitions/metricsDefinition/definitions/cpuDefinitions"
            },
            "disk": {
              "$ref": "#/definitions/metricsDefinition/definitions/diskDefinitions"
            },
            "diskio": {
              "$ref": "#/definitions/metricsDefinition/definitions/diskioDefinitions"
            },
            "statsd": {
              "$ref": "#/definitions/metricsDefinition/definitions/statsdDefinitions"
            },
            "swap": {
              "$ref": "#/definitions/metricsDefinition/definitions/swapDefinitions"
            },
            "mem": {
              "$ref": "#/definitions/metricsDefinition/definitions/memDefinitions"
            },
            "net": {
              "$ref": "#/definitions/metricsDefinition/definitions/netDefinitions"
            },
            "netstat": {
              "$ref": "#/definitions/metricsDefinition/definitions/netstatDefinitions"
            },
            "processes": {
              "$ref": "#/definitions/metricsDefinition/definitions/processesDefinitions"
            },
            "procstat": {
              "$ref": "#/definitions/metricsDefinition/definitions/procstatDefinitions"
            },
            "ethtool": {
              "$ref": "#/definitions/metricsDefinition/definitions/ethtoolDefinitions"
            },
            "nvidia_smi": {
              "$ref": "#/definitions/metricsDefinition/definitions/nvidiaGpuDefinitions"
            }
          },
          "minProperties": 1,
          "additionalProperties": {
            "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
          }
        },
        "force_flush_interval": {
          "description": "Max time to wait before batch publishing the metrics, unit is second.",
          "$ref": "#/definitions/timeIntervalDefinition"
        },
        "credentials": {
          "description": "The credentials with which agent can access aws resources",
          "$ref": "#/definitions/credentialsDefinition"
        },
        "endpoint_override": {
          "description": "The override endpoint to use to access cloudwatch",
          "$ref": "#/definitions/endpointOverrideDefinition"
        }
      },
      "additionalProperties": false,
      "required": [
        "metrics_collected"
      ],
      "definitions": {
        "basicMetricDefinition": {
          "type": "object",
          "properties": {
            "metrics_collection_interval": {
              "$ref": "#/definitions/timeIntervalDefinition"
            },
            "append_dimensions": {
              "$ref": "#/definitions/generalAppendDimensionsDefinition"
            },
            "measurement": {
              "$ref": "#/definitions/metricsDefinition/definitions/metricsMeasurementDefinition"
            }
          },
          "required": [
            "measurement"
          ]
        },
        "basicResourcesDefinition": {
          "type": "object",
          "properties": {
            "resources": {
              "type": "array",
              "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 4096
              },
              "maxItems": 256
            }
          }
        },
        "collectdDefinitions": {
          "type": "object",
          "properties": {
            "service_address": {
              "type": "string",
              "minLength": 1,
              "maxLength": 255
            },
            "name_prefix": {
              "type": "string",
              "minLength": 1,
              "maxLength": 255
            },
            "collectd_auth_file": {
              "type": "string",
              "minLength": 1,
              "maxLength": 4096
            },
            "collectd_security_level": {
              "type": "string",
              "enum": [
                "none",
                "sign",
                "encrypt"
              ]
            },
            "collectd_typesdb": {
              "type": "array",
              "maxItems": 10,
              "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 4096
              }
            },
            "metrics_aggregation_interval": {
              "$ref": "#/definitions/timeIntervalWithZeroDefinition"
            }
          },
          "additionalProperties": false
        },
        "cpuDefinitions": {
          "type": "object",
          "allOf": [
            {
              "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
            },
            {
              "type": "object",
              "properties": {
                "resources": {
                  "type": "array",
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "enum": [
                      "*"
                    ]
                  }
                },
                "totalcpu": {
                  "type": "boolean"
                }
              }
            }
          ]
        },
        "diskDefinitions": {
          "type": "object",
          "allOf": [
            {
              "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
            },
            {
              "$ref": "#/definitions/metricsDefinition/definitions/basicResourcesDefinition"
            },
            {
              "type": "object",
              "properties": {
                "ignore_file_system_types": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 4096
                  },
                  "maxItems": 256
                },
                "drop_device": {
                  "type": "boolean"
                }
              }
            }
          ]
        },
        "diskioDefinitions": {
          "type": "object",
          "allOf": [
            {
              "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
            },
            {
              "$ref": "#/definitions/metricsDefinition/definitions/basicResourcesDefinition"
            }
          ]
        },
        "statsdDefinitions": {
          "type": "object",
          "properties": {
            "allowed_pending_messages": {
              "type": "integer",
              "minimum": 1,
              "maximum": 2147483647
            },
            "service_address": {
              "type": "string",
              "minLength": 1,
              "maxLength": 255
            },
            "metrics_collection_interval": {
              "$ref": "#/definitions/timeIntervalDefinition"
            },
            "metrics_aggregation_interval": {
              "$ref": "#/definitions/timeIntervalWithZeroDefinition"
            },
            "metric_separator": {
              "type": "string",
              "minLength": 1,
              "maxLength": 255
            }
          },
          "additionalProperties": false
        },
        "swapDefinitions": {
          "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
        },
        "memDefinitions": {
          "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
        },
        "netDefinitions": {
          "type": "object",
          "allOf": [
            {
              "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
            },
            {
              "$ref": "#/definitions/metricsDefinition/definitions/basicResourcesDefinition"
            }
          ]
        },
        "netstatDefinitions": {
          "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
        },
        "processesDefinitions": {
          "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
        },
        "procstatDefinitions": {
          "type": "array",
          "minItems": 1,
          "maxItems": 255,
          "items": {
            "allOf": [
              {
                "$ref": "#/definitions/metricsDefinition/definitions/basicMetricDefinition"
              },
              {
                "type": "object",
                "properties": {
                  "pid_file": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 255,
                    "descriptions": "the path of pid_file"
                  },
                  "exe": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 255,
                    "descriptions": "a regex matches the names of processes"
                  },
                  "pattern": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 255,
                    "descriptions": "a regex matches the whole command of processes"
                  },
                  "measurement": {
                    "$ref": "#/definitions/metricsDefinition/definitions/metricsMeasurementWithoutDecorationDefinition"
                  }
                },
                "anyOf": [
                  {
                    "required": [
                      "pid_file"
                    ]
                  },
                  {
                    "required": [
                      "exe"
                    ]
                  },
                  {
                    "required": [
                      "pattern"
                    ]
                  }
                ]
              }
            ]
          }
        },
        "ethtoolDefinitions": {
          "type": "object",
          "properties": {
            "interface_include": {
              "type": "array",
              "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 255
              }
            },
            "interface_exclude": {
              "type": "array",
              "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 255
              }
            },
            "metrics_include": {
              "type": "array",
              "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 255
              }
            }
          },
          "additionalProperties": false
        },
        "nvidiaGpuDefinitions": {
          "type": "object",
          "properties": {
            "measurement": {
              "type": "array",
              "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 255
              }
            }
          },
          "metrics_collection_interval": {
            "$ref": "#/definitions/timeIntervalDefinition"
          }
        },
        "metricsMeasurementWithoutDecorationDefinition": {
          "type": "array",
          "items": {
            "type": "string",
            "minLength": 1,
            "maxLength": 255
          },
          "uniqueItems": true
        },
        "metricsMeasurementDefinition": {
          "type": "array",
          "items": {
            "oneOf": [
              {
                "type": "string",
                "minLength": 1,
                "maxLength": 255
              },
              {
                "type": "object",
                "properties": {
                  "name": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 255
                  },
                  "rename": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 255
                  },
                  "unit": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 256
                  }
                }
              }
            ]
          },
          "uniqueItems": true
        }
      }
    },
    "logsDefinition": {
      "type": "object",
      "descriptions": "configuration for collecting logs and upload to cloudWatch log service",
      "properties": {
        "logs_collected": {
          "type": "object",
          "properties": {
            "files": {
              "$ref": "#/definitions/logsDefinition/definitions/logsFilesDefinition"
            },
            "windows_events": {
              "$ref": "#/definitions/logsDefinition/definitions/logsWindowsEventsDefinition"
            }
          },
          "minProperties": 1,
          "additionalProperties": false
        },
        "metrics_collected": {
          "type": "object",
          "properties": {
            "ecs": {
              "type": "object",
              "properties": {
                "metrics_collection_interval": {
                  "$ref": "#/definitions/timeIntervalDefinition"
                }
              },
              "additionalProperties": false
            },
            "kubernetes": {
              "type": "object",
              "properties": {
                "cluster_name": {
                  "type": "string",
                  "minLength": 1,
                  "maxLength": 512
                },
                "metrics_collection_interval": {
                  "$ref": "#/definitions/timeIntervalDefinition"
                }
              },
              "additionalProperties": true
            },
            "prometheus":{
              "type": "object",
              "properties": {
                "cluster_name": {
                  "type": "string"
                },
                "log_group_name": {
                  "type": "string"
                },
                "prometheus_config_path": {
                  "type": "string"
                },
                "emf_processor": {
                  "$ref": "#/definitions/emfProcessorDefinition"
                },
                "ecs_service_discovery": {
                  "$ref": "#/definitions/ecsServiceDiscoveryDefinition"
                }
              },
              "additionalProperties": false
            }
          },
          "additionalProperties": true
        },
        "log_stream_name": {
          "$ref": "#/definitions/logsDefinition/definitions/logStreamNameDefinition"
        },
        "force_flush_interval": {
          "description": "Max time to wait before batch publishing the log, unit is second.",
          "$ref": "#/definitions/timeIntervalDefinition"
        },
        "credentials": {
          "description": "The credentials with which agent can access aws resources",
          "$ref": "#/definitions/credentialsDefinition"
        },
        "endpoint_override": {
          "description": "The override endpoint to use to access cloudwatch logs",
          "$ref": "#/definitions/endpointOverrideDefinition"
        }
      },
      "additionalProperties": false,
      "anyOf": [
        {
          "required": [
            "logs_collected"
          ]
        },
        {
          "required": [
            "metrics_collected"
          ]
        }
      ],
      "definitions": {
        "logsFilesDefinition": {
          "type": "object",
          "descriptions": "Specifies the log files to be collected",
          "properties": {
            "collect_list": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "file_path": {
                    "type": "string",
                    "maxLength": 4096
                  },
                  "log_group_name": {
                    "$ref": "#/definitions/logsDefinition/definitions/logGroupNameDefinition"
                  },
                  "log_stream_name": {
                    "$ref": "#/definitions/logsDefinition/definitions/logStreamNameDefinition"
                  },
                  "multi_line_start_pattern": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 4096
                  },
                  "timestamp_format": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 4096
                  },
                  "timezone": {
                    "type": "string",
                    "enum": [
                      "Local",
                      "LOCAL",
                      "UTC"
                    ]
                  },
                  "encoding": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 4096
                  },
                  "auto_removal": {
                    "type": "boolean"
                  },
                  "blacklist": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 4096
                  },
                  "publish_multi_logs": {
                    "type": "boolean"
                  },
                  "retention_in_days": {
                    "$ref": "#/definitions/logsDefinition/definitions/retentionInDaysDefinition"
                  },
                  "filters": {
                    "type": "array",
                    "items": {
                      "$ref": "#/definitions/logsDefinition/definitions/filterDefinition"
                    }
                  }
                },
                "required": [
                  "file_path"
                ],
                "additionalProperties": false
              },
              "minItems": 1,
              "maxItems": 16384,
              "uniqueItems": true
            }
          },
          "required": [
            "collect_list"
          ],
          "additionalProperties": false
        },
        "logsWindowsEventsDefinition": {
          "type": "object",
          "descriptions": "Specifies the logs to collect from servers running Windows Server",
          "properties": {
            "collect_list": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "event_name": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 255,
                    "not": {
                      "type": "string",
                      "enum": [
                        "Forwarded Events"
                      ]
                    }
                  },
                  "event_levels": {
                    "type": "array",
                    "items": {
                      "type": "string",
                      "enum": [
                        "CRITICAL",
                        "ERROR",
                        "WARNING",
                        "INFORMATION",
                        "VERBOSE"
                      ],
                      "minItems": 1,
                      "uniqueItems": true
                    }
                  },
                  "log_stream_name": {
                    "$ref": "#/definitions/logsDefinition/definitions/logStreamNameDefinition"
                  },
                  "log_group_name": {
                    "$ref": "#/definitions/logsDefinition/definitions/logGroupNameDefinition"
                  },
                  "retention_in_days": {
                    "$ref": "#/definitions/logsDefinition/definitions/retentionInDaysDefinition"
                  },
                  "event_format": {
                    "type": "string",
                    "enum": [
                      "text",
                      "xml"
                    ]
                  }
                },
                "required": [
                  "event_name",
                  "event_levels"
                ],
                "additionalProperties": false
              },
              "minItems": 1,
              "maxItems": 16384,
              "uniqueItems": true
            }
          },
          "additionalProperties": false,
          "required": [
            "collect_list"
          ]
        },
        "logGroupNameDefinition": {
          "type": "string",
          "minLength": 1,
          "maxLength": 512
        },
        "logStreamNameDefinition": {
          "type": "string",
          "minLength": 1,
          "maxLength": 512
        },
        "retentionInDaysDefinition": {
          "type": "integer",
          "enum": [
           -1,
            1,
            3,
            5,
            7,
            14,
            30,
            60,
            90,
            120,
            150,
            180,
            365,
            400,
            545,
            731,
            1827,
            3653
          ]
        },
        "filterDefinition": {
          "type": "object",
          "descriptions": "Define filters to apply to the log messages in this log file to determine whether to publish the message or not",
          "additionalProperties": false,
          "properties": {
            "type": {
              "description": "Declares if the specified filter should be used to include or exclude log messages",
              "type": "string",
              "enum": [
                "include",
                "exclude"
              ]
            },
            "expression": {
              "description": "Regular expression to apply to the log message",
              "type": "string"
            }
          }
        }
      }
    },
    "timeIntervalDefinition": {
      "type": "integer",
      "minimum": 1,
      "maximum": 172800
    },
    "timeIntervalWithZeroDefinition": {
      "type": "integer",
      "minimum": 0,
      "maximum": 172800
    },
    "userPortDefinition": {
      "type": "integer",
      "minimum": 1024,
      "maximum": 65535
    },
    "generalAppendDimensionsDefinition": {
      "descriptions": "Additional customized dimensions to use",
      "type": "object",
      "maxProperties": 30,
      "additionalProperties": {
        "type": "string",
        "minLength": 1,
        "maxLength": 1024
      }
    },
    "credentialsDefinition": {
      "type": "object",
      "properties": {
        "role_arn": {
          "description": "The target IAM role with which agent can access aws resources",
          "type": "string",
          "minLength": 20,
          "maxLength": 2048
        }
      },
      "additionalProperties": false
    },
    "endpointOverrideDefinition": {
      "type": "string",
      "minLength": 4,
      "maxLength": 2048
    },
    "ecsServiceDiscoveryDefinition": {
      "type": "object",
      "descriptions": "Define ECS service discovery for Prometheus",
      "properties": {
        "docker_label": {
          "$ref": "#/definitions/ecsServiceDiscoveryDefinition/definitions/dockerLabel"
        },
        "task_definition_list": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ecsServiceDiscoveryDefinition/definitions/taskDefinitionList"
          }
        },
        "service_name_list_for_tasks": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ecsServiceDiscoveryDefinition/definitions/serviceNameListForTasks"
          }
        },
        "sd_cluster_region": {
          "description": "ECS cluster region",
          "type": "string"
        },
        "sd_frequency": {
          "description": "ECS service discovery frequency",
          "type": "string"
        },
        "sd_result_file": {
          "description": "ECS service discovery result file full path",
          "type": "string"
        },
        "sd_target_cluster": {
          "description": "The target ECS cluster to be scanned for Prometheus exporters",
          "type": "string"
        }
      },
      "additionalProperties": false,
      "definitions": {
        "dockerLabel": {
          "type": "object",
          "descriptions": "Define ECS service discovery based on docker labels",
          "properties": {
            "sd_job_name_label": {
              "description": "Docker label name for specifying ECS service discovery job name",
              "type": "string"
            },
            "sd_metrics_path_label": {
              "description": "Docker label name for specifying the Prometheus resource path",
              "type": "string"
            },
            "sd_port_label": {
              "description": "Docker label name for specifying the Prometheus port",
              "type": "string"
            }
          }
        },
        "taskDefinitionList": {
          "type": "object",
          "descriptions": "Define ECS service discovery based on task definitions",
          "properties": {
            "sd_container_name_pattern": {
              "description": "ECS container name pattern which expose the Prometheus metrics",
              "type": "string"
            },
            "sd_job_name": {
              "description": "Service discovery result job name",
              "type": "string"
            },
            "sd_metrics_path": {
              "description": "Prometheus metrics path of the exporters",
              "type": "string"
            },
            "sd_metrics_ports": {
              "description": "Prometheus metrics port list of the exporters",
              "type": "string"
            },
            "sd_task_definition_arn_pattern": {
              "description": "ECS task definition pattern which expose the Prometheus metrics",
              "type": "string"
            }
          }
        },
        "serviceNameListForTasks": {
          "type": "object",
          "descriptions": "Define ECS service discovery based on service names",
          "properties": {
            "sd_container_name_pattern": {
              "description": "ECS container name pattern which expose the Prometheus metrics",
              "type": "string"
            },
            "sd_job_name": {
              "description": "Service discovery result job name",
              "type": "string"
            },
            "sd_metrics_path": {
              "description": "Prometheus metrics path of the exporters",
              "type": "string"
            },
            "sd_metrics_ports": {
              "description": "Prometheus metrics port list of the exporters",
              "type": "string"
            },
            "sd_service_name_pattern":{
              "description": "ECS service name pattern responsible for tasks which expose the Prometheus metrics",
              "type": "string"
            }
          }
        }
      }
    },
    "emfProcessorDefinition": {
      "type": "object",
      "descriptions": "Define EMF Processor to set metric filter",
      "properties": {
        "metric_declaration_dedup": {
          "description": "Enable the de-duplication function for the EMF metric",
          "type": "boolean"
        },
        "metric_namespace": {
          "description": "The namespace to use for the Prometheus metrics collected by the agent",
          "type": "string"
        },
        "metric_unit": {
          "description": "The metric name, metric unit map",
          "type": "object",
          "additionalProperties": {
            "type": "string",
            "minLength": 1,
            "maxLength": 256
          }
        },
        "metric_declaration": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/emfProcessorDefinition/definitions/metricDeclarationDefinition"
          }
        }
      },
      "additionalProperties": false,
      "definitions": {
        "metricDeclarationDefinition": {
          "type": "object",
          "descriptions": "Define metric declaration to set EMF",
          "properties": {
            "source_labels": {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            "label_matcher": {
              "type": "string"
            },
            "label_separator": {
              "type": "string"
            },
            "metric_selectors": {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            "dimensions": {
              "type": "array",
              "items": {
                "type": "array",
                "items": {
                  "type": "string"
                }
              }
            }
          },
          "additionalProperties": false
        }
      }
    }
  }
}
`

func GetJsonSchema() string {
	return schema
}

func OverwriteSchema(newSchema string) {
	schema = newSchema
}

// Translate Sample:
// (root).agent.metrics_collection_interval -> /agent/metrics_collection_interval
// (root).metrics.metrics_collected.cpu.resources.1 -> /metrics/metrics_collected/cpu/resources/1
func GetFormattedPath(rawPath string) string {
	//replace heading (root). to /
	prefixRe := regexp.MustCompile("^\\(root\\).")
	result := prefixRe.ReplaceAllString(rawPath, "/")
	//replace . to /
	result = strings.Replace(result, ".", "/", -1)
	return result
}
