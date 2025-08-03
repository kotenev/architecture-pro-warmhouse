```puml
@startuml
!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Component.puml

title Component Diagram for Smart Home System

' -- Внешние по отношению к компонентам контейнеры для контекста --
Container(api_gateway, "API gateway", "Go / Nginx", "Маршрутизирует внешние запросы к компонентам API Controller")
Container(message_broker, "Message broker (Брокер сообщений)", "RabbitMQ", "Обеспечивает асинхронную и надёжную доставку команд и данных телеметрии")
ContainerDb(user_db, "User DB (База данных пользователей)", "PostgreSQL", "Хранит информацию о профилях пользователей")
ContainerDb(device_db, "Device DB (База данных устройств)", "PostgreSQL", "Хранит информацию об устройствах, их типах и настройках")
ContainerDb(telemetry_db, "Telemetry DB (База данных телеметрии)", "TimescaleDB", "Хранит временные ряды показаний датчиков (температура, состояние и т.д.)")

' =======================================
' User Service (Сервис пользователей)
' =======================================
System_Boundary(user_service_c, "User service (Сервис пользователей)") {
Component(user_api, "API controller", "Go / Gin", "Принимает HTTP-запросы для регистрации и аутентификации")
Component(user_logic, "User logic", "Go", "Реализует бизнес-логику управления пользователями (валидация, JWT)")
Component(user_repo, "User repository", "Go", "Обеспечивает CRUD-операции с данными пользователей в БД")
}

' =======================================
' Device Service (Сервис устройств)
' =======================================
System_Boundary(device_registry_c, "Device service (Сервис устройств)") {
Component(device_api, "API controller", "Go / Gin", "Обрабатывает HTTP-запросы на управление устройствами (CRUD)")
Component(device_logic, "Device logic", "Go", "Реализует бизнес-логику управления реестром устройств")
Component(device_repo, "Device repository", "Go", "Обеспечивает доступ к метаданным устройств в БД")
}

' =======================================
' Control Service (Сервис управления)
' =======================================
System_Boundary(device_control_c, "Control service (Сервис управления)") {
Component(control_api, "API controller", "Go / Gin", "Принимает команды управления (включить, выключить и т.д.)")
Component(control_logic, "Command logic", "Go", "Валидирует команды и права доступа, формирует сообщения для брокера")
Component(amqp_publisher, "Message publisher", "Go / RabbitMQ client", "Публикует команды в брокер сообщений для доставки устройствам")
}

' =======================================
' Telemetry Service (Сервис телеметрии)
' =======================================
System_Boundary(telemetry_service_c, "Telemetry service (Сервис телеметрии)") {
Component(amqp_consumer, "Message consumer", "Go / RabbitMQ client", "Подписывается на сообщения телеметрии от устройств через брокер")
Component(telemetry_logic, "Telemetry processor", "Go", "Обрабатывает и валидирует входящие данные с датчиков")
Component(telemetry_repo, "Telemetry repository", "Go", "Обеспечивает запись обработанных данных телеметрии в Time-Series БД")
}


' -- Связи (Relationships) --

' Связи с API Gateway
Rel(api_gateway, user_api, "Запросы по пользователям (/register, /login)", "HTTP / JSON")
Rel(api_gateway, device_api, "Запросы по устройствам ( /devices )", "HTTP / JSON")
Rel(api_gateway, control_api, "Команды управления ( /devices/{id}/control )", "HTTP / JSON")

' Внутренние связи компонентов
Rel(user_api, user_logic, "Использует")
Rel(user_logic, user_repo, "Использует")
Rel(user_repo, user_db, "Читает/пишет данные", "SQL")

Rel(device_api, device_logic, "Использует")
Rel(device_logic, device_repo, "Использует")
Rel(device_repo, device_db, "Читает/пишет данные", "SQL")

Rel(control_api, control_logic, "Использует")
Rel(control_logic, amqp_publisher, "Отправляет команды")
Rel(amqp_publisher, message_broker, "Публикует сообщения", "AMQP")

Rel(message_broker, amqp_consumer, "Доставляет сообщения", "AMQP")
Rel(amqp_consumer, telemetry_logic, "Передает данные для обработки")
Rel(telemetry_logic, telemetry_repo, "Использует для сохранения обработанных данных")
Rel(telemetry_repo, telemetry_db, "Пишет данные", "SQL")

@enduml
```