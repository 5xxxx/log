# log

```
go get  github.com/5xxxx/log
```

## Init

```	
    log.Init()
    log.Info("test info")
	log.Debug("test debug")
	log.Warn("test warn")
	log.Error("test error")
	log.Panic("test panic")
```

## file

```
    log.Init(options.LogType(options.File),
		options.Path("./logs"),
		options.MaxAges(10),
		options.MaxBackups(10),
		options.MaxSize(30),
		options.Compress(true),
		options.LocalTime(true))
		log.Info("test info")
		
	log.Debug("test debug")
	log.Warn("test warn")
	log.Error("test error")
	log.Panic("test panic")
```

## redis 

```
    r := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", 
		DB:       0,                  
	})
	log.Init(options.LogType(options.Elastic),
		options.RedisClient(r))
	log.Init()
```