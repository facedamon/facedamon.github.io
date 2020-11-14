
## 什么是观察者模式

> 观察者观察被观察者，被观察者通知观察者

&emsp;&emsp;我们用"订阅通知"翻译下[观察者模式]的概念：

> "订阅者订阅主题，主题通知订阅者"

&emsp;&emsp;我们再来拆解下这句话，得到：

- 两个对象
    - 被观察者 -> 主题
    - 观察者   -> 订阅者

- 两个动作
    - 订阅     -> 订阅者订阅主题
    - 通知     -> 主题发生变动通知订阅者

&emsp;&emsp;观察者模式的优势：

- 高内聚 -> 不同业务代码变动互不影响
- 可复用 -> 新的业务(就是新的订阅者)订阅不同接口(主题，就是这里的接口)
- 极易扩展 -> 新增接口(就是新增主题);新增业务(就是新增订阅者)


## 什么真实业务场景可以用

> 所有发生变更，需要通知的业务场景

&emsp;&emsp;比如，订单逆向流，也就是订单成立之后的各种取消操作，主要有如下取消类型：

|订单消息类型|
|------|
|未支付取消订单
|超时关单
|已支付取消订单
|拒收

&emsp;&emsp;在触发这些取消操作都要进行各种各样的子操作，显而易见不同的取消操作所涉及的子操作是存在交集的。其次，已支付取消订单的子操作应该是所有订单取消类型最全的，其它类型的服用代码即可，除了分装成函数片段，还有什么更好的封装方式吗？

&emsp;&emsp;接着我们来分析下订单逆向流业务的变与不变：

- 变：
    - 新增取消类型
    - 新增子操作
    - 修改某个子操作的逻辑
    - 取消类型和子操作的对应关系

- 不变
    - 已存在的取消类型
    - 已存在的子操作

## 怎么用

&emsp;&emsp;**分四步走:**
1. 业务梳理
2. 业务流程图
3. 代码建模
4. 代码demo


### 业务梳理

&emsp;&emsp;梳理出所有存在的逆向业务的子操作：

|所有子操作|
|------|
|修改订单状态
|记录订单状态变更日志
|退优惠券
|还优惠活动资格
|还库存
|还礼品卡
|退钱包余额
|修改发货订单状态
|记录发货单状态变更日志
| 生成退款单
| 生成发票
| 发邮件
| 发短信
| 发微信消息

&emsp;&emsp;找到不同订单取消类型和这些子操作的关系：

|订单取消类型(主题)|子操作(订阅者)|
|------|------|
|取消未支付订单||
|-|修改订单状态|
|-|记录订单状态变更日志|
|-|退优惠券|
|-|还优惠活动资格|
|-|还库存|
|超时关单||
|-|修改订单状态|
|-|记录订单状态变更日志|
|-|退优惠券|
|-|还有会活动资格|
|-|还库存|
|-|发邮件|
|-|发短信|
|-|发微信消息|
|已支付取消订单(未生成发货单)||
|-|取消订单状态|
|-|记录订单状态变更日志|
|-|还优惠活动资格|
|-|还库存|
|-|还礼品卡|
|-|退钱包余额|
|-|生成退款单|
|-|生成发票|
|-|发邮件|
|-|发短信|
|-|发微信消息|
|取消发货单(未发货)||
|-|修改订单状态|
|-|记录订单状态变更日志|
|-|修改发货单状态变更日志|
|-|还库存|
|-|还礼品卡|
|-|退钱包余额|
|-|生成退款单|
|-|生成发票|
|-|发邮件|
|-|发短信|
|拒收||
|-|修改订单状态|
|-|记录订单状态变更日志|
|-|修改发货单状态|
|-|记录发货单状态变更日志|
|-|还库存|
|-|还礼品卡|
|-|退钱包余额|
|-|生成退款单|
|-|生成发票|
|-|发邮件|
|-|发短信|

&emsp;&emsp;结论：

- 不同的订单取消类型的子操作存在交集，子操作可被复用
- 子操作可被看作"订阅者"
- 订单取消类型可被看作"主题"
- 不同子操作订阅订单取消类型
- 订单取消类型通知子操作

### 代码建模

&emsp;&emsp;[观察者模式]的核心是两个接口：

- "主题"接口`Observable`
    - 抽象方法`Attach`: 增加订阅者
    - 抽象方法`Detach`: 删除订阅者
    - 抽象方法`Notify`: 通知订阅者
- "订阅者"接口`ObserverInterface`
    - 抽象方法`Do`: 自身业务

&emsp;&emsp;订单逆向流的业务下，我们需要实现这两个接口：

- 具体订单取消的动作实现"主题"接口`Observable`
- 子逻辑实现`订阅者`接口`ObserverInterface`

&emsp;&emsp;来看看伪代码：

        // ----------这里实现一个具体的"主题"----------
        具体订单取消的动作实现"主题"接口`Observable`.得到一个具体的主题：

        - 订单取消的动作的"主题"结构体`ObservableConcrete`
            + 成员属性`observerList []ObserverInterface`: 订阅者列表
            + 具体方法`Attach`: 增加订阅者
            + 具体方法`Detach`: 删除订阅者
            + 具体方法`Notify`: 通知订阅者

        // ----------这里实现所有具体的"订阅者"---------

        子逻辑实现"订阅者"接口`ObserverInterface`:

        - 具体订阅者`OrderStatus`
            + 实现方法`Do`: 修改订单状态
        - 具体订阅者`OrderStatusLog`
            + 实现方法`Do`: 记录订单状态变更日志
        - 具体订阅者`CouponRefund`
            + 实现方法`Do`: 退优惠券
        - 具体订阅者`PromotionRefund`
            + 实现方法`Do`: 还优惠活动资格
        - 具体订阅者`StockRefund`
            + 实现方法`Do`: 还库存
        - 具体订阅者`GiftCardRefund`
            + 实现方法`Do`: 还礼品卡
        - 具体订阅者`WalletRefund`
            + 实现方法`Do`: 退钱包余额
        - 具体订阅者`DeliverBillStatus`
            + 实现方法`Do`: 修改发货单状态
        - 具体订阅者`DeliverBillStatusLog`
            + 实现方法`Do`: 记录发货单状态变更日志
        - 具体订阅者`Refund`
            + 实现方法`Do`: 生成退款单
        - 具体订阅者`Invoice`
            + 实现方法`Do`: 生成发票
        - 具体订阅者`Email`:
            + 实现方法`Do`: 发邮件
        - 具体订阅者`Sms`
            + 实现方法`Do`: 发短信

### 代码demo

        package main

        import (
            "fmt"
            "reflect"
            "runtime"
        )

        // 主题
        type Observable interface {
            Attach(observer ...ObserverInterface) Observable
            Detach(observer ObserverInterface) Observable
            Notify() error
        }

        // 主题实现
        type ObservableConcrete struct {
            observerList []ObserverInterface
        }

        // 获取正在运行的函数名
        func runFuncName() string {
            pc := make([]uintptr, 1)
            runtime.Callers(2, pc)
            f := runtime.FuncForPC(pc[0])
            return f.Name()
        }

        func (o *ObservableConcrete) Attach(observerInterface) Observable {
            o.observerList = append(o.observerList, observer...)
            return o
        }

        func (o *ObservableConcrete) Detach(observer ObserverInterface) Observable {
            if len(o.observerList) == 0 {
                return o
            }
            for k, observerItem := range o.observerList {
                if observer == observerItem {
                    fmt.Printf("%s 注销: %s\n", runFuncName(), reflect.TypeOf(observer))
                    o.observerList = append(o.observerList[:k], o.observerList[k+1:]...)
                }
            }
            return o
        }

        func (o *ObserverableConcrete) Notify() (err error) {
            for _, observer := range o.observerList {
                if err = observer.Do(o); err != nil {
                    return err
                }
            }
            return nil
        }

        // 订阅者
        type ObserverInterface interface {
            Do(o Observable) error
        }

        type OrderStatus struct{}
        func (observer *ObserStatus) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 修改订单状态...", runFuncName())
            return
        }

        type OrderStatusLog struct{}
        func (observer *OrderStatusLog) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 记录订单状态变更日志...", runFuncName())
            return
        }

        type CouponRefund struct {}
        func (observer *CouponRefund) Do(o Observable) (err error) {
            /// code...
            fmt.Printf("%s, 退优惠券...", runFuncName())
            return
        }

        type PromotionRefund struct {}
        func (observer *PromotionRefund) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 还优惠活动资格...", runFuncName())
            return
        }

        type StockRefund struct {}
        func (observer *StockRefund) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 还库存...", runFuncName())
            return
        }

        type GiftVardRefund struct {}
        func (observer *GiftCardRefund) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 还礼品卡...", runFuncName())
            return
        }

        type WalletRefund struct {}
        func (observer *WalletRefund) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 退钱包余额...", runFuncName())
            return
        }

        type DeliverBillStatus struct {}
        func (observer *DeliverBillStatus) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 修改发货单状态...", runFuncName())
            return
        }

        type DeliverBillStatusLog struct {}
        func (observer *DeliverBillStatusLog) Do(o Observable) (err error) {
            // code..
            fmt.Printf("%s, 记录发货单状态变更日志...", runFuncName())
            return
        }

        type Refund struct {}
        func (observer *Refund) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 生成退款单...", runFuncName())
            return
        }

        type Invoice struct {}
        func (observer *Invoice) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 生成发票...", runFuncName())
            return
        }

        type Sms struct {}
        func (observer *Sms) Do(o Observable) (err error) {
            // code...
            fmt.Printf("%s, 发短信...", runFuncName())
            return
        }


        func main() {
            // 创建 未支付去掉订单主题
            orderUnPaidCancelSubject := &ObservableConcrete{}
            orderUnPaidCancelSubject.Attach(
                &OrderStatus{},
                &OrderStatusLog{},
                &CouponRefund{},
                &PromotionRefund{},
                &StockRefund{},
            )
            orderUnPaidCancelSubject.Notify()
        }

        // 创建 超时关单主题
        orderOverTimeSubject := &ObservableConcrete{}
        orderOverTimeSubject.Attach(
            &OrderStatus{},
            &OrderStatusLog{},
            &CouponRefund{},
            &PromotionRefund{},
            &StockRefund{},
            &Sms{},
        )
        orderOverTimeSubject.Notify()

        // 创建 已支付取消订单主题
        orderPaidCancelSubject := &OrdervableConcrete{}
        orderPaidCancelSubject.Attach(
            &OrderStatus{},
            &OrderStatusLog{},
            &CouponRefund{},
            &PromotionRefund{},
            &StockRefund{},
            &GiftCardRefund{},
            &WalletRefund{},
            &Refund{},
            &Invoice{},
            &Sms{},
        )
        orderPaidCancelSubject.NOtify()

        // 创建 取消发货单主题
        deliverBillCancelSubject := &ObservableConcrete{}
        deliverBillCancelSubject.Attach(
            &OrderStatus{},
            &OrderStatusLog{},
            &DeliverBillStatus{},
            &DeliverBillStatusLog{},
            &StockRefund{},
            &GiftCardRefund{},
            &WalletRefund{},
            &Refund{},
            &Invoice{},
            &Sms{},
        )
        deliverBillCancelSubject.Notify()

        // 创建 拒收主题
        deliverBillRejectSubject := OrdervableConcrete{}
        deliverBillRejectSubject.Attach(
            &OrderStatus{},
		&OrderStatusLog{},
		&DeliverBillStatus{},
		&DeliverBillStatusLog{},
		&StockRefund{},
		&GiftCardRefund{},
		&WalletRefund{},
		&Refund{},
		&Invoice{},
        &Sms{},
        )
        deliverBillRejectSubject.Notify()

        // 未来可以快速的根据业务的变化 创建新的主题
        // 从而快速构建新的业务接口


## 结语

&emsp;&emsp;最后总结一下，[观察者模式]抽象过程的核心是：

- 被依赖的"主题"
- 被通知的"订阅者"
- "订阅者"按需订阅"主题"
- "主题"变化通知"订阅者"

## 变化

&emsp;&emsp;[观察者模式]最显著的特点是：消息的生产者和消费者是直连关系。下一期我们将看看[发布订阅模式]是如何处理这类问题的。