# Desafio
Esta aplicação é um sistema de cadastro de produtos que mantém e sincroniza dados em duas bases de dados distintas (MongoDB para aplicação *ms-go* e SQLite em memória para aplicação *ms-rails*) usando mensageria (Kafka) para comunicação entre as aplicações (tópicos **rails-to-go** e **go_to_rails**).

## Instruções
- Tenha o Docker e o Docker-Compose instalados, caso não tenha siga esta página abaixo:

```
https://docs.docker.com/compose/install/
```

- Clone este repositório para sua máquina no terminal na pasta desejada com este comando:
```
git clone https://github.com/aempinto02/teste-icasei-backend-2024.git
```

- Dentro da pasta *teste-icasei-backend-2024* no terminal execute o comando:

```
docker-compose up
```

As aplicações (**ms-go** e **ms-rails**) irão subir executando os respectivos testes unitários (executados no Dockerfile de cada aplicação).

Caso encerre as aplicações, garanta com o seguinte comando que os serviços todos do docker-compose estão e serão encerrados com o seguinte comando:

```
docker-compose down
```

Este último comando garante que ao subir novamente as aplicações com *docker-compose up* nada se perca nos serviços.

## Testes com o Insomnia

- Na raiz do projeto há o arquivo *Insomnia_teste_backend.json*. Baixe em seu computador e importe com o Insomnia, caso não tenha, utilizar o site abaixo:
```
https://insomnia.rest/download
```

- No Insomnia é possível fazer as requisições às duas aplicações (*ms-go* e *ms-rails*). Caso a requisição seja *_POST_* ou *_PATCH_*, ela irá inserir no banco de dados (*MongoDB* na aplicação *ms-go*; *SQLite* na aplicação *ms-rails*) a criação no caso *_POST_* e atualização no caso de *_PATCH_* e cada aplicação irá enviar para um tópico a criação/atualização que será consumida pela outra aplicação

- A aplicação *ms-go* envia para o tópico *go_to_rails* e a aplicação *ms-rails* envia para o tópico Kafka *rails-to-go*

- O consumer da aplicação *ms-rails* é um pouco mais lento que o consumer da *ms-go*, por isso espere um pouco, caso tenha feito um *_POST_* na *ms-go* e garanta que foi adicionado na *ms-rails* utilizando o *_GET_*_Index_ no Insomnia, antes de enviar um *_POST_* na *ms-rails*

## Testes

- Os testes na aplicação *ms-rails* utilizaram o *_RSpec_* como dependência

- Os testes na aplicação *ms-go* não utilizam nenhuma dependência externa, apenas um Mock para MongoDB

- Ambos os testes rodam no Dockerfile de cada aplicação como passo antes de subir a construção da imagem
