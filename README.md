
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-message-agent)
[![Actions Status](https://github.com/tommzn/hdb-message-agent/actions/workflows/go.image.build.amd64.yml/badge.svg)](https://github.com/tommzn/hdb-message-agent/actions)

# HomeDashboard Message Agent
Home Dashboard message agent, forwards messages from AWS SQS to Kafka.

## Config
By config you can define Kafka servers and a queue to topic forwarding. If no topic is specified given queue name is taken.
```yaml
kafka:
  server: localhost

agent:
  routes:
    - source: sqs.queue.1
      target: kafka.topic
    - source: data.distribution
  
```

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
