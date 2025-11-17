# Практическое задание № 8 Борисов Денис Александрович ЭФМО-01-25

Цели:

1.	Освоить базовые операции работы с Redis из Go-приложения.
2.	Научиться использовать команды SET, GET, задавать время жизни ключей (TTL).
3.	Реализовать кэширование данных для ускорения работы API.
4.	Понять, в каких случаях кэш помогает снизить нагрузку на базу данных.

# Выполнение практического задания.

1.	Подготовка окружения.

  	Для выполнения задания был создан файл docker-compose.yml
  	
<img width="341" height="280" alt="Снимок экрана 2025-11-16 183018" src="https://github.com/user-attachments/assets/0d671ac1-a599-4085-92a8-75db2023eec8" />

   После был запущен сервис MongoDB, при помощи Docker

<img width="710" height="271" alt="Снимок экрана 2025-11-16 183036" src="https://github.com/user-attachments/assets/8eb935ea-30da-4260-851b-815411e55265" />

   Запущенныйй сервис MongoDB

<img width="816" height="83" alt="Снимок экрана 2025-11-16 183045" src="https://github.com/user-attachments/assets/d2a251f7-9de6-466f-8909-034115680db1" />

   Проверка работы сервиса

<img width="1071" height="313" alt="Снимок экрана 2025-11-16 183145" src="https://github.com/user-attachments/assets/01d49fad-b6bf-46a2-b321-bac4a19b04b5" />

2. Структура проекта

  Для выполнения практической работы была сделана следующая структура проекта

<img width="181" height="383" alt="image" src="https://github.com/user-attachments/assets/fc2cd1b8-f6ce-4e4a-88d1-56b78817bb45" />

   Так же были установлены все необходимые расширения для выполнения практики

<img width="664" height="484" alt="Снимок экрана 2025-11-16 184314" src="https://github.com/user-attachments/assets/762f445f-ec9e-4bb2-8679-10b4b4ac9cf3" />

3.	Реализуем подключение к MongoDB (таймауты + ping).
  Для реализации подключения к MongoDB был создан файл mongo.go

<img width="590" height="649" alt="Снимок экрана 2025-11-16 184508" src="https://github.com/user-attachments/assets/c2c334bc-251b-43ed-a728-8c6a7ce2b49e" />


4. Модель данных и репозиторий.
 Для работы с заметками была написана модель данных в файле model.go, а так же репозиторий в файле repo.go. В котором реализован CRUD для заметок.

   Файл model.go

<img width="521" height="308" alt="Снимок экрана 2025-11-16 184619" src="https://github.com/user-attachments/assets/59b64a61-b318-4db9-ab1e-3d17010c0622" />

   Файл repo.go

<img width="583" height="942" alt="Снимок экрана 2025-11-16 184824" src="https://github.com/user-attachments/assets/1025239d-6f67-4dc9-9fb9-ec0923b2280f" />

5. HTTP-обработчики и маршрутизация

   Затем был создан файл handler.go, где происходит обработка запросов

<img width="588" height="954" alt="Снимок экрана 2025-11-16 185112" src="https://github.com/user-attachments/assets/069a38ad-4de0-47c6-9e8a-04b8635b2b80" />


6. Точка входа и запуск сервера

   Для запуска сервера был написан main.go

<img width="631" height="808" alt="Снимок экрана 2025-11-16 185926" src="https://github.com/user-attachments/assets/2adb080d-1822-4440-9265-e26fc766d874" />

7. Мини-тесты (Go testing)

   Для проверки основной работоспособности сервера был написан файл notes_test

<img width="762" height="820" alt="image" src="https://github.com/user-attachments/assets/83dd159b-df6f-4107-8b8c-629625605961" />


# Задание «со звёздочкой» (любая 1–2 опции — бонус к баллам)

3. Пагинация по курсору. Вместо skip/limit — «после id» (передавайте after=<ObjectID> и фильтруйте {"_id": {"$lt": ObjectID(after)}} + сортировка по _id:-1).

   Для реализации пагинации был изменен код функциии List в файле repo.go. Где получается значения after, а после происходит фильтрация по id

   <img width="629" height="478" alt="Снимок экрана 2025-11-16 204419" src="https://github.com/user-attachments/assets/5386abef-8e99-4a66-8211-013b451d1382" />

   После для создания запроса на вывод данных, был обновлен код функции listCursor в файле handler.go

   <img width="510" height="273" alt="Снимок экрана 2025-11-16 205235" src="https://github.com/user-attachments/assets/af02dad0-26f0-4fe1-a549-898263605f47" />

4.	Aggregation pipeline. Верните статистику: количество заметок, средняя длина content.

   Для создания вывода статистики была написана функция Stats, в которую собирается количество записей и высчитывается средняя длина заметки.

   <img width="705" height="687" alt="image" src="https://github.com/user-attachments/assets/9030b728-4c77-4e5e-8db2-055b1beb0fd5" />

   Также была добавлена функция для формирования запроса на выдачу статистики в handler.go

   <img width="471" height="229" alt="Снимок экрана 2025-11-16 205332" src="https://github.com/user-attachments/assets/8c74a878-af89-47a6-8fb6-789e986f8e03" />

# Проверка работоспособности

  Для проверки работоспособности был включен контейнер MongoDB, а так же запущеные основные тесты

  <img width="911" height="189" alt="image" src="https://github.com/user-attachments/assets/f661a157-d550-4869-833d-06a4f5555ab3" />

  Посел в Postman были проверено:

Создание заметки

<img width="453" height="656" alt="image" src="https://github.com/user-attachments/assets/ee4fb742-00a2-4187-bb96-3382a753a612" />

Вывод по id

<img width="711" height="678" alt="image" src="https://github.com/user-attachments/assets/c68db3dd-e8a9-4190-83cd-6496b72ff4f1" />

Частичное обновление

<img width="710" height="660" alt="image" src="https://github.com/user-attachments/assets/ad5c753f-7e35-43ab-9309-f6a96f9744f9" />

Удаление

<img width="711" height="581" alt="image" src="https://github.com/user-attachments/assets/406caf00-dc23-4a84-9355-cfc65e9d1e4b" />

Пагинация по курсору

Первая страница

<img width="709" height="810" alt="image" src="https://github.com/user-attachments/assets/f3d0c95e-08d1-4b8e-bc04-ef2df163842a" />

Следующая страница

<img width="700" height="792" alt="image" src="https://github.com/user-attachments/assets/a2cc852d-c2ad-4f9a-8cce-60015126dfe5" />


Aggregation pipeline

<img width="711" height="614" alt="image" src="https://github.com/user-attachments/assets/9a7d7494-97a0-4861-8db6-1e5a64e62d0c" />

