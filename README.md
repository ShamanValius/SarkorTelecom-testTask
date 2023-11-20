# SarkorTelecom-testTask
<details>
      <summary>Тестовое задание </summary>
      <body> <h1>Описание задачи</h1> <ol> <li> <p>Необходимо запустить локальный веб-сервер на Go (подойдет любой Http, Gin, Echo).</p> <p>Не использовать ORM, только чистый SQL. Использовать БД Sqlite, предварительно создав таблицы или вывести SQL как миграцию при запуске веб-сервера.</p> </li> <li> <p>Реализовать POST обработчик /user/register, который обрабатывает форму (PostForm).</p> <p>Поля формы:</p> <ul> <li>Login (string)</li> <li>Password (string)</li> <li>Name (string)</li> <li>Age (int)</li> </ul> <p>(Это не JSON.)</p> <p>Эти значения необходимо сохранить в файловую БД. Пароль сохранить хэшем (bcrypt).</p> </li> <li> <p>Метод аутентификации /user/auth в JSON.</p> <p>Поля JSON:</p> <ul> <li>Login (string)</li> <li>Password (string)</li> </ul> <p>Проверить логин и пароль, в случае успеха вернуть куки SESSTOKEN=&lt;JWT TOKEN&gt;.</p> <p>JWT токен обычный (HS256), с логином и user_id в полезной нагрузке.</p> </li> <li> <p>Реализовать обработчик GET /user/:name.</p> <p>(Необходим middleware авторизации, который проверяет куки SESSTOKEN.)</p> <p>Где по :name ищется в БД и возвращается в формате JSON результат:</p> <pre>{ "id": &lt;id из БД&gt;, "name": "имя", "age": 25 }</pre> </li> <li> <p>Реализовать метод добавления номера, POST /user/phone (также хранить в БД).</p> <p>(Необходим middleware авторизации, который проверяет куки SESSTOKEN.)</p> <p>В JSON:</p> <ul> <li>phone - string (max: 12)</li> <li>description - string (описание номера)</li> <li>is_fax - bool (факс или нет)</li> </ul> <p>Добавляется за пользователем, в таблице должна быть колонка user_id, который будет браться из JWT.</p> <p>Должна быть проверка на дубликат.</p> </li> <li> <p>Реализовать метод получения номера GET /user/phone?q=&lt;номер&gt;.</p> <p>(Необходим middleware авторизации, который проверяет куки SESSTOKEN.)</p> <p>Ответ в JSON: список тех, у кого есть этот номер</p> <pre>user_id, phone, description, is_fax</pre> <p>Примечание: вводится может часть номера, поиск должен возвращать массив с подходящими номерами.</p> </li> <li> <p>Реализовать метод PUT /user/phone для обновления данных номера.</p> <p>(Необходим middleware авторизации, который проверяет куки SESSTOKEN.)</p> <p>Поля:</p> <ul> <li>phone_id</li> <li>phone</li> <li>is_fax</li> <li>description</li> </ul> <p>Обновляются поля, user_id из JWT.</p> </li> <li> <p>Реализовать метод DELETE для удаления номера /user/phone/&lt;phone_id&gt;.</p> <p>(Необходим middleware авторизации, который проверяет куки SESSTOKEN.)</p> </li> <p>Приветствуется использование каких-либо архитектур и комментирования кода.</p> </ol> </body>
  </details>
<details>
      <summary>Структура папок</summary>
      <pre>
.
├── Dockerfile
├── cmd
│   └── SarkorTelecom-testTask
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   ├── DTO
│   │   │   ├── phone.go
│   │   │   └── user.go
│   │   ├── app.go
│   │   ├── http
│   │   │   ├── handlers
│   │   │   │   ├── handlerbuilder.go
│   │   │   │   ├── phone.go
│   │   │   │   └── user.go
│   │   │   ├── middleware
│   │   │   │   └── auth
│   │   │   │       └── checkAuth.go
│   │   │   └── routes
│   │   │       └── route.go
│   │   └── service
│   │       ├── jwtmanager
│   │       │   └── jwt.go
│   │       └── passHash
│   │           ├── bcrypt
│   │           │   └── bcryptHash.go
│   │           └── hasher.go
│   ├── config
│   └── storage
│       ├── sqlite
│       │   ├── sqlite.go
│       └── storage.go
└── tools
    └── Компания СП ООО -Sarkor Telecom- ==- Golang Junior-Junior+ ==- тестовое задание_2023.postman_collection.json


</pre>
  </details>

Структуру папок делал ориентируясь на [# Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README.md) и [# Практичный Go:...](https://habr.com/ru/articles/441842/).  
Для такого маленького "проекта" получилось, возможно, избыточно.

Также попытался реализовать Dependency Injection (DI), например, для работы с базой данных или сервисом хэширования пароля. Однако для сервиса обработки JWT подобное не было реализовано.    

Роуты вынесены в отдельный файл, а их обработчики разделены по работе с основными сущностями и также вынесены. Файл `handlerbuilder` инициализирует набор обработчиков.   

База данных создается при инициализации в корне приложения. 

Конфигурационный файл реализован с использованием `.env`.  
В `main.go` конфигурируются основные параметры, а в `app.go` происходит настройка роутов, обработчиков и запуск сервера.  

"Проект" можно запустить с помощью Docker, используя команду:  
docker build -t sktest . && docker run sktest  
Соединения будут прослушиваться на порту 5000.

Также в папке :file_folder:`tools` находится экспортированный файл коллекции Postman [<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 90px; height: 20px;">](https://app.getpostman.com/run-collection/27753311-95bd29f1-804d-4983-909d-9f2a9412c2b0?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D27753311-95bd29f1-804d-4983-909d-9f2a9412c2b0%26entityType%3Dcollection%26workspaceId%3D04d1f14b-3754-410e-b0d0-91d12be6ff6c)
. С его помощью можно протестировать запросы.

К сожалению, "проект" не покрыт тестами, даже модульными. Валидация входных данных не применяется, кроме как на стороне БД, что недостаточно, так как это создает лишнюю нагрузку на нее. Также отсутствует нормальное логирование (выведено только в консоль стандартными средствами).

По коду есть места с пометкой :bookmark:TODO, где явно требуется улучшение.
