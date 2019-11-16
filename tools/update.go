package main

import (
	"runtime"

	"ebook/common"
	"ebook/engines/book"
	"ebook/engines/book/robot"
	res "ebook/resource"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	common.Init("../config.json")

	robot := robot.Robot{}
	robot.Init(res.NewResource().Index) // 初始化解析器

	robot.Run() // 运行解析器命令

	// fmt.Printf("sokindle 新现图书列表共 %d 册，明细 %d 册\n", SokindleCount1, SokindleCount2)
	// 将未下载的图书下载并上传至七牛云
	robot.Each(func(bok book.Book) bool {
		robot.Command(robot.SokindleDownload, bok)

		return true
	})
	println("waiting")
	robot.Wait() // 等待任务完成

	println("======================================================")
	// if books, f := robot.Find("Isbn", "9787111314943"); f {
	// 	for _, b := range books {
	// 		println(b.Name)
	// 	}
	// }

	if books, f := robot.Match("幸存者", 0, -1); f {
		for _, b := range books {
			println(b.Name)
		}
	}
}
