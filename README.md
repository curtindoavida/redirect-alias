# redirect-alias (Go)

Servidor HTTP em Go que replica as regras de redirecionamento do `JsMain` (Java Servlet).

## Requisitos
- Go 1.22+

## Rodar localmente
```bash
# Windows PowerShell
$env:PORT="8080"; go run .
# ou
remove-item env:PORT -ErrorAction SilentlyContinue; go run .
```

A aplicação inicia em `http://localhost:8080`.

## Build
```bash
go build -o redirect-alias.exe .
```

## Deploy
- Defina a variável de ambiente `PORT` (padrão 8080).
- Rode o binário:
```bash
./redirect-alias.exe
```

## Observação
- `robots.txt` retorna:
```
User-agent: *
Disallow: /
```
- A regra tipo 1 preserva o caminho e querystring, substituindo o host no URL.
- A regra tipo 2 redireciona para `https://<destino>` já com UTM embutido quando aplicável.

## Configuração externa
- O app lê as regras de `config/rules.json` por padrão.
- Você pode definir `RULES_PATH` para outro caminho.

Estrutura do JSON:
```json
{
  "defaultRedirect": "startupventurebuilder.com",
  "rules": {
    "host.exemplo.com": { "redirectTo": "destino.com", "type": 1 },
    "host2.exemplo": { "redirectTo": "www.site.com?utm_source=x", "type": 2 }
  }
}
```


