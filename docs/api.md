# API Reference

Todas las rutas bajo `/api/v1`. Autenticación vía `Authorization: Bearer <token>`.

---

## Autenticación

### `POST /api/v1/auth/login`

**Body:**
```json
{ "username": "admin", "password": "admin123" }
```

**Respuesta:**
```json
{
  "token": "eyJ...",
  "user": { "id": 1, "username": "admin", "role": "admin", ... },
  "expires_at": 1700000000
}
```

### `GET /api/v1/auth/me`

Devuelve el usuario autenticado.

---

## Usuarios (solo admin)

### `POST /api/v1/users`

```json
{ "username": "editor1", "password": "pass123", "role": "editor", "email": "e@x.com" }
```

### `GET /api/v1/users`

Lista todos los usuarios.

---

## Políticas

### `GET /api/v1/policies`

Lista todas las políticas.

### `POST /api/v1/policies`

Crear una política con reglas.

### `GET /api/v1/policies/{id}`

Obtener una política por ID.

### `PUT /api/v1/policies/{id}`

Actualizar política (admin/editor).

### `DELETE /api/v1/policies/{id}`

Eliminar política (admin).

### `POST /api/v1/policies/{id}/validate`

Validar reglas de una política.

### `POST /api/v1/policies/{id}/dry-run?driver=iptables`

Previsualizar comandos sin aplicar.

### `POST /api/v1/policies/{id}/apply?driver=iptables`

Aplicar política al firewall (admin/editor).

---

## Sistema

### `GET /api/v1/system/firewall`

Captura el estado actual del firewall desde `iptables-save`.

**Respuesta:**
```json
{
  "tables": [
    {
      "name": "filter",
      "chains": {
        "INPUT": {
          "policy": "ACCEPT",
          "packets": 1234,
          "bytes": 567890,
          "rules": [
            {
              "chain": "INPUT",
              "target": "ACCEPT",
              "protocol": "tcp",
              "source": "10.0.0.0/8",
              "destination": "",
              "src_port": "",
              "dst_port": "80",
              "in_interface": "eth0",
              "out_interface": "",
              "state": "NEW",
              "log_prefix": "",
              "extra": [],
              "raw": "-A INPUT -s 10.0.0.0/8 -p tcp --dport 80 -j ACCEPT"
            }
          ]
        }
      }
    }
  ],
  "timestamp": "2026-07-30T12:00:00Z"
}
```

---

## Auditoría (solo admin)

### `GET /api/v1/audit?limit=100&action=apply_policy&from=2026-01-01&to=2026-07-30`

Parámetros: `user_id`, `action`, `resource`, `from`, `to`, `limit`, `offset`.
