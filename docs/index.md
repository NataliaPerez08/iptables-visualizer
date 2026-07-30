# IPTables Visualizer — Documentación

## Índice

1. [Arquitectura](architecture.md)
2. [API Reference](api.md)
3. [Captura de Estado del Firewall](capture.md)
4. [Desarrollo](development.md)

## Descripción General

**IPTables Visualizer** es una herramienta web para gestionar, visualizar, validar, compilar y desplegar reglas de firewall en Linux. Soporta los motores **iptables** y **nftables**.

### Tecnologías

| Capa       | Tecnología                                                    |
|------------|---------------------------------------------------------------|
| Backend    | Go — Chi router, JWT, bcrypt                                 |
| BD         | SQLite3 (WAL mode)                                            |
| Frontend   | JavaScript vanilla, HTML5, CSS3 (sin framework)               |
| Firewall   | iptables / nftables vía `os/exec`                            |

### Inicio Rápido

```bash
cd backend
go build -o server .
export DRY_RUN_ONLY=true
./server
```

Abrir http://localhost:8080 — Credenciales: `admin` / `admin123`
