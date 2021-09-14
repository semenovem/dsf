## Логирование 

используется библиотека   
github.com/sirupsen/logrus   



### Настройка / запуск
```
var logHolder = logger.New()

# дефолтные установки 
err := logHolder.SetDefLevel("WARN")
if err != nil {
  fmt.Println(err)
}

err := logHolder.SetDefMode("TEXT")
if err != nil {
  fmt.Println(err)
}

# создание / настройка логера системы
logApi := logHolder.GetLog("api")
logHolder.SetLevel("api", "DEBUG")
```


### форматирование / вывод
```

# по умолчанию
# mod=json
{"level":"info","msg":"sftp server running on port: '[::]:2022'","sys":"sftp","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Connected to IBM MQ manager QM1","sys":"mq  ","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.1]","sys":"mq  ","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Connected to IBM MQ manager QM1","sys":"mq  ","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.2]","sys":"mq  ","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Application started","sys":"main","service_start":"start","time":"2021-04-06T14:23:30Z"}


# mod=short
[INFO] [sys:sftp] sftp server running on port: '[::]:2022'
[INFO] [sys:mq  ] Connected to IBM MQ manager QM1
[INFO] [sys:mq  ] Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.2]
[INFO] [sys:mq  ] Connected to IBM MQ manager QM1
[INFO] [sys:mq  ] Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.1]
[INFO] [sys:main] [service_start:start] Application started


# mod=text
2021-04-06T14:22:53Z[INFO] [sys:sftp] sftp server running on port: '[::]:2022'
2021-04-06T14:22:53Z[INFO] [sys:mq  ] Connected to IBM MQ manager QM1
2021-04-06T14:22:53Z[INFO] [sys:mq  ] Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.2]
2021-04-06T14:22:53Z[INFO] [sys:mq  ] Connected to IBM MQ manager QM1
2021-04-06T14:22:53Z[INFO] [sys:mq  ] Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.1]
2021-04-06T14:22:55Z[INFO] [sys:main] [service_start:start] Application started
```
