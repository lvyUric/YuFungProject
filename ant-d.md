开始使用
Ant Design Pro 是基于 Ant Design 和 umi 的封装的一整套企业级中后台前端/设计解决方案，致力于在设计规范和基础组件的基础上，继续向上构建，提炼出典型模板/业务组件/配套设计资源，进一步提升企业级中后台产品设计研发过程中的『用户』和『设计者』的体验。随着『设计者』的不断反馈，我们将持续迭代，逐步沉淀和总结出更多设计模式和相应的代码实现，阐述中后台产品模板/组件/业务场景的最佳实践，也十分期待你的参与和共建。

开发前的输入
Ant Design Pro 作为一个前端脚手架，默认读者已经懂了一些前端的基础知识，并且了解 umi 和 Ant Design, 如果你是纯粹的新手，第一次来跑项目建议读一下 新手需知。磨刀不误砍柴工，了解一些基础知识可以让学习曲线更加平滑。

准备工作
由于国内网络和前端的特殊性，在安装依赖等方面可能会失败或导致无法启动，浪费大量的时间，因此我们推荐如下的技术栈来帮助我们顺畅的开发。

包管理器
推荐使用 tyarn 来进行包管理，可以极大地减少 install 的时间和失败的概率，并且完全兼容 npm。

如果喜欢使用 npm 而不是 yarn，可以使用 cnpm, 安装速度比 tyarn 更快，但是与 npm 不是完全兼容。

Terminal
非 Windows 用户没什么好说的，iTerm2 和 Oh My Zsh 是最强选择。

对于 Windows 用户而言，Bash on Linux 是最好的选择，但是可能会造成一些性能问题。这里推荐使用 Windows Terminal 和 PowerShell。Windows Terminal 可以直接在微软商店中下载，美貌与实力并存，不逊于 iTerm2，微软官方维护品质也值得信赖。PowerShell 是 Windows 7 以来内置的命令行工具，被 Linux 创始人称赞为不那么烂的命令行。并且可以配置 posh-git，能得到部分 zsh 的能力。

配置好后效果

PowerShell

初始化
我们提供了 pro-cli 来快速的初始化脚手架。

# 使用 npm
npm i @ant-design/pro-cli -g
pro create myapp
选择 umi 的版本

? 🐂 使用 umi@4 还是 umi@3 ? (Use arrow keys)
❯ umi@4
umi@3
如果选择了 umi@4 版本，暂时还不支持全量区块。

如果选择了 umi@3，还可以选择 pro 的模板，simple 是基础模板，只提供了框架运行的基本内容，complete 包含所有区块，不太适合当基础模板来进行二次开发

? 🚀 要全量的还是一个简单的脚手架? (Use arrow keys)
❯ simple
complete
安装依赖：

$ cd myapp && tyarn
// 或
$ cd myapp && npm install
开发
脚手架初始化成功之后就可以开始进行开发了，我们提供了一些命令来辅助开发。

start
运行这个脚本会启动服务，自动打开默认浏览器展示你的页面。当你重新编辑代码后，页面还会自动刷新。

start

build
运行这个脚本将会编译你的项目，你可以在项目中的 dist 目录中找到编译后的文件用于部署。

build

如果你需要部署，可以查阅部署。

analyze
analyze 脚本做的事情与 build 的相同，但是他会打开一个页面来展示你的依赖信息。如果需要优化性能和包大小，你需要它。

analyze

lint
我们提供了一系列的 lint 脚本，包括 TypeScript，less，css，md 文件。你可以通过这个脚本来查看你的代码有哪些问题。在 commit 中我们自动运行相关 lint。

lint

lint:fix
与 lint 相同，但是会自动修复 lint 的错误。

lint:fix

i18n-remove
这个脚本将会尝试删除项目中所有的 i18n 代码，对于复杂的运行时代码，表现并不好，慎用。

更多的命令，可以看这里 https://umijs.org/zh-CN/docs/cli#umi-build
