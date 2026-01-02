#Serviço de monitoramento simples em Linux (nível júnior/trainee DevOps)

#Estrutura Mínima do Projeto:;
linux-monitoring-agent/
├── agent/
│   └── main.go  (ou app.py)
├── systemd/
│   └── monitoring-agent.service
├── README.md
└── Makefile (opcional)
# Linux Monitoring Agent

Agente leve de monitoramento para Linux, escrito em Go, que expõe métricas via HTTP
em formato compatível com Prometheus.

## Funcionalidades
- Health check (`/health`)
- Métricas (`/metrics`)
- Uso real de CPU (via /proc/stat)
- Uso real de memória (via /proc/meminfo)
- Porta configurável por variável de ambiente
- Execução como serviço systemd

## Métricas expostas
- agent_uptime_seconds
- go_goroutines
- process_pid
- node_memory_total_kb
- node_memory_used_kb
- node_memory_available_kb
- node_cpu_usage_percent

## Execução local
```bash
go run agent/main.go

porta custumizada
PORT=9300 go run agent/main.go

#instalação como serviço
go build -o monitoring-agent agent/main.go
sudo cp monitoring-agent /usr/local/bin/
sudo cp monitoring-agent.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now monitoring-agent

