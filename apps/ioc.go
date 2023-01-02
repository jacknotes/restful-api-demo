package app

import "github.com/jacknotes/restful-api-demo/apps/host"

// 控制反转（Inversion of Control）是一种是面向对象编程中的一种设计原则，用来减低计算机代码之间的耦合度。
//其基本思想是：借助于“第三方”实现具有依赖关系的对象之间的解耦
var (
	Host host.ServiceServer //在使用start命令时，会把ipml.Service对象传给此Host
)
