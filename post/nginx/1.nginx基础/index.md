
## 准备环境

&emsp;&emsp;首先需要一个内核为Linux 2.6及以上版本的操作系统，因为Linux 2.6及以上内核才支持epoll。在linux上使用select或poll来解决事件的多路复用，是无法解决高并发压力问题的。

		$ uname -a
		Linux VM_0_8_centos 3.10.0-1062.9.1.el7.x86_64 #1 SMP Fri Dec 6 15:49:49 UTC 2019 x86_64 x86_64 x86_64 GNU/Linux

## 必备软件

- GCC

&emsp;&emsp;GCC可用来编译C语言程序。我们可以使用yum或者采用tar.gz源码的方式安装，由于我们的侧重点是Nginx，所以在此我们采用yum安装

		$ sudo yum install -y gcc

&emsp;&emsp;如果想使用C++来编写Nginx HTTP模块，这时就需要用到G++编译器了。

		$ sudo yum install -y gcc-c++

- PCRE

&emsp;&emsp;PCRE(Perl Compatible Regular Expressions)是由Philip Hazel开发的函数库，目前为很多软件所使用，该库支持正则表达式。它由RegEx演化而来，实际上，Perl正则也是源自Henry Spencer写的RegEX。

&emsp;&emsp;如果我们在配置文件nginx.conf中使用了正则表达式，那么我们在编译Nginx时就必须把Pcre库编译进Nginx。

		$ sudo yum install -y pcre pcre-devel

- zlib

&emsp;&emsp;zlib库用于对HTTP包的内容做gzip格式的压缩，如果我们在nginx.conf里配置了gzip on，并指定对于某些类型(content-type)的HTTP响应使用gzip来进行压缩以减少网络传输量，那么，在编译时就必须把zlib编译进Nginx

		$ sudo yum install -y zlib zlib-devel

- OpenSSL

&emsp;&emsp;如果我们的服务器需要在更安全的SSL协议上传输HTTP，那么就需要拥有OpenSSL。另外，如果我们想要使用MD5、SHA1等散列函数，那么也需要安装它。

		$ sudo yum install -y openssl openssl-devel

---------------------------------------------------------

&emsp;&emsp;上面所列的4个库致使完成Web服务器最基本功能所必需的。Nginx是高度自由化的web服务器，它的功能是由需要模块来支持的。而这些模块可根据我们的使用需求来定制，如果某些模块不需要使用则完全不必理会它。同样，如果使用了某个模块，而这个模块使用了一些类似zlib或者OpenSSL等的第三方库，那么就必须先安装这些软件。

## 了解磁盘目录

- Nginx源码存放目录 -> 一般存放在/usr/local/src
- Nginx编译阶段产生的中间文件存放目录 -> /usr/local/src/nginx/objs
- 部署目录 -> /usr/local/nginx
- 日志目录 -> /usr/local/nginx/logs
- 更多目录更查看[configure参数详解]()

## 获取源码

		$ sudo mkdir /usr/local/src
		$ cd /usr/local/src
		$ sudo wget http://nginx.org/download/nginx-1.0.14tar.gz
		$ sudo tar -zxvf nginx-1.0.14.tar.gz

## configure详解

- configure支持的路径相关参数

|参数名称|意义|默认值|
|:---|:---:|:---:|
|\--prefix=PATH|Nginx安装部署后的根目录|默认为/usr/local/nginx.注意:这个目标的设置会影响其它参数中的相对目录。例如，如果设置了--sbin-path=sbin/nginx,那么实际上可执行文件会被放到/usr/local/nginx/sbin/nginx|
|\--sbin-path=PATH|可执行文件的放置路径|\<prefix\>/sbin/nginx|
|\--conf-path=PATH|配置文件的放置路径|\<prefix\>/conf/nginx.conf|
|\--error-path=PATH|error日志的放置路径,error日志用于定位问题，可输出多种级别(包括debug调试级别)的日志.它的配置非常灵活，可以在nginx.conf里配置为不同请求的日志并输出到不同的log文件中.这里是默认的nginx核心日志路径|\<prefix>/logs/error.log|
|\--pid-path=PATH|pid文件的存量路径，这个文件里仅以ASCII码存放这Nginx master的进程ID。有了这个进程ID，在使用命令行(eg:nginx -s reload)通过读取master进程ID向master进程发送信号时，才能对运行中的Ngin服务产生作用|\<prefix>/logs/nginx.pid|
|\--lock-path=PATH|lock文件的放置路径|\<prefix>/logs/nginx.lock|
|\--builddir=DIR|configure执行时与编译期间产生的临时文件放置的目录，包括产生的MakeFile,C源文件，目标文件，可执行文件|\<nginx source path>/objs|
|\--with-perl=PATH|perl binary放置路径,如果配置的nginx会执行Perl脚本，那么就必须要设置此路径|无|
|\--http-log-path=PATH|access日志放置位置，每一个HTTP请求在结束时都会记录的访问日志|\<prefix>/logs/access.log|
|\--http-client-body-temp-path=PATH|处理HTTP请求时如果请求的包体需要暂时存放到临时磁盘文件中，则把这样的临时文件放置到该路径下|\<prefix>/client_body_temp|
|\--http-proxy-temp-path=PATH|Nginx作为HTTP反向代理服务器时，上游服务器产生的HTTP包体在需要临时存放在磁盘文件时，这样的临时文件将存放到该路径下|\<prefix>/proxy_temp|

- configure支持的编译相关参数

|编译参数|意义|
|:--|:--:|
|\--with-cc=PATH|C编译器的路径|
|\--with-cpp=PATH|C与编译器的路径|
|\--with-cc-opt=OPTIONS|如果希望在Nginx编译期间指定加入一些编译选项，如指定宏或者使用-I加入某些需要包含的目录，这时可以使用该参数达成目的|
|\--with-Id-opt=OPTIONS|最终的二进制可执行文件是由编译后生成的目标文件与一些第三方库链接生成的，在执行链接操作时可能会需要指定链接参数，\--with-Id-opt就是用于加入链接时的参数，例如，如果我们希望将某个库链接到Nginx中，需要在这里加入\--with-Id-opt=libraryName-LibraryPath,其中libraryName是目标库的名称，librarypath则是目标库所在路径|
|\--with-cpu-opt=CPU|指定CPU处理器架构，只能从以下取值中选择：pentium, pentiumpro, pentium3, pentium4, athlon, opteron, sparc32, sparc64, ppc64|

- PCRE参数设置

|PCRE库的设置参数|意义|
|:--|:--:|
|\--without-pcre|如果确认Nginx不用解析正则，nginx.conf配置文件中不会出现正则，那么可以使用这个参数|
|\--with-pcre|强制使用PCRE库|
|\--with-pcre=DIR|指定PCRE库的源码位置，在编译Nginx时会进入该目录编译PCRE源码|
|\--with-pcre-opt=OPTIONS|编译PCRE源码时希望加入的编译选项|

- OpenSSL参数设置

|OpenSSL库的设置参数|意义|
|:--|:--:|
|\--with-openssl=DIR|指定OpenSSL库的源码位置，在编译Nginx时会进入该目录编译OpenSSL源码。如果要使用HTTPS，则必须使用OpenSSL|
|\--with-openssl-opt=OPTIONS|编译OpenSSL源码时希望加入的编译选项|

- 散列函数库设置参数

|散列函数库|意义|
|:--|:--:|
|\--with-MD5=DIR|指定MD5库源码位置，在编译时Nginx会进入该目录编译MD5源码。注意：Nginx源码内部已经有了MD5的算法实现，如果没有特殊需求，那么完全可以使用Nginx自身实现的MD5算法|
|\--with-MD5-opt=OPTIONS|编译MD5源码时希望加入的编译选项|
|\--with-MD5-asm|使用MD5的汇编源码|
|\--with-SHA1=DIR|指定SHA1库源码位置，在编译时Nginx会进入该目录编译SHA1源码。注意：OpenSSL中已经有了SHA1算法的实现，如果已经安装了OpenSSL那么完全可以使用OpenSSL实现的SHA1算法|
|\--with-SHA1-opt=OPTIONS|编译SHA1源码时希望加入的编译选项|
|\--with-SHA-asm|使用SHA1的汇编源码|


- zlib库设置参数

|zlib库参数|意义|
|:--|:--:|
|\--with-zlib=DIR|指定zlib库的源码位置，在编译Nginx时会进入该目录编译zlib源码，如果使用了gzip压缩功能，则需要zlib库的支持|
|\--with-zlib-opt=OPTIONS|编译zlib源码时希望加入的编译选项|

### 模块参数

- 事件模块
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/nginx/epub_811745_16.jpeg)

- 默认编译进入Nginx的HTTP模块
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/nginx/epub_811745_17.jpeg)

- 默认不编译进入Nginx的HTTP模块
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/nginx/epub_811745_18.jpeg)

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/nginx/epub_811745_19.jpeg)

- 邮件代理服务器mail模块
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/nginx/epub_811745_20.jpeg)

- 其它模块
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/nginx/epub_811745_21.jpeg)
