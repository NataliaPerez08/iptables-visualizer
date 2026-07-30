# Desarrollo

## Requisitos

- Go 1.21+
- Linux con `iptables` (opcional para dry-run)
- `gcc` (requerido por `mattn/go-sqlite3` — CGO)

## Compilar

```bash
cd backend
go build -o server .
```

## Ejecutar

```bash
# Modo dry-run (seguro, sin cambios reales)
export DRY_RUN_ONLY=true
./server

# Modo real (requiere root)
export DRY_RUN_ONLY=false
sudo ./server
```

## Variables de Entorno

| Variable            | Default                   | Descripción                         |
|---------------------|---------------------------|-------------------------------------|
| `SERVER_HOST`       | `0.0.0.0`                | Dirección de bind                   |
| `SERVER_PORT`       | `8080`                    | Puerto HTTP                         |
| `DB_PATH`           | `./data/firewall.db`      | Ruta a la base de datos SQLite      |
| `JWT_SECRET`        | `change-me-in-production` | Secreto para firmar JWT             |
| `JWT_EXPIRATION`    | `24h`                     | Duración del token                  |
| `FIREWALL_DRIVER`   | `iptables`                | Driver por defecto                  |
| `DRY_RUN_ONLY`      | `true`                    | Si true, no ejecuta comandos reales |

## Tests

```bash
go test ./...
```

## Añadir una Nueva Vista al Frontend

1. Crear `web/js/views/mivista.js`:

```js
window.views = window.views || {};
views.mivista = {
    async render(main) {
        main.innerHTML = '<div class="loading">Loading</div>';
        const data = await api.get('/ruta');
        main.innerHTML = `<div class="container">...</div>`;
    }
};
```

2. Agregar `<script src="js/views/mivista.js"></script>` en `web/index.html`

3. Agregar ruta en `web/js/app.js`:

```js
case '#/mi-ruta': views.mivista.render(main); break;
```

## Estructura del Frontend

```
web/
├── index.html                        # HTML principal + script tags
├── css/
│   ├── style.css                     # Reset, tokens, layout, componentes base
│   └── components.css                # Estilos de componentes específicos
└── js/
    ├── api.js                        # Cliente HTTP con JWT
    ├── app.js                        # Router SPA
    ├── components/
    │   ├── navbar.js                 # Barra de navegación
    │   ├── rules-table.js            # Tabla de reglas
    │   ├── rule-form.js              # Formulario de regla
    │   └── audit-log.js              # Entrada de log de auditoría
    └── views/
        ├── login.js                  # Pantalla de login
        ├── dashboard.js              # Dashboard con estadísticas
        ├── policies.js               # CRUD de políticas
        └── audit.js                  # Visor de logs de auditoría
```
