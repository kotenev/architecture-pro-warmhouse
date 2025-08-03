```puml
@startuml
!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Container.puml

title Container diagram for Smart Home System

' Определяем действующих лиц и внешние системы
Person(user, "User (Пользователь)", "Клиент, использующий систему 'Умный дом' через веб-интерфейс.")
System_Ext(smart_device, "Smart device (Умное устройство)", "Любое поддерживаемое устройство (датчик, реле, камера и т.д.)")

' Определяем границы системы
System_Boundary(c1, "Smart Home Ecosystem (Экосистема 'Умный Дом')") {

    ' --- Клиентское приложение ---
    Container(spa, "Single-Page Application", "JavaScript / React", "Предоставляет пользователю интерфейс для управления умным домом")

    ' --- Внутренние сервисы (микросервисы) ---
    Container(api_gateway, "API gateway", "Go / Nginx", "Маршрутизирует внешние запросы к нужным микросервисам, обеспечивает безопасность")

    Container(user_service, "User service (Сервис пользователей)", "Go", "Отвечает за регистрацию, аутентификацию и управление данными пользователей")
    ContainerDb(user_db, "User DB (База данных пользователей)", "PostgreSQL", "Хранит информацию о профилях пользователей")

    Container(device_service, "Device service (Сервис устройств)", "Go", "Отвечает за регистрацию, настройку и управление метаданными устройств")
    ContainerDb(device_db, "Device DB (База данных устройств)", "PostgreSQL", "Хранит информацию об устройствах, их типах и настройках")

    Container(telemetry_service, "Telemetry service (Сервис телеметрии)", "Go", "Сбор, обработка и предоставление данных с датчиков")
    ContainerDb(telemetry_db, "Telemetry DB (База данных телеметрии)", "TimescaleDB", "Хранит временные ряды показаний датчиков (температура, состояние и т.д.)")

    Container(control_service, "Control service (Сервис управления)", "Go", "Отправляет команды на умные устройства")

    ' --- Брокер сообщений для асинхронного взаимодействия ---
    Container(message_broker, "Message broker (Брокер сообщений)", "RabbitMQ", "Обеспечивает асинхронную и надёжную доставку команд и данных телеметрии")
}

' --- Определяем связи ---

' Пользователь и система
Rel(user, spa, "Управляет устройствами, просматривает данные", "HTTPS")
Rel(spa, api_gateway, "Отправляет API-запросы", "HTTPS/JSON")

' API Gateway и внутренние сервисы
Rel(api_gateway, user_service, "Проксирует запросы аутентификации и управления профилем", "HTTP")
Rel(api_gateway, device_service, "Проксирует запросы на управление устройствами (CRUD)", "HTTP")
Rel(api_gateway, telemetry_service, "Проксирует запросы на получение истории данных", "HTTP")
Rel(api_gateway, control_service, "Проксирует запросы на отправку команд", "HTTP")

' Связи сервисов с их базами данных
Rel(user_service, user_db, "Читает/пишет данные", "TCP/IP")
Rel(device_service, device_db, "Читает/пишет данные", "TCP/IP")
Rel(telemetry_service, telemetry_db, "Читает/пишет данные", "TCP/IP")

' Асинхронное взаимодействие через брокер
Rel(control_service, message_broker, "Отправляет команду на выполнение", "AMQP/MQTT")
Rel(telemetry_service, message_broker, "Подписывается на данные телеметрии", "AMQP/MQTT")

' Устройства и система
Rel(smart_device, message_broker, "Отправляет данные телеметрии и получает команды", "AMQP/MQTT")

@enduml
```