# [Go为什么能成功？](https://mp.weixin.qq.com/s/02Dtj94yOjy2WTLko-5JXQ)

[Go语言中文网](javascript:void(0);) *2022-12-20 08:52* *Posted on 北京*

The following article is from TonyBai Author 白明的赞赏账户

[![img](http://wx.qlogo.cn/mmhead/Q3auHgzwzM7SvW23NJrPgbQAdoXIq2nIlCY3bc9S4kOg7UOiap48H2w/0)**TonyBai**.与技术博客tonybai.com同源。近期关注Kubernetes、Docker、Golang、儿童编程、DevOps、云计算平台和机器学习。](https://mp.weixin.qq.com/s/02Dtj94yOjy2WTLko-5JXQ#)

大家在入门Go语言时，多埋头于Go语法，忙于练手或快速完成公司的项目，无暇思考。

但当大家到了要进阶，要冲刺高级阶段的时候，我建议你不能再稀里糊涂了。既然入了Go这个坑，在进入高级阶段前，我们最好在门口的“影壁墙”前驻留一下。

![Image](https://mmbiz.qpic.cn/mmbiz_jpg/cH6WzfQ94mbE1j1BTNg6rvgTBiap4eicP2S6NFV5phBHpTibJ5H2AL6A9xro1FstHiabkyibs55ystHc1CXZ01wGUiaQ/640?wx_fmt=jpeg&wxfrom=5&wx_lazy=1&wx_co=1)

仔细思考一下**我们投入这么多精力研究的Go为什么能成功**，后续还能否持续成功下去。你要有自己的基本的判断，自我暗示也好，坚定信心也罢，**我们要为继续攀登Go高峰进行蓄能**。

## 一. 头脑风暴一下Go成功的因素

相信无论针对哪个gopher群体做头脑风暴，让大家列举Go成功的因素，大家的主流答案也无外乎下图中这些：

![Image](https://mmbiz.qpic.cn/mmbiz_png/cH6WzfQ94mbE1j1BTNg6rvgTBiap4eicP2Um543edKDh112IL0VPc6YibenibyGT5ibM6LD3zVaqn24Uf8SRiaG368wg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

图中的各个因素与Go的成功都不无干系，但是究竟哪个或哪几个是决定性的呢？

## 二. Go成功的根本因素

很显然，这个问题是没有标准答案，是见仁见智的。这里我列举一下我的观点，供大家参考。

直接上结论，我认为**Go成功的根本因素就一个：Google**。

为什么这么说呢？下面我们展开来看(见下图)！

![Image](https://mmbiz.qpic.cn/mmbiz_png/cH6WzfQ94mbE1j1BTNg6rvgTBiap4eicP2yia8coAxAh7sjOxOQerxXmdWHwAI8In6llPBDxbBPsNBucvicjfBadiaQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

我将Go社区比做一支军队，而Go就是Go社区的武器，与其他编程语言搏杀，占地盘(fans)。下面我们就来解构一下这支军队的构成以及为什么这支军队目前有诸多成功案例！

### 1. Google为Go社区提供了统帅与武器

众所周知，2007年Google的三名员工Robert Griesemer、Rob Pike和Ken Thompson(retire很早，精神上领袖，给予Go名誉上的背书)一起发明了Go语言，2009年Go开源后，Go社区逐渐形成。统帅是一支军队的灵魂，他们做出了影响Go和Go社区的最初的也是最重要的决策和这背后的Go设计哲学！

#### a) 设计决策

在2022年，Go团队在**美国计算机学会通讯(Communications of the ACM)**[1]期刊上发表paper：《**Go编程语言与环境**[2]》，对当年做出的诸多决策做了细致说明，这里对其中两个最重要的决策做简单说明：

- Go旨在成为一个编程环境

Go语言之父们认为语言特性仅是编程语言的一部分，而编程环境特性与语言特性同等重要，这些环境特性包括：库、工具、惯例和针对软件工程的整体做法，它们都对使用Go语言编程提供了支持，不可或缺。而这些环境特性恰恰是在传统的编程语言设计中并没有受到应有的重视的。

这样的决策让Go在开源之初就为开发者提供了使用Go进行编程所需的几乎一切：包括功能丰富、开箱即用的标准库以及全面的工具集，代码格式化、代码静态检查、依赖关系管理、构建(包括跨平台交叉编译)、测试、性能剖析、查看和生成文档等，并且这些工具集在今天都统一放在了go命令的下面。这个决策也帮助Go在开源后吸引了第一批Go社区用户。

- Go的一致性的表现

Go的一个目标是**让它在不同的实现、执行环境中，甚至在不同的时间内表现出相同的行为**。所以，Go语言尽可能地规定了一致的结果。比如：Go程序生命周期内一致的性能(相对于使用JIT慢启动的语言）、一致的GC的开销等。甚至对于最常见的编程错误提供了明确定义的语义，这有助于可理解性和调试，而不是像C/C++中那样，充斥着各种未定义的行为。

而我认为最重要的一致性则是从2012年发布的Go 1.0开始，Go团队**公开承诺只对语言和标准库进行向后兼容的修改**，这样程序在编译到较新的Go版本时可以继续运行而不发生变化。这一承诺对业界产生了吸引力，它不仅鼓励了那些长声明周期的工程项目(比如Google内部的一些大型项目或者像Kubernetes这样的社区顶级项目)，也鼓励了其他努力，如书籍、培训课程和第三方软件包的繁荣生态系统。这一一致性的决策也为Go招募了相当数量的拥趸。Go1兼容性，同样可以避免社区分裂（像python2/python3那样），即便是**10多年来变更最大的泛型语法落地**[3]，也没有违反Go1兼容性，这实属不易。

#### b) 设计哲学

上述的设计决策的背后蕴含着Go语言之父们的设计哲学。

- 简单

Tony Hoare在1980年图灵奖演讲中说了这样的观点：“我的结论是，构建软件设计有两种方法：一种方法是让它变得如此简单，显然没有缺陷，另一种方法是让它变得如此复杂，以至于没有明显的缺陷。第一种方法要困难得多。它需要同样的技能，奉献，洞察力，甚至灵感，就像发现作为自然复杂现象基础的简单物理定律一样。它还要求愿意接受受物理，逻辑和技术限制的目标，并在无法实现冲突目标时接受妥协。”

Go选择的正是Tony Hoare演进中的第一种构建软件的设计方法。Rob Pike说过的一句Go流行谚语"less is exponentially more"与此异曲同工。Go的语法简单，API简单，这些为Gopher提供了极大便利，但这些简单的背后其实是Go团队长时间的复杂的思考与实现，努力将语法和API简化为最小、最有用、最接近本质的努力工作。

同时，简单意味着可读性、可维护性，意味着代码的清晰。另一句Go谚语“Clear is better than clever”告诫Gopher们编写平淡如水的Go代码才是“政治正确”的，不要炫技。

- 并发

多核时代，Go将并发作为语言内置特性。Go内置并发原语，包括goroutine、channel、select等。

Go鼓励在较高级别使用并发性，特别是通过通信的方式。我们耳熟能详的一句Go谚语是“Don't communicate by sharing memory. Share memory by communicating”就是并发哲学的外在体现。

- 组合

Go拥有类型，类型可以有method，这似乎像是一种面向对象style的实现，但Go并没有OO语言那种类型层次体系(type hierarchy)，在Go中，组合才是Go类型之间建立联系的最主要手段，而interface和类型嵌入恰是这种组合哲学的具体体现。

- 面向工程

2012年, Go开源元年，Rob Pike就在SPLASH 2012大会上以**“Google的Go：为软件工程服务的语言设计”**[4]为题，讲解了Go是如何围绕Google内部存在的软件工程问题进行有针对性的语言设计的。可以看出，Go从诞生伊始就将解决软件工程领域问题作为语言的目标。同时，我们看到面向工程这个哲学与上面的旨在成为一个编程环境的决策息息相关。

除了统帅之外，Go社区的治理架构也是以Google“将领”为核心的，我们继续来看。

### 2. Google出钱：以Google“将领”(googler and ex-googler)为核心的Go社区治理架构

Go开源10年了，Go社区形成了以Googler和ex-googler(前google员工)为核心的Go社区治理架构，这些人就是上图中的那些“将领”，他们是Go项目某个细分领域，比如：编译器、运行时goroutine调度、GC、内存管理、网络、安全等的领头人。根据Go项目一名产品经理的描述：**2021年，Google Go项目的专职人员多达50多人**[5]，Google这个“亲爹”在金钱的投入上显然表现的十分大方，不得不承认：**在编程语言领域里，有个有钱的“亲爹”就是好**[6]。

这种以googler和Ex-googler为开源社区治理核心的架构决定了Go社区采用的是一种我称之为“民主集中制”的决策机制。在Go社区你不要幻想会有绝对的公平投票，Go项目决策向来是由少数Googler和ex-googler主导的。这样意味着很多情况下，核心治理团队的人提出的proposal以及Google内部gopher提出proposal很容易被accept，而来自外部社区的proposal要想被accept，可能难度就要大一些。怎么说呢？Google的方案不一定总是最好的，但我们也不能不承认多数情况下，Googler提的proposal还是更优的，并且通常这些proposal对应的实现都已经在google内部测试过了，甚至和Go决策组在公司内部“吹过风”，如果你是Go社区的决策人，你会怎么做呢？你是更相信Googler还是外部一个没有任何背景的gopher呢？

我觉得在Google依然引领IT前沿的今天以及未来若干年，这种机制可能还是有利于Go的蓬勃发展的。

### 3. Google为Go社区提供战场/试验场

就像上面所说的那样，Go是有着非常鲜明Google烙印的编程语言，除了Go语言之父都来自google，Go社区治理架构的核心都来自Google和前google员工外，**Google内部为Go的设计提供了足够的一流的问题域，也为Go的真实应用提供了试验场和真实战场**，即便Go至今没有成为Google内部的第一语言。面向Google的一手且一流问题域，让Go设计者和Go开发者能够获得一手的反馈，从而对Go做进一步的打磨。

举几个例子：

- Google内部的单一代码仓库让Go最初设计了不带版本的go get(后在社区的强烈要求下引入了go module，go get才支持版本号)；
- googler反馈，google内部工具超好用，这一定程度让Go团队认识到向Gopher提供完善的go工具链的重要性；
- Google内部的多核与网络服务让Go设计者决定内置原生goroutine以应对多核时代的应用开发；
- Google内性能与开发效率并重让Go设计者决定设计一门带gc的静态编程语言，将内存管理、并发管理下沉到runtime，这与近两年出现的服务网格, dapr等概念的思路一致；



- Google内部大规模人员协作让Go决定面向软件工程，不仅要设计好语言特性，还要提供体验良好的编程环境(工具链、标准库等)；
- Google超大规模的系统构建慢让Go决定提供快速的构建能力，为此对包格式与包依赖做了精心的设计；
- Google内部长期维护的系统(生命周期长) 让Go团队决定支持Go1兼容性并提供支持重构的语法，比如type alias等；
- Google认为安全十分重要，促使Go提供了go sumdb和对sbom的良好支持；

同时Google内部系统为了支持Go的内部试验也是不遗余力，比如：每当Go发布大版本的RC版本，甚至是Beta版本时，Google App Engine都会首当其冲的充当“小白鼠”，在生产环境支持尚未发布正式版的Go。

另外Google在业内的领先性也让“近水楼台”的Go受益，比如像容器调度编排这样的平台，Google十年前就有了(borg)，后续Googler以另起开源项目的方式将其中经验外溢输出，让Kubernetes最终选择了Go作为开发语言，从而成为Go的最大的也是最典型的成功战例。

综上，我们看到Google对Go成功的决定性作用，这种作用可决不能被理解为简单的金钱上的支撑。

## 三. Go语言演进历史

进入Go高级阶段后，对Go语言的演化历史要知道，当然能做到如数家珍更佳，即便不能，也要能记住Go语言的主要演化历史：

- 2007年9月，Go语言诞生；
- 2009年11月，Go正式开源；
- 2012年3月，Go 1.0发布，同时Go1兼容性承诺官宣；
- 2014年12月，Go 1.4版本发布，这是最后一个编译器和runtime由C语言实现的版本；
- 2015年8月，Go 1.5版本发布，这个版本Go实现了自举(用go编译go)，同时编译器和runtime中的绝大部分c代码都换成了go，新版gc让延迟大幅降低；
- 2018年8月，Go 1.11版本发布，go module被正式引入；
- 2022年3月，Go 1.18版本发布，Go泛型语法正式落地。

## 四. 小结

C++之父说过：“世上只有两种编程语言：一种是总是被人抱怨的，一种是从来没人用的”。

Go属于前者。世界上没有完美的编程语言，Go经过十年的打磨已经有了长足的进步，并且取得了不错的战绩，尤其是在云基础设施和云原生因公领域，就连Rob Pike也承认**Go确实已成为云基础架构的语言**[7]。而这个Go走向成功的过程中，Google起着根本性的作用。

不过中国古语有云：成也萧何，败也萧何！目前Google仍然引领IT技术前沿，这对Go的发展来说是一个利好，也会不断推动Go向着好的方向发展。

但我大胆预测一下：“成也Google，败也Google”，一旦Google开始走下坡路的那天，Go语言成功的根基就不在了，Go还能像今天这样顺风顺水么？如果Go社区治理结构不重构，很可能不会再有今天这样的良好状态。大家觉得呢？

## 五. 参考资料

-《**Go编程语言与环境：万字长文复盘导致Go语言成功的那些设计决策**[8]》 -《Go内存模型》- https://research.swtch.com/gomm -《Go语言真正的问题》 - https://vanitynotes.com/posts/20221101-the-real-problem-with-go



### 参考资料

[1] 美国计算机学会通讯(Communications of the ACM): *https://cacm.acm.org/*[2] Go编程语言与环境: *https://tonybai.com/2022/05/04/the-paper-of-go-programming-language-and-environment*[3] 10多年来变更最大的泛型语法落地: *https://tonybai.com/2022/04/20/some-changes-in-go-1-18*[4] “Google的Go：为软件工程服务的语言设计”: *https://go.dev/talks/2012/splash.article*[5] 2021年，Google Go项目的专职人员多达50多人: *https://tonybai.com/2022/01/16/the-2021-review-of-go-programming-language*[6] 在编程语言领域里，有个有钱的“亲爹”就是好: *https://tonybai.com/2012/10/08/the-new-age-of-programming-language/*[7] Go确实已成为云基础架构的语言: *https://tonybai.com/2020/05/01/rob-pike-interview-go-become-the-language-of-cloud-infrastructure/*[8] Go编程语言与环境：万字长文复盘导致Go语言成功的那些设计决策: *https://tonybai.com/2022/05/04/the-paper-of-go-programming-language-and-environment*

------

**推荐阅读**

- [一文彻底理解Go语言栈内存/堆内存](http://mp.weixin.qq.com/s?__biz=MzAxMTA4Njc0OQ==&mid=2651453875&idx=1&sn=883e1b4ac26d62e2d15f96b426885cb0&chksm=80bb2741b7ccae5749f5e7dabc3607f01d88572efe5bbb915db6f0aa235c1bc923cbd4597edc&scene=21#wechat_redirect)

