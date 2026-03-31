# Guía de Instalación — Conpas AI

Guía rápida para instalar Conpas AI en **Windows** y **Linux**.

---

## 📋 Requisitos Previos

### Opción 1: Binario Precompilado (Recomendado)
- **Windows**: PowerShell 5.1+ (incluido en Windows 10/11)
- **Linux**: `curl` o `wget` instalado
- **No requiere Go** — el binario ya está compilado

### Opción 2: Instalación desde Código (Desarrollo)
- **Go 1.24+** instalado y configurado
- Git (para clonar el repositorio)

---

## ⚡ Instalación Rápida (Binario Precompilado)

### Windows (PowerShell)

```powershell
# Descargar e instalar desde GitHub Releases
irm https://raw.githubusercontent.com/Thrasno/conpas-ai/main/scripts/install.ps1 | iex
```

Esto descarga el binario para Windows y lo coloca en tu `PATH`.

### Linux

```bash
# Descargar e instalar desde GitHub Releases
curl -fsSL https://raw.githubusercontent.com/Thrasno/conpas-ai/main/scripts/install.sh | bash
```

El script detecta tu arquitectura (amd64/arm64) y descarga el binario correspondiente.

---

## 🛠️ Instalación desde Código (Desarrollo)

### 1. Clonar el Repositorio

```bash
git clone https://github.com/Thrasno/conpas-ai.git
cd conpas-ai
```

### 2. Compilar el Binario

#### Windows (PowerShell)
```powershell
go build -o conpas-ai.exe ./cmd/conpas-ai
```

#### Linux / macOS
```bash
go build -o conpas-ai ./cmd/conpas-ai
```

### 3. (Opcional) Instalar Globalmente

#### Windows
Mover el binario a un directorio en tu `PATH`:
```powershell
Move-Item conpas-ai.exe C:\Windows\System32\
```

O añadir el directorio actual a tu `PATH` en Variables de Entorno del Sistema.

#### Linux / macOS
```bash
sudo mv conpas-ai /usr/local/bin/
```

---

## 🚀 Primer Uso

Una vez instalado, ejecuta el comando principal:

```bash
conpas-ai
```

Esto lanza el **TUI interactivo** donde puedes:
- Seleccionar tu agente AI (Claude Code, OpenCode, Gemini CLI, Cursor, etc.)
- Elegir un preset de instalación (`minimal`, `full`, `custom`)
- Seleccionar una persona (gentleman, littleYoda, zen-master, etc.)
- Configurar componentes, skills y MCP servers

### Modo CLI (Sin TUI)

Si prefieres el modo no-interactivo:

```bash
# Instalar preset completo para Claude Code con persona gentleman
conpas-ai --agent claude-code --preset full --persona gentleman

# Ver todas las opciones
conpas-ai --help
```

---

## 🔍 Verificar Instalación

Comprobar que el binario está disponible:

```bash
conpas-ai --version
```

Debería mostrar la versión instalada.

---

## 🐛 Troubleshooting

### Windows

#### Error: "No se puede ejecutar scripts en este sistema"
PowerShell está bloqueando la ejecución de scripts. Ejecuta:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### Error: "conpas-ai no se reconoce como comando"
El binario no está en tu `PATH`. Opciones:
1. Reinicia PowerShell/Terminal
2. Añade manualmente el directorio donde está `conpas-ai.exe` al `PATH` del sistema
3. Ejecuta usando la ruta completa: `.\conpas-ai.exe`

### Linux

#### Error: "Permission denied"
El binario no tiene permisos de ejecución:
```bash
chmod +x conpas-ai
```

#### Error: "command not found"
El binario no está en tu `PATH`. Opciones:
1. Instálalo en `/usr/local/bin/` (ver sección anterior)
2. Añade el directorio actual a tu `PATH`:
   ```bash
   export PATH=$PATH:$(pwd)
   ```
3. Ejecuta usando la ruta completa: `./conpas-ai`

### Ambas Plataformas

#### Error de dependencias Go (solo si compilas desde código)
Asegúrate de tener Go 1.24 o superior:
```bash
go version
```

Si necesitas actualizar Go, descárgalo desde [go.dev/dl](https://go.dev/dl/).

---

## 📚 Documentación Completa

Para más información sobre uso avanzado, componentes, skills y arquitectura, consulta:
- [README.md](README.md) — Visión general del proyecto
- [docs/usage.md](docs/usage.md) — Uso detallado del CLI y TUI
- [docs/platforms.md](docs/platforms.md) — Notas específicas por plataforma
- [docs/components.md](docs/components.md) — Catálogo de skills y presets

---

## 🆘 Soporte

Si encuentras problemas:
1. Revisa la sección [Troubleshooting](#-troubleshooting) arriba
2. Consulta los [Issues](https://github.com/Thrasno/conpas-ai/issues) en GitHub
3. Abre un nuevo issue siguiendo la plantilla de bug report

---

**Nota**: Este proyecto es una herramienta interna corporativa de Thrasno para estandarizar configuraciones de agentes AI. Cada miembro del equipo puede configurar su stack preferido sin romper convenciones del equipo.
