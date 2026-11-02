# Registro de uso de IA — Luiz (R1, R2, R3)

> **Leia antes de usar.** Este arquivo é um **rascunho de apoio** para a Declaração de Uso de Inteligência Artificial. O texto foi redigido pelo Claude a partir do histórico da própria sessão. Por isso:
>
> - Revise cada linha e reescreva com as suas palavras o que for entrar na Declaração oficial. A Declaração é sua, e o professor pode perguntar sobre ela na arguição.
> - As seções 5 (autoavaliação) e 6 (termo e assinatura) foram deixadas **em branco de propósito**. Só você pode respondê-las.
> - O fato de este rascunho ter sido escrito com IA também precisa ser declarado (há uma linha sobre isso no fim da seção 4).
> - O modelo oficial da Declaração é o documento enviado pelo professor. A estrutura abaixo segue as seis seções descritas no Guia de Uso Ético (seção 7). Confira se ela bate com o modelo.

---

## 1. Identificação

| Campo | Valor |
|---|---|
| Nome | _(preencher)_ |
| Matrícula | _(preencher)_ |
| Disciplina | T310 — Redes Convergentes (Prof. Bruno Lopes) |
| Atividade | Servidor HTTP em Go com sockets TCP |
| Dupla | Luiz e Nelson |
| Requisitos sob minha responsabilidade | R1 (servidor TCP concorrente), R2 (parsing da requisição), R3 (arquivos estáticos) |

## 2. Nível de uso de IA

**Níveis 1, 2 e 3**, dentro do autorizado pelo enunciado (níveis 1 a 3):

- **Nível 1 (consulta conceitual):** explicações sobre sockets, goroutines, `defer`, conexões persistentes do HTTP/1.1, IPv6/`::1`, Content-Length, TCP como fluxo de bytes, path traversal.
- **Nível 2 (planejamento):** um roteiro de estudo sem código, com conceitos, perguntas-guia e checklists para cada requisito.
- **Nível 3 (revisão):** sete revisões do código que eu escrevi. A IA apontou problemas e fez perguntas, sem reescrever o código. Todas as alterações foram feitas por mim.

Nenhuma linha de código do projeto foi gerada pela IA.

## 3. Ferramentas utilizadas

| Ferramenta | Uso |
|---|---|
| Claude (Anthropic), modelo Claude Opus 5.5, no app Claude Code para desktop | Toda a interação descrita abaixo, numa única sessão, de 30/10/2026 a 02/11/2026 |
| Roteiro gerado (página do Claude) | https://claude.ai/artifact/89qgZqRSifoTTEXvMyEpJ7 |

Durante as revisões, a ferramenta **leu os arquivos do meu projeto** (`main.go`, `server.go`, `public/`) e **compilou e executou cópias do meu código em uma pasta temporária** para testá-lo com `curl` (requisições válidas, malformadas e de path traversal). Ela não editou nenhum arquivo do projeto. A única exceção é este arquivo `docs/`, criado a meu pedido.

## 4. Descrição do uso

### 4.1 Pergunta inicial: o uso era permitido?

Enviei ao Claude o enunciado do T310 e o Guia de Uso Ético de IA. Perguntei se seria permitido pedir um guia que explicasse como começar o projeto **sem entregar código**. A resposta foi que sim, dentro dos níveis 1 e 2, desde que eu declarasse o uso e escrevesse todo o código. A ferramenta observou que o enunciado já explica bem o lado HTTP e que faltava o lado Go: setup, pacote `net`, goroutines e argumentos de linha de comando.

**O que recebi:** um roteiro de estudo (link acima) com o setup do projeto, os conceitos de cada requisito, perguntas para eu responder antes de codar, armadilhas de teste, checklists e materiais de estudo (A Tour of Go, Go by Example, documentação em pkg.go.dev, Kurose & Ross cap. 2, MDN, RFC 9112, Beej's Guide). Depois, a meu pedido, as seções do R2 e do R3 foram ampliadas no mesmo formato, também sem código.

**Como usei:** estudei pelos materiais indicados e pelos exemplos de código do professor (echo-server/echo-client) e respondi às perguntas do roteiro antes de implementar.

### 4.2 Revisões do meu código (nível 3)

Em cada rodada, eu escrevi o código e pedi uma revisão. A ferramenta apontou problemas em forma de perguntas, sem reescrever o código. A tabela registra o que foi apontado e o que eu decidi.

| # | Requisito | O que a IA apontou | Minha decisão |
|---|---|---|---|
| 1 | R1 | Uma aba do navegador gerava várias conexões. A IA explicou que o Chrome abre conexões antecipadas e pede o `favicon.ico`, e que o meu servidor ainda era o echo do exemplo: devolvia a própria requisição em vez de uma resposta HTTP. | Concluí que HTTP é "um pedido, uma resposta". Removi o laço `for` de leitura e passei a escrever uma resposta HTTP válida e a fechar a conexão. |
| 2 | R1 | A porta estava fixa no código. | Pesquisei `os.Args` e `flag`. Escolhi `flag` por ter valor padrão e descrição prontos, e porque `os.Args` exigiria tratar a ausência de argumento à mão. |
| 3 | R1 | O log "Encerrada" só saía quando o erro era `io.EOF`. Os logs tinham espaço duplo e textos fora da especificação. | Usei `defer` para o log de encerramento e ajustei os logs ao formato do enunciado. |
| 4 | R1 | Perguntas de arguição: ordem dos `defer`, `Connection: close`, endereço `[::1]`, Content-Length escrito à mão. | Respondi por pesquisa própria. A IA corrigiu nuances: quem escolhe o IPv6 é o cliente, e não o servidor; e o servidor fecha a conexão por causa do `defer`. |
| 5 | R2 | Uma única leitura de 1024 bytes não garante a requisição inteira. As conexões vazias podiam derrubar o parser. | Passei a usar `bufio.Reader`, a ler linha a linha até a linha em branco, e a validar a request line antes de acessar as partes. |
| 6 | R2 | A resposta HTTP tinha sumido na reescrita, o servidor dava "Empty reply". `ProtoVersion` carregava `\r\n`. O Content-Length não chegava ao R5. Content-Length negativo era aceito. O path incluía a query string. | Recoloquei a resposta, limpei a versão com `TrimSpace`, coloquei `ContentLength` e o leitor (`BodyReader`) na struct por causa do buffer do `bufio`, rejeitei valores negativos, separei a query com `strings.Cut` e exigi exatamente 3 partes na request line. |
| 7 | R3 | Sugeriu separar o código em funções. | Criei `parseRequest` e `serveFile`, e um `switch` por método no `HandleClient`. |
| 8 | R3 | **Path traversal:** a ferramenta criou, numa cópia temporária, uma pasta vizinha `public_secreto` e mostrou que `/../public_secreto/senha.txt` era servido, porque a comparação de prefixo aceitava `...\public_secreto`. | Corrigi acrescentando o separador de caminho ao final do caminho absoluto de `public` antes da comparação. Retestado: bloqueado. |
| 9 | R3 | O `serveFile` não devolvia nada, então os erros não viravam resposta. Depois de eu passar a devolver o erro, a IA mostrou que o `HandleClient` ainda descartava esse erro. | Fiz o `serveFile` devolver `error`, criei o sentinel error `ErrAccessDenied` e um `switch` com `errors.Is` no `HandleClient` (404, 403, 500), além de uma função genérica de resposta de erro. |
| 10 | R3 / R4 | Alertou que a resposta 404 com página HTML e o log `[404]` fazem parte do **R4, que é do Nelson**, e que eu não deveria implementar o requisito dele. | Deixei o tratamento das respostas de erro HTTP (R4) e do POST (R5/R6) para o Nelson. A minha parte termina no gancho que identifica o erro. |

Sugestões que **não adotei** ou deixei como opcionais: usar `filepath.Ext` no lugar da cadeia de `HasSuffix` (mantive a minha versão), o header `Allow` no 405 e o deadline de leitura em `net.Conn`.

_(Atualize esta linha antes de entregar.)_ Pendências apontadas na última revisão: o `if !found` redundante do `strings.Cut`, os prints de depuração (inclusive os `println` embutidos), rodar o `gofmt` e remover a pasta de teste `public_secreto/`.

### 4.3 O que não foi feito pela IA

- Nenhum código do projeto foi escrito, corrigido ou colado pela IA.
- As respostas às perguntas conceituais foram pesquisadas por mim. A IA apenas comentou se estavam corretas ou incompletas.
- A IA rodou o código para testes, mas toda correção foi feita por mim no meu editor.

### 4.4 Este documento

Este registro foi redigido pelo Claude, a meu pedido, a partir do histórico da sessão, para eu enviar ao Nelson (que vai montar o relatório da dupla). Eu revisei e ajustei o conteúdo.

## 5. Autoavaliação ética

_(Responda com as suas palavras. Perguntas que podem ajudar:)_

- Consigo explicar cada linha do R1, R2 e R3 sem consultar anotações?
- Em algum momento a IA me deu algo que eu usei sem entender? Se sim, o que fiz a respeito?
- As sugestões aceitas foram verificadas na documentação oficial?
- O uso ficou dentro dos níveis autorizados (1 a 3)?

## 6. Termo de responsabilidade e assinatura

_(Usar o texto do modelo oficial. Assinar via gov.br em https://assinador.iti.br ou de forma manuscrita.)_
