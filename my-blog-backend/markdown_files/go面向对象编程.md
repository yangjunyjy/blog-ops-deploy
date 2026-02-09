面向对象编程应用实例

```go
步骤
1)声明(定义)结构体，确定结构体名
2)编写结构体的字段
3)编写结构体的方法
学生案例：
1)编写一个Student结构体，包含name、gender、age、id、score字段，分别为string、string、int、int、float64类型。
2)结构体中声明一个say方法，返回string类型，方法返回信息中包含所有字段值。
3）在main方法中，创建Student结构体实例(变量)，并访问say方法，并将调用结果打印输出。
```

代码：

```go
package main
import (
	"fmt"
)
type Student struct{
    name string
    gender string
    age int
    id int
    score float64
}
func (student Student)say()(name string, gender string,age int,id int,score float64){
    name = student.name
    gender = student.gender
    age = student.age
    id = student.id
    score = student.score
}
func main(){
    jack := Student{
        name : "jack",
        gender : "男",
        age : 22,
        id : 1,
        score : 95.5,
    }
    fmt.Println(jack.say())
}

jack 男 22 1 95.5
```

```go
工厂模式
说明
Golang的结构体没有构造函数，通常可以使用工厂模式来解决这个问题。
看一个需求,一个结构体的声明是这样的：
package model
type Student struct{
Name string...
}
因为这里的student的首字母s是大写的，如果我们想在其它包创建Student的实例(比如main包)，引入model包后，就可以直接创建Student结构体的变量（实例)。但是问题来了，如果首字母是小写的，比如是type student struct{..}就不不行了，怎么办--->工厂模式来解决.
```

案例

```go
package model

type student struct {
	Name string
	Age  int
}

func NewStudnt(n string, a int) *student {
	return &student{
		Name: n,
		Age:  a,
	}
}
func (s *student) GetName() string {
	return s.Name
}
```

```go
package main

import (
	"fmt"
	"src/go_code/main/model"
)

func main() {
	stu := model.NewStudnt("付存云", 25)
	fmt.Printf("%s的年龄是%d", stu.Name, stu.Age)
}
```

面向对象编程思想-抽象

```go
如何理解抽象
我们在前面去定义一个结构体时候，实际上就是把一类事物的共有的属性和行为提取
出来，形成一个物理模型(模板)。这种研究问题的方法称为抽象
```

![](F:\笔记\后端\golang\golang图片\面向对象-抽象.png)

面向对象编程-封装

```go
封装介绍
封装(encapsulation)就是把抽象出的字段和对字段的操作封装在一起，数据被保护在
内部，程序的其它包只有通过被授权的操作（方法），才能对字段进行操作
面向对象编程-封装
封装的理解和好处
1)隐藏实现细节
2)提可以对数据进行验证，保证安全合理
如何体现封装
1)对结构体中的属性进行封装
2)通过方法，包实现封装
```

封装的实现步骤

```go
1)将结构体、字段(属性)的首字母小写(不能导出了，其它包不能使用，类似private)

2)给结构体所在包提供一个工厂模式的函数，首字母大写。类似一个构造函数

3)提供一个首字母大写的Set方法(类似其它语言的public)，用于对属性判断并赋值
func (var结构体类型名) SetXxx(参数列表) (返回值列表){
//加入数据验证的业务逻辑
var.字段=参数
4）提供一个首字母大写的Get方法(类似其它语言的public)，用于获取属性的值
func (var结构体类型名) GetXxx(){
return var.字段;
特别说明：在Golang开发中并没有特别强调封装，这点并不像Java.所以提醒学过java的朋友，
不用总是用java的语法特性来看待Golang,Golang本身对面向对象的特性做了简化的.
```

课堂练习

```go
创建程序,在model包中定义Account结构体：在main函数中体会Golang的封装性。
1）Account结构体要求具有字段：账号（长度在6-10之间）、余额（必须>20)、密码（必须是六位）
2)通过SetXxx的方法给Account的字段赋值。(同学们自己完成）
3)在main函数中测试
```

继承

嵌套匿名结构体基本语法

```go
type Goods struct{
    Name string
    price int
}
type Book struct{
    Goods //嵌套匿名结构体Goods
    write string
}
```

继承的深入讨论

```go
1)结构体可以使用嵌套匿名结构体所有的字段和方法，即：首字母大写或者小写的字段、方法，都可以使用。
2）匿名结构体字段访问可以简化
3)当结构体和匿名结构体有相同的字段或者方法时，编译器采用就近访问原则访问，如希望访问匿名结构体的字段和方法，可以通过匿名结构体名来
区分
4)结构体嵌入两个（或多个）匿名结构体，如两个匿名结构体有相同的字段和方法（同时结构体本身没有同名的字段和方法)，在访问时，就必须明确指定匿名结构体名字，否则编译报错。
5）如果一个struct嵌套了一个有名结构体，这种模式就是组合，如果是组合关系，那么在访问组金的结构体的字段或方式时，必须带上结构体的名字
6）嵌套匿名结构体后，也可以在创建结构体变量（实例）时，直接指定各个匿名结构体字段的值.
```

多重继承说明

```go
如一个struct嵌套了多个匿名结构体，那么该结构体可以直接访问嵌套的匿名结构体的字段和方法，从而实现了多重继承。
多重继承细节说明
1)如嵌入的匿名结构体有相同的字段名或者方法名，则在访问时，需要通过匿名结构体类型名来区分。【案例演示】
2）为了保证代码的简洁性，建议大家尽量不使用多重继承
```

面向对象编程-多态
基本介绍

```go
变量（实例）具有多种形态。面向对象的第三大特征，在Go语言，多态特征是通过接口
实现的。可以按照统一的接口来调用不同的实现。这时接口变量就呈不同的形态。

快速入门
在前面的Usb接口案例，Usbusb，既可以接收手机变量，又可以接收相机变量，就体现了Usb接口多态特性。
```

接口体现多态特征

```
1）多态参数
在前面的Usb接口案例，Usbusb，即可以接收手机变量，又可以接收相机变量，就体现了Usb接口多态
2)多态数组
演示一个案例：给Usb数组中，存放Phone结构体和Camera结构体变量，Phone还有一个特有的方法call()，请遍历Usb数组，如果是Phone变量，除了调用usb接口声明的方法外，还需要调用Phone特有方法call.
```

