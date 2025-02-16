ТЗ - Нужно написать прототип (консольный) с RTP - 97%, и 3% - возврат провайдеру.

Механика для клиента: 
1) Минимальный коэффициент - 1.0х
2) Игрок делает ставку. В Pakakumi дается на это 5 секунд - мое мнение это мало, нужно хотя бы 10-15 (на подумать). А то прям сильно быстрый конвеер получается, плюс больше людей подключится к партии одной. Некоторые африканцы суеверные может, хотят успеть поклониться богам на удачу))
3) Auto cash out - это коэффициент, при котором ставка атоматически срабатывает, если коэффициент остановился на цифре меньше указаной пользователем - то он проигрывает.
4) Чистый выигрыш пользователя = Коэффициент * Ставку


Особенности со стороны провайдера:
1) Минимальная сумма для участия в игре (10 кенийских шиллингов в Pakakumi)
2) Минимальный коэффициент Auto cash out-a (1.01x в Pakakumi)
3) Минимальная сумма пополнения при депозите (50 кенийских шиллингов в Pakakumi)
4) Предлагаю брать небольшую комиссию за вывод средств, к примеру 0,5% + комиссия платежного агрегатора
5) Минимальная сумма вывода (100 шиллингов в Pakakumi)
6) Реферальная система - работает по ссылке с уникальным кодом и дает 30% прибыли от пополнений "друзей"


### Frontend часть

Предлагаю сделать на Angular 19 версии (самой последней). За счет использования новых технологий Signals (по обновлению состояний и отхода от Zone.js и Change detection) что положительно сказывается на оптимизации и скорости. Также используя Standalone компоненты, мы уменьшим размер билда фронта за счет технологии TreeShaking (автоматическое удаление из билда не используемого кода, как собственного так и в импортируемых библиотеках) 
- Настроить Докер для фронта и nginx сервер для статики
- СI/CD скрипты для деплоймента и коммита/пуша в репозиторий

### Backend часть

Будет писаться на Golang -  потому что это быстро, типизированно и билдится в бинарник. Так же на производительность будет очень положительно влиять многопоточность через GoRoutines и расположения языка к микросервисной архитектуре (для перспективы развития проекта)
- Базу данных используем PostgreSQL
- Для кеширования данных статистики - Redis
- Для передачи данных realtime - WebSockets
- Написание авто тестов для бизнес логики.
- Тоже будет настроен Docker для бекенда с несколькими environment-ами (dev, stage, prod) - можно пока просто dev и prod.

### Devops часть

Так как проект является высоконагруженым (потенциально) и пользователю нужен максимальный отклик, а нам как организаторам - защита от DDOS, то будут применены следующие технологии и подходы:

- Для ддос защиты можно поставить Cloudflare или AWS Shield
- Настроить фронтенд через CDN - Cloudflare или AWS CloudFront
- Для масштабирования можно использовать Kubernetes
- Базу данных для безопасности разместить на RDS (Amazon)


### Логика и вычисление коэффициента конца игры

коэффициент конца игры = 1 / 0.97 (RTP) = 1.03 (наша обязательная комиссия)

для добавления определенного рандома в диапазоне  будет использоваться криптографически безопасный генератор случайных чисел (RNG) 

В этом случае будет использоваться библиотека `"crypto/rand"`,  так же есть вариант заморочиться, и получать рандомный диапазон при помощи api с сайта random.org
к примеру через Gaussian Random Number Generator, который берет случайные данные из 
атмосферного шума по датчикам.

основная формула генерации условно рандомного коэффициента конца игры:

```
const RTP = 0.97 // ртп клиента
const houseEdge = 0.03 // гарантированные 3 процента наши


n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // безопасный рандом
randomNumber := float64(n.Int64()) / 1000000.0  // дробное число в диапазоне от 0 до 1

`endGameMultiplier := 1.0 / (1.0 - (1.0 - RTP - houseEdge) * randomNumber)`


```


### Честность игры

Пользователь сможет удостовериться честности игры со стороны сервера при помощи технологии *Provably Fair*.  В данном случае используются хешированные сиды от клиента и сервера, по такому принципу:

`SHA256("ServerSeed" + "ClientSeed" + "Nonce")`
Nonce - счетчик запросов, который инкрементируется с каждый раундом

т.е пример будет выглядеть так:

`SHA256("secretServerSeed" + "player123" + "0") = "f5c3a0b1e4b3d0a5f0bfe5bfa19d9c8f7c8a9f8e..."`