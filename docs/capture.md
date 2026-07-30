# Captura de Estado del Firewall

## Propósito

El paquete `internal/capture` permite leer la configuración **actual y en vivo** del firewall del sistema ejecutando `iptables-save -c`. Esto complementa el flujo existente de políticas almacenadas en BD con una vista de la realidad del firewall.

## Uso desde API

```
GET /api/v1/system/firewall
```

Requiere autenticación (cualquier rol). Ejecuta `iptables-save -c` y devuelve el resultado parseado como JSON.

## Uso desde Go

```go
import "github.com/anomalyco/iptables-visualizer/internal/capture"

state, err := capture.Capture()
if err != nil {
    log.Fatal(err)
}
// state.Tables contiene todas las tablas, cadenas y reglas
```

## Arquitectura

### `capture.Capture()`
1. Ejecuta `iptables-save -c` vía `os/exec`
2. Parsea la salida con `parse()`
3. Devuelve `*FirewallState`

### `parse(data string)`
Procesa línea por línea:
- `*table` → nueva tabla
- `:chain POLICY [pkts:bytes]` → nueva cadena con política y contadores
- `-A chain ...` → nueva regla parseada con `parseRule()`
- `COMMIT` → fin de tabla

### `parseRule(line string)`
Tokeniza respetando quotes dobles y extrae campos conocidos:
- `-j` → target (ACCEPT, DROP, etc.)
- `-p` → protocolo
- `-s` / `-d` → direcciones
- `--sport` / `--dport` → puertos
- `-i` / `-o` → interfaces
- `-m state --state` → estado de conexión
- `--log-prefix` → prefijo de log
- Campos no reconocidos van a `Extra`

### Estructuras

```go
type FirewallState struct {
    Tables    []Table   `json:"tables"`
    Timestamp time.Time `json:"timestamp"`
}

type Table struct {
    Name   string           `json:"name"`
    Chains map[string]Chain `json:"chains"`
}

type Chain struct {
    Policy  string `json:"policy"`
    Rules   []Rule `json:"rules"`
    Packets uint64 `json:"packets,omitempty"`
    Bytes   uint64 `json:"bytes,omitempty"`
}

type Rule struct {
    Chain        string   `json:"chain"`
    Target       string   `json:"target"`
    Protocol     string   `json:"protocol,omitempty"`
    Source       string   `json:"source,omitempty"`
    Destination  string   `json:"destination,omitempty"`
    SrcPort      string   `json:"src_port,omitempty"`
    DstPort      string   `json:"dst_port,omitempty"`
    InInterface  string   `json:"in_interface,omitempty"`
    OutInterface string   `json:"out_interface,omitempty"`
    State        string   `json:"state,omitempty"`
    LogPrefix    string   `json:"log_prefix,omitempty"`
    Extra        []string `json:"extra,omitempty"`
    Raw          string   `json:"raw"`
}
```

## Limitaciones

- Requiere que `iptables-save` esté instalado y accesible
- Requiere permisos de root (o `CAP_NET_ADMIN`) para ejecutar `iptables-save`
- Solo captura iptables, no nftables
- El parseo de reglas es heurístico: campos no reconocidos caen en `Extra`
