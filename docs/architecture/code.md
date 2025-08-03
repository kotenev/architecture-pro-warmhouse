```puml
@startuml
title "C4 Code: Sequence Diagram - Отправка команды устройству"

!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Sequence.puml

autonumber

actor "Пользователь" as user
participant "SPA" as spa
participant "API Gateway" as gateway
participant "Control Service" as control_service
participant "Device Service" as device_service
participant "Message Broker" as broker
participant "Умное устройство" as device

user -> spa: Нажимает "Включить свет"
spa -> gateway: POST /api/v1/devices/123/control\n{"command": "turn_on"}
gateway -> control_service: POST /control\n(с данными пользователя и команды)

control_service -> device_service: GET /devices/123/auth\n(Проверить, что юзер X может управлять устройством 123)
device_service --> control_service: 200 OK

control_service -> control_service: Сформировать сообщение команды\n(e.g., {"device_id": 123, "action": "set_power", "value": "on"})

control_service -> broker: Publish(topic: "commands.light.123", message)
broker -> device: Доставляет сообщение
activate device
device -> device: Включает свет
deactivate device

broker --> control_service: Ack (сообщение получено)
control_service --> gateway: 202 Accepted (команда принята к исполнению)
gateway --> spa: 202 Accepted
spa -> user: Показывает "Команда отправлена"

@enduml
```
