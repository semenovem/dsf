## Логгирование 

используется библиотека   
github.com/sirupsen/logrus   


TODO 
переименовать p (pkg) в sys - название системы


```
# Настройка / запуск
log  = logger.New().Pkg("main").Apl()

lg.SetDefLev(lg.ParseLev(cfg.LogLevel, logrus.InfoLevel))
lg.SetDefMod(lg.ParseMod(cfg.LogMode, lg.ModeJson))

lg.


форматирование / вывод

по умолчанию
# mod=json
{"level":"info","msg":"sftp server running on port: '[::]:2022'","p":"sftp","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Connected to IBM MQ manager QM1","p":"mq  ","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.1]","p":"mq  ","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Connected to IBM MQ manager QM1","p":"mq  ","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.2]","p":"mq  ","time":"2021-04-06T14:23:28Z"}
{"level":"info","msg":"Application started","p":"main","service_start":"start","time":"2021-04-06T14:23:30Z"}


# mod=short
[INFO] [p:sftp] sftp server running on port: '[::]:2022'
[INFO] [p:mq  ] Connected to IBM MQ manager QM1
[INFO] [p:mq  ] Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.2]
[INFO] [p:mq  ] Connected to IBM MQ manager QM1
[INFO] [p:mq  ] Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.1]
[INFO] [p:main] [service_start:start] Application started


# mod=text
2021-04-06T14:22:53Z[INFO] [p:sftp] sftp server running on port: '[::]:2022'
2021-04-06T14:22:53Z[INFO] [p:mq  ] Connected to IBM MQ manager QM1
2021-04-06T14:22:53Z[INFO] [p:mq  ] Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.2]
2021-04-06T14:22:53Z[INFO] [p:mq  ] Connected to IBM MQ manager QM1
2021-04-06T14:22:53Z[INFO] [p:mq  ] Opened the manager/channel/queue: [QM1 / DEV.APP.SVRCONN / DEV.QUEUE.1]
2021-04-06T14:22:55Z[INFO] [p:main] [service_start:start] Application started




```
