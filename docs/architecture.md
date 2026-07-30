# Arquitectura

## Estructura del Proyecto

```
backend/
├── main.go                          # Punto de entrada, migraciones, SPA handler
├── internal/
│   ├── api/
│   │   ├── router.go                # Definición de rutas (Chi)
│   │   ├── handlers/                # Handlers HTTP
│   │   │   ├── auth_handler.go      # Login, users
│   │   │   ├── policy_handler.go    # CRUD de políticas, apply, dry-run
│   │   │   ├── audit_handler.go     # Consulta de logs de auditoría
│   │   │   └── system_handler.go    # Estado del sistema (firewall)
│   │   └── middleware/              # Auth, roles, logging
│   ├── capture/                     # Captura del estado vivo del firewall
│   │   └── capture.go               # Scanner iptables-save + parser
│   ├── config/                      # Config vía variables de entorno
│   ├── deployment/                  # Ejecutor de comandos firewall
│   ├── drivers/                     # Compiladores iptables / nftables
│   ├── engine/                      # Compilador y validador de políticas
│   ├── models/                      # Estructuras de datos
│   └── repository/                  # Capa de persistencia SQLite
└── web/                             # Frontend SPA embebido
    ├── index.html
    ├── css/
    └── js/
        ├── api.js                   # Cliente HTTP con JWT
        ├── app.js                   # Router SPA
        ├── components/              # Componentes reutilizables
        └── views/                   # Vistas (páginas)
```

## Flujo de Datos

```
Usuario → Chi Router → Middleware (Auth, Roles, Logging)
                            ↓
                     Handler HTTP
                        ↙     ↘
              Repository     Engine
              (SQLite)     (Compiler + Validator)
                              ↓
                          Drivers
                    (iptables / nftables)
                              ↓
                       Deployment
                      (os/exec real)
```

## Captura de Estado Actual

El paquete `capture` ejecuta `iptables-save -c` en el sistema y parsea la salida en estructuras Go. Esto permite leer la configuración actual del firewall en lugar de depender solo de la base de datos.

Flujo:

```
GET /v1/system/firewall
  → SystemHandler.GetFirewall()
    → capture.Capture()
      → exec("iptables-save -c")
      → parse(output)
      → FirewallState (JSON)
```

## Autenticación

JWT con claims: `user_id`, `username`, `role`. Roles: `admin`, `editor`, `viewer`.
