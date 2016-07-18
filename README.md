# log
日志系统，可以根据panic精确输出是哪行代码出错，是调试的好帮手


第一步先下载代码到工程环境下


第二步:在项目最初调用初始化方法：log.Init()

第三步，在你需要调试的地方log.InitListner(),注意,如果开了新的协程,并且需要在新的协程里输出日志,请在协程内部再调用log.InitListner()


第三步，在你需要调试的地方panic或者log.Write(err),其实调用log.InitListner()方法以后,就会在本协程内部进行事件监听,可以直接调用panic方法,无需更改



请注意：panic调用后,后面的代码将不能执行，当前协程退出
       log.Write(err) 在drawin下后面的代码可以继续执行,windows后面的代码将不能执行，当前协程退出

mac与linux完美进行
